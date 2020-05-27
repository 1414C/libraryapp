package gmsrv

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/1414C/libraryapp/group/gmcl"
	"github.com/1414C/libraryapp/group/gmcom"

	"github.com/1414C/lw"
	"golang.org/x/net/websocket"
)

// GMServInt outlines the core group membership interface.
type GMServInt interface {
	processCmdChannel()
	Serve(myIPAddress string, lsg gmcom.GMLeaderSetterGetter, actUsrs *gmcom.ActUsrsH, groupAuths *gmcom.GroupAuthsH, auths *gmcom.AuthsH, usrGroups *gmcom.UsrGroupsH, logging bool, evlCycle uint, failureThreshold uint64)
}

// GMServHandlerInt outlines the group-membership related web-socket handlers.
type GMServHandlerInt interface {
	PingHandler(ws *websocket.Conn)
	JoinHandler(ws *websocket.Conn)
	CoordinatorHandler(ws *websocket.Conn)
	UsrUpdateHandler(ws *websocket.Conn)
}

// GMServSenderInt outlines the message senders.  These methods mostly push messages into the group-membership serialization channels.
type GMServSenderInt interface {
	SendGetDBLeader(m *gmcom.GMMessage)
	SendSetDBLeader(m *gmcom.GMMessage)
	SendSetLeader(m *gmcom.GMMessage)
	SendGetLocalLeader(m *gmcom.GMMessage)
	SendJoin(l *gmcom.GMLeader) error
	SendPing()
	SendFailure(f gmcom.GMMessage)
	SendDeparting()
	SendCoordinator(c gmcom.GMMessage) error
}

// GMServTxRxInt outlines the common message TxRx and message Tx only methods.
type GMServTxRxInt interface {
	TxRxGMMessage(m gmcom.GMMessage) (*gmcom.GMMessage, error)
	TxGMMessage(m gmcom.GMMessage) error
}

// GMServUsrCacheSupportInt outlines convenience methods for the Usr cache
// to get information regarding the state of the group processes.
type GMServUsrCacheSupportInt interface {
	SendGetLocalDetails() *gmcom.GMMessage
}

// GMServ is the server access struct
type GMServ struct {
	Log              bool
	memberMap        *gmcom.OMap
	Leader           gmcom.GMLeader
	MyID             uint
	MyIPAddress      string
	MyIncarnation    uint
	MyNetworkOffline bool // true = unable to ping myself
	FailureThreshold uint32
	chin             chan gmcom.GMMessage
	chout            chan gmcom.GMMessage
	count            int
	HTTPServer       *http.Server
	LSG              gmcom.GMLeaderSetterGetter
	ActUsrsH         *gmcom.ActUsrsH
	GroupAuthsH      *gmcom.GroupAuthsH
	AuthsH           *gmcom.AuthsH
	UsrGroupsH       *gmcom.UsrGroupsH
	InElection       bool // true/false
	GMServInt
	GMServHandlerInt
	GMServSenderInt
	GMServTxRxInt
}

// ensure consistency against interfaces
var _ GMServInt = &GMServ{}
var _ GMServTxRxInt = &GMServ{}
var _ GMServSenderInt = &GMServ{}
var _ GMServHandlerInt = &GMServ{}

// init an empty gm server
func (gm *GMServ) initialize(myIPAddress string, lsg gmcom.GMLeaderSetterGetter, actUsrs *gmcom.ActUsrsH, groupAuths *gmcom.GroupAuthsH, auths *gmcom.AuthsH, usrGroups *gmcom.UsrGroupsH, logging bool, failureThreshold uint32) error {

	//log.SetFlags(0)

	// check that a port was included in myIPAddress (1-255.0-255.0-255.0.255:10-99999)
	match, _ := regexp.MatchString("(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(:{1})([1-9][0-9]|[1-9][0-9]{2}|[1-9][0-9]{3}||[1-9][0-9]{4})$", myIPAddress)
	if !match {
		return fmt.Errorf("supplied IPv4 address:port combination is illegal - got: %s", myIPAddress)
	}

	// this is the number of aggregate SUSPECT notifications the process (Pi) gets for another process (Pj)
	// this will likely change to something like (failureThreshold * number of active processes) -1
	if failureThreshold == 0 {
		gm.FailureThreshold = 30
	} else {
		gm.FailureThreshold = failureThreshold
	}

	// initialize the local server attributes
	gm.MyID = 0
	gm.MyIPAddress = myIPAddress
	gm.MyIncarnation = 1
	gm.LSG = lsg
	gm.Log = logging
	gm.ActUsrsH = actUsrs
	gm.GroupAuthsH = groupAuths
	gm.AuthsH = auths
	gm.UsrGroupsH = usrGroups

	// initialize the members ordered map
	if gm.memberMap == nil {
		gm.memberMap = &gmcom.OMap{}
		gm.memberMap.Init()
	}

	// initialize the command/serialization channel
	gm.chin = make(chan gmcom.GMMessage)
	gm.chout = make(chan gmcom.GMMessage)

	// start blocking on the command/serialization channel
	go gm.processCmdChannel()
	return nil
}

