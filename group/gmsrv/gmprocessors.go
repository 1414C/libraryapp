package gmsrv

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:57 CDT
//=============================================================================================

import (
	"fmt"

	"github.com/1414C/libraryapp/group/gmcom"

	"github.com/1414C/lw"
)

// All methods contained herein must be called within the processCmdChannel method.
// It would likely be better to forgo the use of channels and simply use the
// sync.RWMutex or maybe even the new sync.Map to mitigate data-races on the
// shared gm data.  channels are go doctrine for shared access, but the numbers
// (at least in this case) support the use of RWMutex instead.

// processPing evaluates the attached process list being disseminated
// by the pinger (SrcID).  An Ack message containing the appropriate
// content is returned to the called.  Should ACK messages contain
// (Pi)'s suspect/failure list?
func (gm *GMServ) processPing(p gmcom.GMMessage) (rm gmcom.GMMessage) {

	il := false
	// am I the leader?
	// sender has an id?
	// if I am the leader, grant the sender an id (current max(pList) + 1)
	// If I am not the leader, should I do it anyway?
	// evaluate the sender's process-list
	// make updates to my own process list
	// send an Ack to the sender
	// am I a supect?  If so, increment my incarnation number and continue to ping
	rm = gmcom.GMMessage{
		Type:                    gmcom.CAck,
		SrcID:                   gm.MyID,
		SrcIPAddress:            gm.MyIPAddress,
		SrcIncarnationNumber:    gm.MyIncarnation,
		TargetID:                p.SrcID,
		TargetIPAddress:         p.SrcIPAddress,
		TargetIncarnationNumber: p.SrcIncarnationNumber,
	}

	// am I the leader?
	if gm.Leader.LeaderID == gm.MyID {
		il = true
	}

	// if the sender does not have an id try to assign one
	if p.SrcID == 0 {
		m := gmcom.GMMember{
			ID:                p.SrcID,
			IPAddress:         p.SrcIPAddress,
			Status:            gmcom.CStatusAlive,
			StatusCount:       1,
			IncarnationNumber: p.SrcIncarnationNumber,
		}

		// if I am the leader, allocate an id for the sender and add them
		// to the (local) membership list.
		// 1. If I am not the leader, send an internal NoAck message?
		// 2. If the address exists in the leader's memberMap, and the status
		//    of that process is ACTIVE or SUSPECT, then we cannot add assign
		//    a new process-id for the address.
		if il {

			// 1
			da := gm.memberMap.ReadByAddress(p.SrcIPAddress)
			if da != nil {
				for _, v := range da {
					if v.Status == gmcom.CStatusAlive ||
						v.Status == gmcom.CStatusSuspect {
						e := fmt.Errorf("address %s is in use with proc-id: %d status: %s", v.IPAddress, v.ID, v.Status)
						lw.Error(e)
						rm.ErrCode = gmcom.CerrJoinAddrInUse
						return rm // TargetIPAddress == 0
					}
				}
			}

			// 2
			err := gm.memberMap.AddWithoutID(&m)
			if err != nil {
				lw.Error(err)
				rm.ErrCode = gmcom.CerrJoinUnknownErr
				return rm
			} else {
				rm.TargetID = m.ID
				rm.TargetIncarnationNumber = m.IncarnationNumber
				rm.TargetStatus = m.Status
				rm.ErrCode = ""
				rm.MemberMap = gm.memberMap
				return rm
			}
		} else {
			// remember that rm.TargetID == 0 here
			lw.Info("non-leader process-id %d received JOIN message from ip-address: %s", gm.MyID, p.SrcIPAddress)
			rm.ErrCode = gmcom.CerrJoinNotLeader
			return rm
		}

	} else {
		// this ping came from a numbered process, so we have to update
		// our local memberMap with status, status count and incarnation
		// number.

		// check that the ping was intended for this process-id
		if p.TargetID != gm.MyID {
			rm.ErrCode = gmcom.CerrPingIncorrectReceiver
			return rm
		}

		// read the sender process-id from the memberMap
		pj, ok := gm.memberMap.Read(p.SrcID)

		// if the pinging process-id is not in the local memberMap - add it, but do
		// not use its membership list to update our membership list.  skipping a
		// single update will not hurt.
		if !ok {
			m := gmcom.GMMember{
				ID:                p.SrcID,
				IPAddress:         p.SrcIPAddress,
				Status:            gmcom.CStatusAlive,
				StatusCount:       1,
				IncarnationNumber: p.SrcIncarnationNumber,
			}

			// add the process-id to the memberMap
			err := gm.memberMap.Add(p.SrcID, m)
			if err != nil {
				lw.Error(err)
			}

			// if the pinging process-id (pj.ID)is in the local memberMap - update it
		} else {
			switch pj.Status {
			case gmcom.CStatusAlive:
				pj.IncarnationNumber = p.SrcIncarnationNumber
				pj.StatusCount++

				// update the pinger's information in the local memberMap
				err := gm.memberMap.Add(pj.ID, pj)
				if err != nil {
					lw.Error(err)
				}

				// perform local updates based on incoming ping memberMap
				if p.MemberMap != nil {
					// BETA: is the local process-id (Pi) marked as suspect or failed in the sender's memberMap?
					err = gm.updateOwnStatusInPingMap(&p)
					if err != nil {
						lw.Error(err)
					}

					// update the local memberMap based on the incoming ping's memberMap
					rm.CheckElection, err = gm.memberMap.UpdFromPing(gm.MyID, pj.ID, p.MemberMap, uint(gm.FailureThreshold), gm.Leader.LeaderID)
					if err != nil {
						lw.Error(err)
					}
				}

				// the problem is that P1 thinks P2 is suspect, while P2 thinks that P1 is suspect.
				// this leads to a stalemate.
			case gmcom.CStatusSuspect:
				// regardless of the pinger's incarnation number - if the local process-id
				// is marked as suspect in the pingers memberMap, increase the local process
				// incarnation number.
				if p.MemberMap != nil {
					// BETA: is the local process-id (Pi) marked as suspect or failed in the sender's memberMap?
					err := gm.updateOwnStatusInPingMap(&p)
					if err != nil {
						lw.Error(err)
					}
				}

				if p.SrcIncarnationNumber > pj.IncarnationNumber {
					pj.Status = gmcom.CStatusAlive
					pj.StatusCount = 1
					pj.IncarnationNumber = p.TargetIncarnationNumber

					// update the pinger's information in the local memberMap
					err := gm.memberMap.Add(pj.ID, pj)
					if err != nil {
						lw.Error(err)
					}

					// perform local updates based on incoming ping memberMap
					if p.MemberMap != nil {
						// // BETA: is the local process-id (Pi) marked as suspect or failed in the sender's memberMap?
						// err = gm.updateOwnStatusInPingMap(&p)
						// if err != nil {
						// 	// TODO: err
						// 	log.Println(err)
						// }

						// update the local memberMap based on the incoming ping's memberMap
						rm.CheckElection, err = gm.memberMap.UpdFromPing(gm.MyID, pj.ID, p.MemberMap, uint(gm.FailureThreshold), gm.Leader.LeaderID)
						if err != nil {
							lw.Error(err)
						}
					}
				}

			case gmcom.CStatusFailed:
				// do nothing

			default:
				// do nothing

			}
		}
	}
	rm.MemberMap = gm.memberMap
	return rm
}

