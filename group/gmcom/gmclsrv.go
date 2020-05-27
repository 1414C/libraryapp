package gmcom

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"golang.org/x/net/websocket"
)

const (
	CPing        = "PING"
	CAck         = "ACK"
	CJoin        = "JOIN"
	CJoinAck     = "JOINACK"
	CLeave       = "LEAVE"
	CFailure     = "FAILURE"
	CCoordinator = "COORDINATOR"
	COkay        = "OKAY"
	CDeparting   = "DEPARTING"
)

const (
	CinGetDBLeader      = "inGetDBLeader"
	CinSetDBLeader      = "inSetDBLeader"
	CinSetLeader        = "inSetLeader"
	CinGetLocalLeader   = "inGetLocalLeader"
	CinStartElection    = "inStartElection"
	CinRunElection      = "inRunElection"
	CinSetElectionState = "inSetElectionState" // true/false - NEEDED?
	CinGetElectionState = "inGetElectionState" // true/false
	CinGetLocalDetails  = "inGetLocalDetails"  // into Target ID, IPAddress, MemberMap - NEEDED?
	CinGetMySrcInfo     = "inGetMySrcInfo"
	CinGetMyTargetInfo  = "inGetMyTargetInfo"
	CinNoAck            = "inNoAck"
	CinDoSendPrep       = "inDoSendPrep"
	CinFlushMemberMap   = "inFlushMemberMap"
)

// election 'error' codes
const (
	CElectOK         = "OK"
	CElectNR         = "NR" // no response
	CElectComplete   = "COMPLETE"
	CElectIncomplete = "INCOMPLETE"
)

// ProcStatus is the process-status constant type
type ProcStatus string

const (
	CStatusAlive    ProcStatus = "ALIVE"
	CStatusSuspect  ProcStatus = "SUSPECT"
	CStatusFailed   ProcStatus = "FAILED"
	CStatusDeparted ProcStatus = "DEPARTED"
)

// GMMember is the struct for membership recording
type GMMember struct {
	ID                uint
	IPAddress         string
	Status            ProcStatus
	StatusCount       uint
	IncarnationNumber uint
}

// GMMessage is the general message struct
type GMMessage struct {
	Type                    string
	SrcID                   uint
	SrcIPAddress            string
	SrcIncarnationNumber    uint
	SrcStatus               ProcStatus
	SrcStatusCount          uint
	TargetID                uint   // ping-ack
	TargetIPAddress         string // ping-ack
	TargetIncarnationNumber uint
	TargetPath              string
	TargetStatus            ProcStatus
	TargetStatusCount       uint
	SubjectID               uint   // ping-ack etc.
	SubjectIPAddress        string // ping-ack etc.
	ExpectResponse          bool
	MemberMap               *OMap
	CheckElection           bool
	InElection              bool
	ErrCode                 string // optional error code
	Value                   []byte // general carrier field - use Type as a guide to decode
}

// GMLeader holds leader information and is used during the initialization of a process.
type GMLeader struct {
	LeaderID        uint
	LeaderIPAddress string
}

// GMTarget is a convenience struct to hold target process information
type GMTarget struct {
	ID        uint
	IPAddress string
}

// GMLeaderSetterGetter contains methods to facilitate r/w access to the persisted leader
// record.  Applications using this library may choose to store the persisted 'current'
// leader in a number of ways; flat-file, db table record, NVS etc.  GMLeaderSetterGetter
// is used as a parameter in the GMServ.Serve(...) method.
type GMLeaderSetterGetter interface {
	GetDBLeader() (*GMLeader, error)
	SetDBLeader(l GMLeader) error
	Cleanup() error
}

// EncodeGMMessage is used to gob-encode a GMMessage for transmission
// to interested parties.
func EncodeGMMessage(m GMMessage) ([]byte, error) {

	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(m)
	if err != nil {
		return nil, fmt.Errorf("EncodeGMMessage failed to gob-encode the GMMessage: %s", err)
	}
	return encBuf.Bytes(), nil
}

// DecodeGMMessage is used to decode any incoming GMMessage from
// its gob-encoded format into something usable.
func DecodeGMMessage(ws *websocket.Conn) (*GMMessage, error) {

	// gob decoding
	var m GMMessage
	var msg = make([]byte, 2048)
	l, err := ws.Read(msg)
	if err != nil {
		return nil, fmt.Errorf("DecodeGMMessage() ws.Read() error: %s", err)
	}
	raw := msg[0:l]
	decBuf := bytes.NewBuffer(raw)
	err = gob.NewDecoder(decBuf).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("DecodeGMMessage() gob.Decode() error: %s", err)
	}
	return &m, nil
}