// processCmdChannel is used to process commands from the command/serialization channel.
// operations on gmsrv artifacts are considered goroutine/race safe.
func (gm *GMServ) processCmdChannel() {

	// block on the channel
	for v := range gm.chin {
		gm.count++
		lw.Info("gm.count: %v", gm.count)

		switch v.Type {
		case gmcom.CJoin:
			lw.Info("JOIN received by processCmdChannel()")
			rm := gm.processPing(v)
			gm.chout <- rm

		case gmcom.CJoinAck:
			lw.Info("JOIN has been accepted")
			jm := gmcom.GMMember{
				ID:                v.TargetID,
				IPAddress:         v.TargetIPAddress,
				Status:            gmcom.CStatusAlive, //v.TargetStatus,
				StatusCount:       1,                  // v.TargetStatusCount,
				IncarnationNumber: v.TargetIncarnationNumber,
			}
			gm.memberMap.Add(v.TargetID, jm)
			gm.MyID = v.TargetID
			gm.chout <- v

		case gmcom.CDeparting:
			lw.Info("DEPARTING: message received from: %v %v", v.SrcID, v.SrcIPAddress)
			dm := gmcom.GMMember{
				ID:                v.SrcID,
				IPAddress:         v.SrcIPAddress,
				Status:            gmcom.CStatusDeparted, //v.TargetStatus,
				StatusCount:       1,                     // v.TargetStatusCount,
				IncarnationNumber: v.TargetIncarnationNumber,
			}
			var err error
			v.CheckElection, err = gm.memberMap.UpdFromDeparting(dm, gm.Leader.LeaderID)
			if err != nil {
				lw.ErrorWithPrefixString("error: gm.memberMap.UpdFromDeparting got:", err)
			}
			lw.Debug("========================================================================")
			all := gm.memberMap.ReadAll()
			for _, mem := range all {
				lw.Debug("member: %v", mem)
			}
			lw.Debug("========================================================================")
			gm.chout <- v

		case gmcom.CFailure:
			lw.Info("FAILURE received")
			ok := gm.memberMap.Delete(v.TargetID)
			if !ok {
				e := fmt.Errorf("failed to delete failed process-id: %v", v.TargetID)
				lw.Error(e)
			}

		case gmcom.CPing:
			lw.Info("PING received from process-id: %v", v.SrcID)
			rm := gm.processPing(v)
			gm.chout <- rm

		case gmcom.CAck:
			lw.Info("ACK received from process-id: %v", v.SrcID)
			// lw.Debug("ACK received")
			// lw.Debug("v.SrcID: %v", v.SrcID)
			// lw.Debug("v.SrcIPAddress: %v", v.SrcIPAddress)
			// lw.Debug("v.SrcStatus: %v", v.SrcStatus)
			// lw.Debug("v.SrcIncarnationNumber: %v", v.SrcIncarnationNumber)
			// lw.Debug("v.TargetID: %v", v.TargetID)
			// lw.Debug("v.TargetIPAddress: %v", v.TargetIPAddress)
			// lw.Debug("v.TargetStatus: %v", v.TargetStatus)
			// lw.Debug("v.TargetIncarnationNumber: %v", v.TargetIncarnationNumber)
			if v.SrcID != 0 {
				am := gmcom.GMMember{
					ID:                v.SrcID,
					IPAddress:         v.SrcIPAddress,
					Status:            v.SrcStatus,
					StatusCount:       v.SrcStatusCount,
					IncarnationNumber: v.SrcIncarnationNumber,
				}
				gm.memberMap.Add(v.SrcID, am)
				// lw.Debug("post ACK memberMap: %v", gm.memberMap)
			}
			gm.chout <- v

		case gmcom.CinNoAck:
			lw.Info("NoACK received for process-id: %v", v.TargetID)
			// lw.Debug("No ACK received for PING")
			// lw.Debug("v.SrcID: %v", v.SrcID)
			// lw.Debug("v.SrcIPAddress: %v", v.SrcIPAddress)
			// lw.Debug("v.TargetID: %v", v.TargetID)
			// lw.Debug("v.TargetIPAddress: %v", v.TargetIPAddress)
			// lw.Debug("v.TargetStatus: %v", v.TargetStatus)
			// lw.Debug("v.TargetIncarnationNumber: %v", v.TargetIncarnationNumber)

			// update the local member list
			tm := gmcom.GMMember{
				ID:                v.TargetID,
				IPAddress:         v.TargetIPAddress,
				Status:            v.TargetStatus,
				StatusCount:       v.TargetStatusCount,
				IncarnationNumber: v.TargetIncarnationNumber,
			}

			if tm.IncarnationNumber == 0 {
				tm.IncarnationNumber = 1
			}
			gm.memberMap.Add(v.TargetID, tm)

			// lw.Debug("post-SUSPECT gm.memberMap update for process-id: %v", v.TargetID)
			// for _, z := range gm.memberMap.IDMap {
			// 	lw.Debug("local membermap post-NoAck processing: %v %v %v", z.ID, z.Status, z.StatusCount)
			// }
			gm.chout <- v

		case gmcom.CinFlushMemberMap:
			ok := gm.memberMap.Flush()
			if !ok {
				lw.Error(fmt.Errorf("error: Requested gm.memberMap.Flush() did not succeed"))
				// panic("Requested gm.memberMap.Flush() did not succeed")
				v.ErrCode = gmcom.CerrMemberMapFlushFailed
			}
			v.ErrCode = ""
			gm.chout <- v

		case gmcom.CinDoSendPrep:

			// first increment failure status counts and flush failed processes that have been
			// disseminated ((gm.memberMap.Count()*2)-1) times.  this is overkill, as 2n-1 is
			// gives two rounds of pinging but it doesn't cost much to do it a few extra times.
			// failed and departed processes are updated or deleted from the local membermap here.
			err := gm.memberMap.UpdFailedDepartedProcesses()
			if err != nil {
				// TODO: log
				lw.ErrorWithPrefixString("error: gm.memberMap.UpdFailedDepartedProcesses() got:", err)
			}

			// populate the dissemination membermap with a deep-copy of the local membermap in order
			// to avoid goroutine contention during the encoding of the outgoing message.
			v.MemberMap = new(gmcom.OMap)
			v.MemberMap.Init()
			for k, mb := range gm.memberMap.IDMap {
				v.MemberMap.Add(k, mb)
			}

			// don't include failed processes in the ping list, (but include them in the membermap)
			c := gm.memberMap.Count()
			for i := 0; i < c; i++ {
				m, _ := gm.memberMap.ReadByIndex(i)
				if m.Status != gmcom.CStatusFailed {
					v.SrcID = gm.MyID
					v.SrcIPAddress = gm.MyIPAddress
					v.SrcIncarnationNumber = gm.MyIncarnation
					v.TargetID = m.ID
					v.TargetIPAddress = m.IPAddress
					v.TargetStatus = m.Status
					v.TargetStatusCount = m.StatusCount
					v.TargetIncarnationNumber = m.IncarnationNumber
					gm.chout <- v
				}
			}
			// don't want to close the channel, so send an invalid target-id for transmission in order to
			// inform the reader that this is the last msg (for now).
			v.TargetID = 0
			v.TargetIPAddress = ""
			gm.chout <- v

		case gmcom.CinGetMySrcInfo:
			v.SrcID = gm.MyID
			v.SrcIPAddress = gm.MyIPAddress
			v.SrcIncarnationNumber = gm.MyIncarnation
			gm.chout <- v

		case gmcom.CinGetMyTargetInfo:
			v.TargetID = gm.MyID
			v.TargetIPAddress = gm.MyIPAddress
			v.TargetIncarnationNumber = gm.MyIncarnation
			gm.chout <- v

		case gmcom.CinGetDBLeader:
			leader, err := gm.LSG.GetDBLeader()
			if err != nil {
				// if logging is enabled drop a message to stderr.
				// TODO_LOG
				// send unchanged message back to caller
				gm.chout <- v
			} else {
				rm := v
				rm.TargetID = leader.LeaderID
				rm.TargetIPAddress = leader.LeaderIPAddress
				gm.chout <- rm
			}

		case gmcom.CinSetDBLeader:

			gm.Leader.LeaderID = v.TargetID
			gm.Leader.LeaderIPAddress = v.TargetIPAddress

			err := gm.LSG.SetDBLeader(gm.Leader)
			if err != nil {
				lw.ErrorWithPrefixString("CinSetDBLeader failed:", err)
			}
			lw.Info("db and local leader set to %d, %s", v.TargetID, v.TargetIPAddress)
			gm.chout <- v

		case gmcom.CinSetLeader:
			gm.Leader.LeaderID = v.TargetID
			gm.Leader.LeaderIPAddress = v.TargetIPAddress
			gm.InElection = v.InElection // should be false
			lw.Info("local leader set to %d, %s", v.TargetID, v.TargetIPAddress)
			gm.chout <- v

		case gmcom.CinGetLocalLeader:
			v.TargetID = gm.Leader.LeaderID
			v.TargetIPAddress = gm.Leader.LeaderIPAddress
			gm.chout <- v

		case gmcom.CinStartElection:
			lw.Info("ELECTION: in gmcom.CinStartElection")
			gm.chout <- v

		case gmcom.CinRunElection:
			// set the InElection flag in the GMServ struct
			if gm.InElection == true {
				lw.Debug("ELECTION: CinRunElection detected active election and will exit")
				v.InElection = true
				v.ErrCode = gmcom.CElectComplete // election is running, so don't continue
				gm.chout <- v
				continue
			}
			gm.InElection = true

			// read the details of the largest non-failed ID from the local memberMap
			topMember := gm.memberMap.ReadMaxActiveID()
			if topMember == nil { // || (topMember.ID == gm.MyID && topMember.IPAddress == gm.MyIPAddress) {

				// local process will become the leader
				lw.Info("COORDINATOR: no topMember found in memberMap.  Adding local process-id to local memberMap.")
				gm.MyIncarnation++
				gm.memberMap.Add(gm.MyID, gmcom.GMMember{
					ID:                gm.MyID,
					IPAddress:         gm.MyIPAddress,
					Status:            gmcom.CStatusAlive,
					StatusCount:       1,
					IncarnationNumber: gm.MyIncarnation,
				})

				lw.Info("COORDINATOR: local process setting local process-id %d as leader in persistent store.", gm.MyID)
				err := gm.LSG.SetDBLeader(gmcom.GMLeader{
					LeaderID:        gm.MyID,
					LeaderIPAddress: gm.MyIPAddress,
				})
				if err != nil {
					panic("failed calling gm.LSG.SetDBLeader() in gmcom.CinRunElection")
				}
				gm.Leader.LeaderID = gm.MyID
				gm.Leader.LeaderIPAddress = gm.MyIPAddress
				gm.InElection = false
				v.ErrCode = gmcom.CElectComplete
				gm.chout <- v
				continue
			}
			lw.Info("COORDINATOR: found topMember: %d", topMember.ID)

			// if the local process has the highest active ID in the local memberMap,
			// make an attempt to become the leader
			if topMember.ID == gm.MyID {
				lw.Debug("COORDINATOR: topMember.ID == gm.MyID == %d", gm.MyID)
				lw.Info("COORDINATOR: local process setting local process-id %d as leader in persistent store", gm.MyID)
				lw.Info("COORDINATOR: local process %d will send coordinator message to group", gm.MyID)
				gm.MyIncarnation++
				gm.memberMap.Add(gm.MyID, gmcom.GMMember{
					ID:                gm.MyID,
					IPAddress:         gm.MyIPAddress,
					Status:            gmcom.CStatusAlive,
					StatusCount:       1,
					IncarnationNumber: gm.MyIncarnation,
				})

				lw.Info("COORDINATOR: local process setting local process-id %d as leader in persistent store", gm.MyID)
				err := gm.LSG.SetDBLeader(gmcom.GMLeader{
					LeaderID:        gm.MyID,
					LeaderIPAddress: gm.MyIPAddress,
				})
				if err != nil {
					panic("failed calling gm.LSG.SetDBLeader() in gmcom.CinRunElection")
				}

				gm.Leader.LeaderID = gm.MyID
				gm.Leader.LeaderIPAddress = gm.MyIPAddress
				gm.InElection = false

				lw.Debug("COORDINATOR: sending coordinator messages to group members")
				c := gmcom.GMMessage{
					Type:                 gmcom.CCoordinator,
					SrcID:                gm.MyID,
					SrcIPAddress:         gm.MyIPAddress,
					SrcIncarnationNumber: gm.MyIncarnation,
					SubjectID:            gm.MyID,
					SubjectIPAddress:     gm.MyIPAddress,
				}
				// send coordinator messages to all non-failed processes with lower-ids than that of the local process.
				for i := 0; i < gm.memberMap.Count(); i++ {
					m, ok := gm.memberMap.ReadByIndex(i)
					if ok && m.Status != gmcom.CStatusFailed && m.Status != gmcom.CStatusDeparted && m.ID != gm.MyID {
						c.TargetID = m.ID
						c.TargetIPAddress = m.IPAddress
						err = gm.SendCoordinator(c)
						if err != nil {
							lw.Warning("COORDINATOR: send of coordinator message to process-id %d with %d as leader failed.  got: %s", c.TargetID, c.SubjectID, err.Error())
						} else {
							lw.Info("COORDINATOR: send coordinator message to process-id %d with %d as leader", c.TargetID, c.SubjectID)
						}
					}
				}
				v.ErrCode = gmcom.CElectComplete
				gm.chout <- v
				lw.Info("COORDINATOR: COMPLETE - set leader-id: %d", topMember.ID)
				continue
			}

			if topMember.ID > gm.MyID {
				// do nothing here.
				gm.chout <- v
			}
			// // local process-id is lower than the highest active ID in the local memberMap,
			// // so trigger an election by sending Election messages to all higher non-failed
			// // process-id's.
			// if topMember.ID > gm.MyID {
			// 	log.Printf("ELECTION: local process-id %d is lower than topMember.ID %d\n", gm.MyID, topMember.ID)
			// 	log.Debug("ELECTION: triggering election via Election message dispatch")
			// 	e := gmcom.GMMessage{
			// 		Type:                 gmcom.CElection,
			// 		SrcID:                gm.MyID,
			// 		SrcIPAddress:         gm.MyIPAddress,
			// 		SrcIncarnationNumber: gm.MyIncarnation,
			// 		SubjectID:            gm.MyID,
			// 		SubjectIPAddress:     gm.MyIPAddress,
			// 	}
			// 	// send election messages to all non-failed processes with higher-ids than that of the local process.
			// 	// there are three possibilities here; i) OK ii) Ack iii) no response
			// 	for i := 0; i < gm.memberMap.Count(); i++ {
			// 		m, ok := gm.memberMap.ReadByIndex(i)
			// 		if ok && m.Status != gmcom.CStatusFailed && m.ID > gm.MyID {
			// 			e.TargetID = m.ID
			// 			e.TargetIPAddress = m.IPAddress
			// 			err := gm.SendElection(&e) // SendElection should return COK, CAck, CNoACK?
			// 			if err != nil {
			// 				log.Printf("ELECTION: send of election message to process-id %d with %d as proposed leader failed.  got: %s", e.TargetID, e.SubjectID, err)
			// 				continue
			// 			}
			// 			log.Printf("ELECTION: sent election message to process-id %d with %d as proposed leader\n", e.TargetID, e.SubjectID)
			// 			log.Printf("ELECTION: got response code of %s from process-id %d in response to propsal of process-id %d as leader\n", e.ErrCode, e.TargetID, e.SubjectID)
			// 			switch e.ErrCode {
			// 			case gmcom.CElectOK:
			// 				continue

			// 			case gmcom.CElectNR:
			// 				topMember = gm.memberMap.ReadMaxActiveID()
			// 				if topMember != nil && topMember.ID == gm.MyID {
			// 					// set local process-id as leader in persistent store
			// 					// send coordinator messages
			// 					// unset election flag
			// 					log.Printf("ELECTION: local process setting local process-id %d as leader in persistent store\n", gm.MyID)
			// 					err := gm.LSG.SetDBLeader(gmcom.GMLeader{
			// 						LeaderID:        gm.MyID,
			// 						LeaderIPAddress: gm.MyIPAddress,
			// 					})
			// 					if err != nil {
			// 						panic("failed calling gm.LSG.SetDBLeader() in gmcom.CinRunElection")
			// 					}
			// 					gm.Leader.LeaderID = gm.MyID
			// 					gm.Leader.LeaderIPAddress = gm.MyIPAddress
			// 					gm.InElection = false

			// 					log.Debug("ELECTION: sending coordinator messages to group members")
			// 					c := gmcom.GMMessage{
			// 						Type:                 gmcom.CCoordinator,
			// 						SrcID:                gm.MyID,
			// 						SrcIPAddress:         gm.MyIPAddress,
			// 						SrcIncarnationNumber: gm.MyIncarnation,
			// 						SubjectID:            gm.MyID,
			// 						SubjectIPAddress:     gm.MyIPAddress,
			// 					}
			// 					// send coordinator messages to all non-failed processes with lower-ids than that of the local process.
			// 					for i := 0; i < gm.memberMap.Count(); i++ {
			// 						m, ok := gm.memberMap.ReadByIndex(i)
			// 						if ok && m.Status != gmcom.CStatusFailed && m.ID != gm.MyID {
			// 							c.TargetID = m.ID
			// 							c.TargetIPAddress = m.IPAddress
			// 							err = gm.SendCoordinator(c)
			// 							if err != nil {
			// 								log.Printf("ELECTION: send of coordinator message to process-id %d with %d as leader failed.  got: %s", c.TargetID, c.SubjectID, err)
			// 							} else {
			// 								log.Printf("ELECTION: send coordinator message to process-id %d with %d as leader\n", c.TargetID, c.SubjectID)
			// 							}
			// 						}
			// 					}
			// 					v.ErrCode = gmcom.CElectComplete
			// 					gm.chout <- v
			// 				} else {
			// 					// restart election via recursive call?
			// 					v.ErrCode = gmcom.CElectIncomplete
			// 					gm.chout <- v
			// 				}
			// 			default:

			// 			}
			// 		}
			// 	}
			// 	gm.chout <- v
			// }

		case gmcom.CinSetElectionState:
			gm.InElection = v.InElection

		case gmcom.CinGetElectionState:
			v.InElection = gm.InElection
			gm.chout <- v

		case gmcom.CinGetLocalDetails:
			v.TargetID = gm.MyID
			v.TargetIPAddress = gm.MyIPAddress

			// deep-copy / race mitigation
			v.MemberMap = new(gmcom.OMap)
			v.MemberMap.Init()
			for k, mb := range gm.memberMap.IDMap {
				v.MemberMap.Add(k, mb)
			}
			gm.chout <- v

		default:
			// do nothing
		}
	}
}