// updateOwnStatusInPingMap is used to check the local process-id's (Pi) state
// from the perspective of the pinger.  If the pinger has the local process-id (Pi)
// marked as suspect, Pi's incarnation number must be updated, and the incoming
// PingMap is updated so that our localMap will get updated correctly in the
func (gm *GMServ) updateOwnStatusInPingMap(p *gmcom.GMMessage) error {

	lw.Info("updateOwnStatusInPingMap() called...")
	// is the local process-id (Pi) marked as suspect or failed in the sender's memberMap?
	pi, ok := p.MemberMap.IDMap[gm.MyID]
	if ok {
		switch pi.Status {
		case gmcom.CStatusAlive:
			lw.Debug("Incoming memberMap from pinger says that I am alive:")
			lw.Debug("Ping map: %v", p.MemberMap.IDMap)
		case gmcom.CStatusSuspect:
			lw.Debug("I am suspected of failure - increasing incarnation number and asserting alive status")
			gm.MyIncarnation++
			pi.Status = gmcom.CStatusAlive
			pi.IncarnationNumber = gm.MyIncarnation
			err := gm.memberMap.Add(gm.MyID, pi)
			if err != nil {
				return err
			}
		case gmcom.CStatusFailed:
			panic(fmt.Errorf("process %d sent ping to process %d (me) indicating process %d has failed", p.SrcID, gm.MyID, gm.MyID))
		default:
			panic(fmt.Errorf("process %d sent ping to process %d (me) and process-id %d was not in memberMap", p.SrcID, gm.MyID, gm.MyID))
		}
	} else {
		panic(fmt.Errorf("updateOwnStatusInPingMap: gm.MyID: %d not read in p.MemberMap.IDMap: %v", gm.MyID, p.MemberMap.IDMap))
	}
	return nil
}
