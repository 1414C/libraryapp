package gmsrv

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/1414C/libraryapp/group/gmcom"
	"golang.org/x/net/websocket"

	"github.com/1414C/lw"
)

// TxRxGMMessage sends message m to the target process (Pj) and
// returns a pointer to the message response.  If no response message
// is expected (for example gmcom.Coordinator), the response message
// pointer and error == nil.  If an error occurs, the error return
// parameter will be populated and the response message pointer will
// == nil.
func (gm *GMServ) TxRxGMMessage(m gmcom.GMMessage) (*gmcom.GMMessage, error) {

	// check for self-ping
	selfPing := false
	lw.Debug("SELFPING-0: gm.MyIPAddress: %v m.TargetIPAddress: %v", gm.MyIPAddress, m.TargetIPAddress)
	if gm.MyIPAddress == m.TargetIPAddress {
		lw.Debug("SELFPING-1: Self-ping scenario identified")
		selfPing = true
	}

	// network connection has been down?  The following ping may fail, but the join
	// gets sent and processed anyway.  remember that if the network connection has
	// been down for any amount of time, it is likely that the membermap will be
	// empty.
	if selfPing && gm.MyNetworkOffline == true {
		lw.Debug("SELFPING-2: Self-ping with existing offline network condition")
		lm := &gmcom.GMMessage{}
		gm.SendGetDBLeader(lm)
		lw.Debug("SELFPING-3: gm.SendGetDBLeader() got: %v", lm)
		if lm.TargetID != 0 {
			lw.Debug("SELFPING-4: calling gm.SendSetLeader()...")
			gm.SendSetLeader(lm)
			gm.MyID = 0
			gm.MyIncarnation = 1
			lw.Debug("SELFPING-5: calling gm.SendJoin()...")
			err := gm.SendJoin(&gm.Leader)
			if err != nil {
				lw.Debug("SELFPING-6: gmcl.TxRxGMMessage() attempted to rejoin group. got: %v", err.Error())
			} else {
				lw.Debug("SELFPING-7: setting gm.MyNetworkOffline = false")
				gm.MyNetworkOffline = false
			}
		}
	}

	// set origin and target
	origin := "http://" + gm.MyIPAddress
	url := "ws://" + m.TargetIPAddress + m.TargetPath

	// encode the source message as bytes
	raw, err := gmcom.EncodeGMMessage(m)
	if err != nil {
		panic(fmt.Errorf("gmcl.TxRxGMMessage error: %s", err))
	}

	// connect to the target
	lw.Info("websocket.Dial: %s", url)
	config, _ := websocket.NewConfig(url, origin)
	config.Dialer = &net.Dialer{
		// Deadline: time.Now().Add(-time.Minute),
		Deadline: time.Now().Add(5 * time.Second),
	}

	ws, err := websocket.DialConfig(config)
	if err != nil {
		lw.ErrorWithPrefixString("gm.TxRxGMMessage() ws connection failed - got:", err)
		if ws != nil {
			ws.Close()
		}
		if selfPing {
			lw.Debug("SELFPING-1B: setting gm.MyNetworkOffline = true")
			gm.MyNetworkOffline = true // out of channel / no lock
		}
		return nil, err
	}

	defer ws.Close()
	lw.Info("websocket.Dial: %s complete", url)

	// send the encoded GMMessage
	_, err = ws.Write(raw)
	if err != nil {
		lw.ErrorWithPrefixString("gm.TxRxGMMessage() ws.Write error - got:", err)
		return nil, err
	}

	var msg = make([]byte, 2048)
	raw = raw[:0]

	// single read from the ws is okay here
	l, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("gm.TxRxGMMessage() ws.Read error - got:", err)
		return nil, err
	}

	// decode the returned message
	rm := gmcom.GMMessage{}
	raw = msg[0:l]
	decBuf := bytes.NewBuffer(raw)
	err = gob.NewDecoder(decBuf).Decode(&rm)
	if err != nil {
		lw.ErrorWithPrefixString("gm.TxRxGMMessage() gob.Decode() error - got:", err)
		return nil, err
	}
	return &rm, nil
}

// TxGMMessage sends message m to the target process (Pj) and
// returns an error occurs if there is a technical issue.
func (gm *GMServ) TxGMMessage(m gmcom.GMMessage) error {

	// set origin and target
	origin := "http://" + gm.MyIPAddress
	url := "ws://" + m.TargetIPAddress + m.TargetPath

	// encode the source message as bytes
	raw, err := gmcom.EncodeGMMessage(m)
	if err != nil {
		panic(fmt.Errorf("gm.TxGMMessage error: %s", err))
	}

	// connect to the target
	lw.Info("websocket.Dial: %s", url)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		lw.ErrorWithPrefixString("gm.TxGMMessage() ws connection failed - got:", err)
		if ws != nil {
			ws.Close()
		}
		return err
	}
	defer ws.Close()
	lw.Info("websocket.Dial: %s complete", url)

	// send the encoded GMMessage - successful write is enough validation
	_, err = ws.Write(raw)
	if err != nil {
		lw.ErrorWithPrefixString("gm.TxGMMessage() ws.Write error - got:", err)
		return err
	}
	return nil
}