// Serve starts the group membership server
func (gm *GMServ) Serve(myIPAddress string, lsg gmcom.GMLeaderSetterGetter, actUsrs *gmcom.ActUsrsH, groupAuths *gmcom.GroupAuthsH, auths *gmcom.AuthsH, usrGroups *gmcom.UsrGroupsH, logging bool, evlCycle uint, failureThreshold uint64) {
	ft := uint32(failureThreshold) // 64-bit atomic alignment mitigation for 32-bit ARM
	err := gm.initialize(myIPAddress, lsg, actUsrs, groupAuths, auths, usrGroups, logging, ft)
	if err != nil {
		panic("Serve()" + err.Error())
	}

	mux := http.NewServeMux()

	// failure detector handlers
	mux.Handle("/ping", websocket.Handler(gm.PingHandler))
	mux.Handle("/join", websocket.Handler(gm.JoinHandler))
	mux.Handle("/departing", websocket.Handler(gm.DepartingHandler))
	mux.Handle("/coordinator", websocket.Handler(gm.CoordinatorHandler))

	mux.Handle("/updateusrcache", websocket.Handler(gm.UsrUpdateHandler))
	mux.Handle("/updategroupauthcache", websocket.Handler(gm.GroupAuthUpdateHandler))
	mux.Handle("/updateauthcache", websocket.Handler(gm.AuthUpdateHandler))
	mux.Handle("/updateusrgroupcache", websocket.Handler(gm.UsrGroupUpdateHandler))
	mux.Handle("/set", websocket.Handler(gm.SetHandler))

	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Done()

	lw.Console("Listening for ws traffic on: %s", gm.MyIPAddress)

	// create the ws server - this can be stopped by calling cs.httpServer.Shutdown(...)
	gm.HTTPServer = &http.Server{
		Addr:    myIPAddress,
		Handler: mux,
	}

	go func() {
		err := gm.HTTPServer.ListenAndServe()
		if err != nil {
			// this will write something in a Shutdown scenario too
			lw.ErrorWithPrefixString("cs.httpServer.ListenAndServe() error - got:", err)
		}
	}()

	// exagerated wait for httpserver spinup
	time.Sleep(2 * time.Second)

	// initialization
	// 1.  Get Leader from DB
	// 2.    succecss == continue
	// 3.    fail     == I am leader?
	// 4.  Contact leader (Ping without ID)
	// 5.    success  == Ack with ID
	// 5b.   success  == Leader sends Join message with current membership list
	// 6.    fail     == back, retry (x2)
	// 7.    fail     == I am leader?
	lastJoinErr := &gmcom.JoinError{}
	lm := gmcom.GMMessage{}
	rand.Seed(time.Now().UnixNano())

	// initialization loops
	// try to join 4 times - here there be dragons
	for z := 0; z < 4; z++ {

		// cleanup last error
		lastJoinErr.Clear()

		// try to read the leader from the persistent store 3 times per join attempt
		for i := 1; i < 4; i++ {

			// 1. get leader from db - 3 attempts
			gm.SendGetDBLeader(&lm)
			if lm.TargetID == 0 || lm.TargetIPAddress == "" {
				lw.Console("get db leader attempt: %d", i)
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				continue
			}
			break
		}

		// is my ip address the same as that of the leader that was read from the persistent
		// store?  this can happen in scenarios where the leader is the last process standing
		// and then terminates without warning.  the persistent store retains the last leader
		// information (Pi) resulting in a successful leader-ping if the first process to restart
		// in the group happens to hold the same ip-address as that of the old leader.  a successful
		// ping (Pi->Pi) looks like there is a functioning leader, so if the joining process sees
		// it's own ip-address in the leader slot within the persistent store, it is best to act
		// like there is no current leader and wipe the read (leader) information clean.
		// i.e. {inGetDBLeader 0 192.168.1.42:5050 0  0 1 192.168.1.42:5050 0   0 0  true <nil>}
		// see if the persisted leader is contactable or not...
		if lm.TargetID != 0 && lm.TargetIPAddress != "" && lm.TargetIPAddress != myIPAddress {
			lw.Info("JOIN: leader read from persistent store: %v", lm)
			gm.SendSetLeader(&lm)

			// try to join (really just a ping)
			l := gmcom.GMLeader{
				LeaderID:        lm.TargetID,
				LeaderIPAddress: lm.TargetIPAddress,
			}
			// need a better way of determining how the join failed - need to be certain that it is
			// okay to attempt to seize leadership here.
			// 1. could not contact leader (i.e. leader in persistent store is incorrect/in transition)
			// 2. leader is okay, but the addresss is in use by an active or suspect(most likely) process-id
			err = gm.SendJoin(&l)
			if err, ok := err.(*gmcom.JoinError); ok {
				lw.Warning("SendJoin got errorCode: %v", err.ErrCode)
				lw.Error(err)
				lastJoinErr = err
				lw.Warning("JOIN: failed to join with leader process-id %d at %s - error code: %s", lm.TargetID, lm.TargetIPAddress, lastJoinErr.ErrorCode())
				switch lastJoinErr.ErrorCode() {
				case gmcom.CerrJoinNoContact, gmcom.CerrJoinNotLeader:
					time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
					continue

				case gmcom.CerrJoinAddrInUse:
					time.Sleep(5 * time.Second)
					continue

				case gmcom.CerrJoinUnknownErr:
					time.Sleep(5 * time.Second)
					continue

				default:

				}
			}
			break
		}

		// if the local process-id was read from the persistent store, wait a bit and check again
		// before deciding to become the real leader.
		if lm.TargetID != 0 && lm.TargetIPAddress == myIPAddress {
			lastJoinErr.ErrCode = gmcom.CerrJoinNoContact
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		// if no leader-id was read from the db, try to become the leader
		if lm.TargetID == 0 && lm.TargetIPAddress == "" {
			lastJoinErr.ErrCode = gmcom.CerrJoinNoContact
			break
		}
	}

	switch lastJoinErr.ErrorCode() {
	case gmcom.CerrJoinNoContact:
		lw.Warning("JOIN: no leader detected - setting myID = 1 as leader")
		gm.MyID = 1
		lm.TargetID = gm.MyID
		lm.TargetIPAddress = gm.MyIPAddress
		lm.Type = gmcom.CinSetDBLeader
		lm.ExpectResponse = true

		// add local server to the memberMap
		mb := gmcom.GMMember{
			ID:                gm.MyID,
			IPAddress:         gm.MyIPAddress,
			Status:            gmcom.CStatusAlive,
			StatusCount:       1,
			IncarnationNumber: gm.MyIncarnation,
		}
		gm.memberMap.Add(gm.MyID, mb)

		// set gm.MyID as leader locally and in the persistent store
		gm.SendSetDBLeader(&lm)

	case gmcom.CerrJoinNotLeader:
		lw.Error(lastJoinErr) //.ErrorSummary())
		os.Exit(-1)

	case gmcom.CerrJoinAddrInUse:
		lw.Error(lastJoinErr) //.ErrorSummary())
		os.Exit(-1)

	case gmcom.CerrJoinUnknownErr:
		lw.Error(lastJoinErr) //.ErrorSummary())
		os.Exit(-1)

	default:

	}

	// ensure there is a pause in the failure detector evl
	if evlCycle == 0 {
		evlCycle = 5
	}

	// start failure detector event loop
	go func() {
		evl := true
		for evl {
			time.Sleep(time.Duration(evlCycle) * time.Second)
			lw.Debug("PING CYCLE: starting ping cycle...")
			gm.SendPing()
			lw.Debug("PING CYCLE: ping cycle complete.")
		}
	}()

	lw.Console("initialization complete...")
	lw.Console("leader info is: %v", gm.Leader)
	wg.Wait()
}

// GetDBLeaderHandler is used to access the persisted leader information, which
// may be stored in any accessible medium.
func (gm *GMServ) getDBLeader(ws *websocket.Conn) {

	// call function
	leader, err := gm.LSG.GetDBLeader()
	if err != nil {
		lw.ErrorWithPrefixString("getDBLeader() error - got:", err)
		return
	}
	if leader != nil {
		lw.Info("getDBLeader() got leader: %v", leader)
	}
}

// SetDBLeaderHandler is used to set the persisted leader information, which
// may be stored in any accessible medium.  It is up to the developer to decide
// what to put in here (accept closure etc.?)
func (gm *GMServ) SetDBLeaderHandler(ws *websocket.Conn) {

	// call function
}

// PingHandler deals with incoming failure detector ping reqests.  Ping messages
// are decoded and then fed into the inbound GMMessage channel for processing.
// Message Type == Ping
func (gm *GMServ) PingHandler(ws *websocket.Conn) {

	// decode the incoming PING message and feed to inbound channel
	m, err := gmcom.DecodeGMMessage(ws)
	if err == nil {

		// process the ping and wait for the outbound channel to provide a ACK message
		gm.chin <- *m
		m2 := <-gm.chout
		lw.Info("PingHandler() channel result: %v", m2)

		// the inclusion of the MemberMap in m2 is problematic as the reference to the
		// map is common with gm.  IDMap/MemberMap access is synchronized via the
		// serialization channels (gm.chin/gm.chout), but calling EncodeGMMessage(m2)
		// forces an unsynchronized read of the m2.MemeberMap to occur. If conditions
		// are right, this will result in a data-race, notably when dealing with the
		// disappearance of (Pj) where (Pj) is pinging (Pi). The value of the outgoing
		// MemberMap is negligable, as dissemination occurs via the PING message only,
		// so for now m2.MemberMap is set to nil in all response messages.
		// If you need it for some reason, the best thing to do would be to create a
		// copy of the map in the upstream PingHandler processor.
		m2.MemberMap = nil

		// check for incorrect receiver problem
		if m2.ErrCode == gmcom.CerrPingIncorrectReceiver {
			return // do not send ACK
		}

		// gob-encode the ACK message for transmission
		rawM2, err := gmcom.EncodeGMMessage(m2)
		if err != nil {
			// timeout the ws connection and drop a log message to stderr
			lw.ErrorWithPrefixString("PingHandler() gobEncodeError:", err)
		} else {
			_, err = ws.Write(rawM2)
			if err != nil {
				lw.ErrorWithPrefixString("PingHandler() ws.Write(rawM2) got error:", err)
				// did the sender disappear?  if so, this is suspicious,
				// but the sender's pinger should pickup on the fact that
				// the process is not responding.  if logging is enabled
				// drop a message to stderr.
			}
		}

		// is there a need to trigger an elelection?
		if !m2.CheckElection {
			return
		}

		lw.Debug("ELECTION: PingHandler() starting election evaluation...")
		lw.Debug("ELECTION: unsafe access to gm.InElection returns: %v", gm.InElection)

		// maybe do all of the following in new CinRunElection?
		s := gmcom.GMMessage{
			Type: gmcom.CinRunElection,
		}
		// for s.ErrCode != gmcom.CElectComplete {
		gm.chin <- s
		s = <-gm.chout
		//	lw.Debug("ELECTION: election not completed.  Waiting 500ms for restart...")
		//	time.Sleep(500 * time.Millisecond)
		// }

		lw.Debug("PingHandler() complete")
		return
	}
}

// CoordinatorHandler deals with incoming coordinator messages sent by processes
// stating that they are the new leader.
func (gm *GMServ) CoordinatorHandler(ws *websocket.Conn) {
	lw.Debug("COORDINATOR: in gm.CoordinatorHandler")

	// decode the incoming COORDINATOR message and feed to inbound channel
	c, err := gmcom.DecodeGMMessage(ws)
	if err != nil {
		lw.Error(fmt.Errorf("COORDINATOR: received /coordinator call with unreadable payload"))
		return
	}

	// check that the incoming message-type is correct
	if c.Type != gmcom.CCoordinator {
		lw.Error(fmt.Errorf("COORDINATOR: received /coordinator call with incorrect message-type.  got:", c.Type))
		return
	}

	// set the new leader info in the GMServ struct and clear the InElection status
	l := gmcom.GMMessage{
		Type:            gmcom.CinSetLeader,
		TargetID:        c.SubjectID,
		TargetIPAddress: c.SubjectIPAddress,
		InElection:      false,
	}
	gm.chin <- l
	l = <-gm.chout
}

// DepartingHandler deals with incoming departing messages.  Departing messages
// are decoded and then fed into the inbound GMMessage channel for processing.
// Message Type == Departing
func (gm *GMServ) DepartingHandler(ws *websocket.Conn) {

	// decode the incoming DEPARTING message and feed to inbound channel
	m, err := gmcom.DecodeGMMessage(ws)
	if err == nil {

		lw.Debug("DEPARTING: In DepartingHandler for: %v", m)

		// process the departing message and wait for the outbound channel to provide a response
		gm.chin <- *m
		m2 := <-gm.chout
		lw.Debug("DepartingHandler() channel result: %v", m2)

		// is there a need to trigger an elelection?
		if !m2.CheckElection {
			return
		}

		lw.Debug("ELECTION: DepartingHandler() starting election evaluation...")
		lw.Debug("ELECTION: unsafe access to gm.InElection returns: %v", gm.InElection)

		s := gmcom.GMMessage{
			Type: gmcom.CinRunElection,
		}

		gm.chin <- s
		s = <-gm.chout
		lw.Debug("DepartingHandler() complete")
		return
	}
}

// JoinHandler deals with join websocket requests that have been accepted by the router.
// It is implied that the caller is attempting to join the group by contacting the leader.
// What to do if the local service is no longer the leader?
func (gm *GMServ) JoinHandler(ws *websocket.Conn) {

	// decode the incoming JOIN message and feed to inbound channel
	m, err := gmcom.DecodeGMMessage(ws)
	if err == nil {
		gm.chin <- *m

		// wait for the outbound channel to provide a ACK message
		m2 := <-gm.chout

		// gob-encode the ACK message for transmission
		lw.Debug("JoinHandler() response GMMessage: %v", m2)
		rawM2, err := gmcom.EncodeGMMessage(m2)
		if err != nil {
			// timeout the ws connection and drop a log message to stderr
			lw.ErrorWithPrefixString("JoinHandler():", err)
		} else {
			_, err = ws.Write(rawM2)
			if err != nil {
				// did the sender disappear?  if so, this is suspicious,
				// but the sender's pinger should pickup on the fact that
				// the process is not responding.  if logging is enabled
				// drop a message to stderr.
				lw.ErrorWithPrefixString("JoinHandler() ws.Write response error:", err)
			}
		}
	} else {
		lw.ErrorWithPrefixString("JoinHandler error:", err)
	}
}

// UsrUpdateHandler handles incoming traffic from other group-members containing information regarding
// changes to the status of application Usr entities.
func (gm *GMServ) UsrUpdateHandler(ws *websocket.Conn) {
	lw.Debug("In UsrUpdateHandler()")

	// gob decoding
	var u gmcom.ActUsrD
	var msg = make([]byte, 1024)
	l, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("UsrUpdateHandler() ws.Read() error:", err)
		return
	}
	raw := msg[0:l]
	decBuf := bytes.NewBuffer(raw)
	err = gob.NewDecoder(decBuf).Decode(&u)
	if err != nil {
		lw.ErrorWithPrefixString("UsrUpdateHandler() gob.Decode() error:", err)
		return
	}

	// update the local server's ActiveUsrs cache (map)
	gm.ActUsrsH.Lock()
	gm.ActUsrsH.ActiveUsrs[u.ID] = u.Active
	gm.ActUsrsH.Unlock()

	// send to other group members?
	if !u.Forward {
		ws.Write([]byte("true"))
		return
	}
	u.Forward = false

	// send the update to all non-failed processes in the process group
	// get a list of the active processes (this is inherently stale)
	m := gm.SendGetLocalDetails()
	if m == nil {
		lw.Warning("Failed to read usr cache group server details in SendGetLocalDetails()")
		ws.Write([]byte("false"))
	}

	// send process-id/article-key to group members
	r := m.MemberMap.ReadActiveProcessList()
	for _, g := range r {
		if g.ID == gm.MyID {
			continue
		}
		lw.Info("CS: SENDING %v to %s", u, g.IPAddress)
		err := gmcl.AddUpdUsrCache(u, g.IPAddress) //TODO go()
		if err != nil {
			lw.ErrorWithPrefixString("wscl.AddUpdUsrCache() error:", err)
		}
	}
	ws.Write([]byte("true"))
}

// GroupAuthUpdateHandler handles incoming traffic from other group-members containing information regarding
// changes to GroupAuth entities.
func (gm *GMServ) GroupAuthUpdateHandler(ws *websocket.Conn) {
	lw.Debug("In GroupAuthUpdateHandler()")

	// gob decoding
	var ga gmcom.GroupAuthD
	var msg = make([]byte, 1024)
	l, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("GroupAuthUpdateHandler() ws.Read() error:", err)
		return
	}
	raw := msg[0:l]
	decBuf := bytes.NewBuffer(raw)
	err = gob.NewDecoder(decBuf).Decode(&ga)
	if err != nil {
		lw.ErrorWithPrefixString("GroupAuthUpdateHandler() gob.Decode() error:", err)
		return
	}

	// lookup the groupName and authName
	ga.GroupName = gm.UsrGroupsH.GroupNames[ga.GroupID]
	ga.AuthName = gm.AuthsH.Auths[ga.AuthID]
	if ga.GroupName == "" || ga.AuthName == "" {
		lw.Warning("GroupAuthUpdateHandler() could not update with GroupName: %s and AuthName: %s", ga.GroupName, ga.AuthName)
		return
	}

	// update the local server's GroupAuths cache (map)
	switch ga.Op {
	case gmcom.COpCreate, gmcom.COpUpdate:
		// do the local cache update
		gm.GroupAuthsH.Lock()
		mapAuth := gm.GroupAuthsH.GroupAuths[ga.GroupName]
		if mapAuth == nil {
			mapAuth = make(map[string]bool)
			mapAuth[ga.AuthName] = true
			gm.GroupAuthsH.GroupAuths[ga.GroupName] = mapAuth
		} else {
			// if the groupName does exist in the top-level map, add the auth to
			// the group's auth map.
			mapAuth[ga.AuthName] = true
		}

		gm.GroupAuthsH.GroupAuthsID[ga.ID] = gmcom.GroupAuthNames{GroupName: ga.GroupName, AuthName: ga.AuthName}
		gm.GroupAuthsH.Unlock()

	case gmcom.COpDelete:
		gm.GroupAuthsH.Lock()
		mapAuth := gm.GroupAuthsH.GroupAuths[ga.GroupName]
		if mapAuth != nil {
			delete(mapAuth, ga.AuthName)
		}
		delete(gm.GroupAuthsH.GroupAuthsID, ga.ID)
		gm.GroupAuthsH.Unlock()

	default:
		// do nothing
	}

	// send to other group members?
	if !ga.Forward {
		ws.Write([]byte("true"))
		return
	}
	ga.Forward = false

	// send the update to all non-failed processes in the process group
	// get a list of the active processes (this is inherently stale)
	m := gm.SendGetLocalDetails()
	if m == nil {
		lw.Warning("Failed to read groupauth cache group server details in SendGetLocalDetails()")
		ws.Write([]byte("false"))
	}

	// send process-id/article-key to group members
	r := m.MemberMap.ReadActiveProcessList()
	for _, g := range r {
		if g.ID == gm.MyID {
			continue
		}
		lw.Info("CS: SENDING %v to %s", ga, g.IPAddress)
		err := gmcl.AddUpdGroupAuthCache(ga, g.IPAddress) //TODO go()
		if err != nil {
			lw.ErrorWithPrefixString("wscl.AddUpdGroupAuthCache() error:", err)
		}
	}
	ws.Write([]byte("true"))
}

// AuthUpdateHandler handles incoming traffic from other group-members containing information regarding
// changes to the status of application Auth entities.
func (gm *GMServ) AuthUpdateHandler(ws *websocket.Conn) {
	lw.Debug("In AuthUpdateHandler()")

	// gob decoding
	var a gmcom.AuthD
	var msg = make([]byte, 1024)
	l, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("AuthUpdateHandler() ws.Read() error:", err)
		return
	}
	raw := msg[0:l]
	decBuf := bytes.NewBuffer(raw)
	err = gob.NewDecoder(decBuf).Decode(&a)
	if err != nil {
		lw.ErrorWithPrefixString("AuthUpdateHandler() gob.Decode() error:", err)
		return
	}

	// update the local server's ActiveUsrs cache (map)
	switch a.Op {
	case gmcom.COpCreate, gmcom.COpUpdate:
		gm.AuthsH.Lock()
		gm.AuthsH.Auths[a.ID] = a.AuthName
		gm.AuthsH.Unlock()
	case gmcom.COpDelete:
		gm.AuthsH.Lock()
		delete(gm.AuthsH.Auths, a.ID)
		gm.AuthsH.Unlock()
	default:
		// do nothing
	}

	// send to other group members?
	if !a.Forward {
		ws.Write([]byte("true"))
		return
	}
	a.Forward = false

	// send the update to all non-failed processes in the process group
	// get a list of the active processes (this is inherently stale)
	m := gm.SendGetLocalDetails()
	if m == nil {
		lw.Warning("Failed to read auth cache group server details in SendGetLocalDetails()")
		ws.Write([]byte("false"))
	}

	// send process-id/article-key to group members
	r := m.MemberMap.ReadActiveProcessList()
	for _, g := range r {
		if g.ID == gm.MyID {
			continue
		}
		lw.Info("CS: SENDING %v to %s", a, g.IPAddress)
		err := gmcl.AddUpdAuthCache(a, g.IPAddress) //TODO go()
		if err != nil {
			lw.ErrorWithPrefixString("wscl.AddUpdAuthCache() error:", err)
		}
	}
	ws.Write([]byte("true"))
}

// UsrGroupUpdateHandler handles incoming traffic from other group-members containing information regarding
// changes to the status of application UsrGroup entities.
func (gm *GMServ) UsrGroupUpdateHandler(ws *websocket.Conn) {
	lw.Debug("In UsrGroupUpdateHandler()")

	// gob decoding
	var ug gmcom.UsrGroupD
	var msg = make([]byte, 1024)
	l, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("UsrGroupUpdateHandler() ws.Read() error:", err)
		return
	}
	raw := msg[0:l]
	decBuf := bytes.NewBuffer(raw)
	err = gob.NewDecoder(decBuf).Decode(&ug)
	if err != nil {
		lw.ErrorWithPrefixString("UsrGroupUpdateHandler() gob.Decode() error:", err)
		return
	}

	// update the local server's UsrGroup cache (map)
	switch ug.Op {
	case gmcom.COpCreate, gmcom.COpUpdate:
		gm.UsrGroupsH.Lock()
		gm.UsrGroupsH.GroupNames[ug.ID] = ug.GroupName
		gm.UsrGroupsH.Unlock()
	case gmcom.COpDelete:
		gm.UsrGroupsH.Lock()
		delete(gm.UsrGroupsH.GroupNames, ug.ID)
		gm.UsrGroupsH.Unlock()
	default:
		// do nothing
	}

	// send to other group members?
	if !ug.Forward {
		ws.Write([]byte("true"))
		return
	}
	ug.Forward = false

	// send the update to all non-failed processes in the process group
	// get a list of the active processes (this is inherently stale)
	m := gm.SendGetLocalDetails()
	if m == nil {
		lw.Warning("Failed to read usrgroup cache group server details in SendGetLocalDetails()")
		ws.Write([]byte("false"))
	}

	// send process-id/article-key to group members
	r := m.MemberMap.ReadActiveProcessList()
	for _, g := range r {
		if g.ID == gm.MyID {
			continue
		}
		lw.Info("CS: SENDING %v to %s", ug, g.IPAddress)
		err := gmcl.AddUpdUsrGroupCache(ug, g.IPAddress) //TODO go()
		if err != nil {
			lw.ErrorWithPrefixString("wscl.AddUpdUsrGroupCache() error:", err)
		}
	}
	ws.Write([]byte("true"))
}

// LeaveHandler handles the announced departure of a process, thereby
// facilitating its graceful exit from the membership group.  the
// LeaveHandler must ensure that all processes are aware that the
// process in question (Pj) has left the group.
// Message Type == Leave
func (gm *GMServ) LeaveHandler(ws *websocket.Conn) {

	// gm.leaveHandler(...)
}

// SetHandler pushes the AU Article into the serialization channel
func (gm *GMServ) SetHandler(ws *websocket.Conn) {

	// decode the incoming Ping message and feed to the inbound
	// GMMessage channel for processing.
	m, err := gmcom.DecodeGMMessage(ws)
	if err == nil {
		gm.chin <- *m

		// wait the outbound channel to provide a Ack message
		// and then gob-encode it for transmission.
		// m2 := <-gm.chout

	}

	m2 := <-gm.chout
	rawM2, err := gmcom.EncodeGMMessage(m2)
	if err != nil {
		// timeout the ws connection
		time.Sleep(2 * time.Second)
	} else {
		_, err = ws.Write(rawM2)
		if err != nil {
			// did the sender disappear?
			// this is suspicious, but the sender's pinger should
			// pickup on the fact that they are not responding?  if
			// logging is enabled you may want to drop a message to
			// stderr.
		}
	}
}

// ==== internal handler code
// getPingTarget returns the salient attributes of process Pj
func (gm *GMServ) getPingTarget() (*gmcom.GMMember, bool) {

	m, ok := gm.memberMap.ReadCWNeighbour(gm.MyID)
	if !ok {
		return nil, ok
	}
	return &m, true
}
