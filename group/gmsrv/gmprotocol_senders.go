package gmsrv

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"fmt"
	"math/rand"
	"sync/atomic"

	"github.com/1414C/libraryapp/group/gmcom"

	"github.com/1414C/lw"
)

// SendGetDBLeader sends a non-protocol message via the serialization channels
// in order to retrieve the current (or last) leader's process-id and contact
// information.  Leader information is persisted outside of the group-context
// in order to provide a common point of reference for joining processes.
// This message is intended only to be sent during process initialization and
// it does not retry if no leader can be read, nor does it test to see if the
// leader can be contacted.
func (gm *GMServ) SendGetDBLeader(m *gmcom.GMMessage) {

	// make some adjustments to the incoming message structure
	m.SrcID = gm.MyID
	m.SrcIPAddress = gm.MyIPAddress
	m.Type = gmcom.CinGetDBLeader
	m.ExpectResponse = true

	// put the inGetDBLeader command into the channel then read the response from the channel
	gm.chin <- *m
	*m = <-gm.chout
}

// SendSetDBLeader sends a non-protocol message via the serialization channels
// in order to set a new leader's process-id and contact information in the
// persistent store.  Leader information is persisted outside of the group-context
// in order to provide a common point of reference for joining processes.
// This message is intended only to be sent during process initialization and
// it does not retry if no leader can be read, nor does it test to see if the
// leader can be contacted.
func (gm *GMServ) SendSetDBLeader(m *gmcom.GMMessage) {

	// put the inSetDBLeader command into the channel / read the inSetDBLeader response from the channel
	m.Type = gmcom.CinSetDBLeader
	gm.chin <- *m
	*m = <-gm.chout
}

// SendSetLeader sends a non-protocol message via the serialization channels
// in order to set a new leader's process-id and contact information in the
// local server store.
func (gm *GMServ) SendSetLeader(m *gmcom.GMMessage) {

	// put the inSetLeader command into the channel / read the inSetLeader response from the channel
	m.Type = gmcom.CinSetLeader
	gm.chin <- *m
	*m = <-gm.chout
}

// SendGetLocalLeader sends a non-protocol message via the serialization channels
// in order to get the leader information as known by the local process.  Although
// the leader info is available directly via the method's structure it is inadvisable
// to access the information in a non thread-safe manner (yeah I know they are goroutines...)
func (gm *GMServ) SendGetLocalLeader(m *gmcom.GMMessage) {

	m.Type = gmcom.CinGetLocalLeader
	gm.chin <- *m
	*m = <-gm.chout
}

// SendGetLocalDetails sends a non-protocol message via the serialization channels
// in order to get the local process information including a deep-copy of the
// current memberMap.  It is inherent that the group-membership map content is
// stale as soon as it has been created, but lets share by communicating
// rather than communicate by sharing ;)
func (gm *GMServ) SendGetLocalDetails() *gmcom.GMMessage {

	m := gmcom.GMMessage{
		Type: gmcom.CinGetLocalDetails,
	}
	gm.chin <- m
	m = <-gm.chout
	return &m
}

// SendJoin is called during the group server initialization using
// the leader information read from the external persistent store.
// A false return indicates that the Join request was not acknowledged
// by the leader.
func (gm *GMServ) SendJoin(l *gmcom.GMLeader) error {

	if l == nil {
		return fmt.Errorf("No Leader")
	}
	return gm.sendPing(l)
}

// SendPing is used to send a ping message to target process Pj. The target-process
// is determined anew for each call, as the group membership list may have changed
// since the last call.
func (gm *GMServ) SendPing() {
	_ = gm.sendPing(nil)
}

// SendFailure sends a failure message to each process in the membership group
// with the exception of the failed process.  SendFailure is also sent to the
// local server and the failed process is removed from the local list on this
// basis.
func (gm *GMServ) SendFailure(f gmcom.GMMessage) {

	// ??? DEPRECATED ???
	panic("Deprecated SendFailure(...) called!")
	// if no subject process was provided leave
	if f.SubjectID == 0 || f.SubjectIPAddress == "" {
		// TODO: log
		return
	}

	// request a list of current process-ids
	f.Type = gmcom.CinDoSendPrep
	f.TargetPath = ""
	f.TargetID = 0
	f.TargetIPAddress = ""
	f.ExpectResponse = true
	gm.chin <- f

	// read the list from the channel one-by-one :(
	requested := true
	for requested {
		f = <-gm.chout
		if f.TargetID == 0 {
			break
		}
		// make some adjustments to the incoming message structure
		f.Type = gmcom.CFailure
		f.TargetPath = "/failure"
		f.ExpectResponse = false

		// send the failure message to the target
		err := gm.TxGMMessage(f)
		if err != nil {
			lw.Error(err)
			return
		}
	}
}

// SendDeparting sends a departing message to each process in the membership group
// including the local process, although the local list will likely not see the
// update (or care about it).
func (gm *GMServ) SendDeparting() {

	lw.Info("DEPARTING-1: SendDeparting called")

	// do this on the channel?
	gm.MyIncarnation++
	gm.memberMap.Add(gm.MyID, gmcom.GMMember{
		ID:                gm.MyID,
		IPAddress:         gm.MyIPAddress,
		Status:            gmcom.CStatusDeparted,
		StatusCount:       1,
		IncarnationNumber: gm.MyIncarnation,
	})

	f := gmcom.GMMessage{}

	// request a list of current process-ids
	f.Type = gmcom.CinDoSendPrep
	f.TargetPath = ""
	f.TargetID = 0
	f.TargetIPAddress = ""
	f.ExpectResponse = true
	gm.chin <- f

	// read the list from the channel one-by-one :(
	requested := true
	for requested {
		f = <-gm.chout
		if f.TargetID == 0 {
			lw.Info("DEPARTING-1B: SendDeparting got: %v", f)
			break
		}
		// make some adjustments to the incoming message structure
		f.SrcIncarnationNumber += 2 // double-down
		f.Type = gmcom.CDeparting
		f.TargetPath = "/departing"
		f.ExpectResponse = false

		lw.Info("DEPARTING-2: SendDeparting sending: %v", f)

		// send the departing message to the target
		err := gm.TxGMMessage(f)
		if err != nil {
			lw.ErrorWithPrefixString("DEPARTING-ERR: SendDeparting got:", err)
		}
	}
}

// SendCoordinator sends the new leader information from the SrcID to the TargetID.
func (gm *GMServ) SendCoordinator(c gmcom.GMMessage) error {
	c.TargetPath = "/coordinator"
	c.ExpectResponse = false
	err := gm.TxGMMessage(c) //(pl[iv])
	if err != nil {
		return err
	}
	return nil
}

// private method implementations
// sendPing is used to send normal pings at the prescribed protocol period
// and is also used to send Join messages to the leader that has been read
// from the persistent store.
// in the ping use-case leader parameter (*l) should be nil and the return
// parameter can be ignored.
func (gm *GMServ) sendPing(l *gmcom.GMLeader) error {

	p := gmcom.GMMessage{}
	pl := []gmcom.GMMessage{}
	o := make([]int, 0)

	// normal ping use case
	if l == nil {

		// get the ping target information
		p.Type = gmcom.CinDoSendPrep
		gm.chin <- p

		// read CinDoSendPrep responses from the channel
		read := true
		c := 0
		for read {
			p = <-gm.chout
			if p.TargetID == 0 {
				break
			}

			p.Type = gmcom.CPing
			p.TargetPath = "/ping"
			p.ExpectResponse = true
			pl = append(pl, p)
			o = append(o, c)
			c++
		}

		// randomize the ping order via reordering sequential slice []o
		rand.Shuffle(len(o), func(i, j int) {
			o[i], o[j] = o[j], o[i]
		})

	} else {
		// join scenario
		lw.Debug("in sendPing Join scenario")
		lw.Info("got leader info: %v", l)

		// populate SrcID, SrcIPAddresss, SrcIncarnationNumber
		p.Type = gmcom.CinGetMySrcInfo
		gm.chin <- p
		p = <-gm.chout

		// setup for join
		p.Type = gmcom.CJoin
		p.TargetID = l.LeaderID
		p.TargetIPAddress = l.LeaderIPAddress
		p.TargetStatus = gmcom.CStatusAlive // make the assumption that the leader is alive
		p.TargetPath = "/join"
		p.ExpectResponse = true
		pl = append(pl, p)
		o = append(o, 0)
	}

	// log.Println("sendPing assembled ping list:", pl)
	lw.Debug("sendPing assembled index slice: %v", o)

	// send the pings
	for _, iv := range o {

		// send the PING/JOIN to non-failed processes via the general use TxRx and get the
		// result from the Rx or an error.  error indicates that an ACK message was not
		// received for the ping.
		p = pl[iv]
		if p.TargetStatus != gmcom.CStatusFailed && p.TargetStatus != gmcom.CStatusDeparted {
			lw.Info("sending PING: %v", p)

			r, err := gm.TxRxGMMessage(p) //(pl[iv])
			if err != nil {
				lw.ErrorWithPrefixString("gm.TxRxGMMessage got error in sendPing:", err)
				// no ACK message received, so set the internal inNoAck message
				// and then push into the serialization channel in order to update
				// the target's status in the member list.
				p.Type = gmcom.CinNoAck
				p.TargetPath = ""

				// if the local process-id could not ping itself the most likely cause is that the
				// network interface is down (unplugged?).  the process is still technically viable
				// if the interface comes back up, and the other group members have not set a failed
				// status for the gm.MyID. keep the local process in an alive state and increase the
				// local incarnation number for each locally failed gm.MyID -> gm.MyID ping.
				if p.SrcID == p.TargetID {
					p.TargetStatus = gmcom.CStatusAlive // still running - network is not reachable
					p.TargetStatusCount++
					p.TargetIncarnationNumber++ // increase incarnation number

					// put the inNoAck command into the channel and read it back (wait)
					gm.chin <- p
					p = <-gm.chout
					continue
				}

				// set target process's status
				switch p.TargetStatus {
				case gmcom.CStatusAlive:
					// old status was alive, ping did not get ACK - process is now suspect
					p.TargetStatus = gmcom.CStatusSuspect
					p.TargetStatusCount = 1

					// put the inNoAck command into the channel and read it back (wait)
					gm.chin <- p
					p = <-gm.chout

				case gmcom.CStatusSuspect:
					// old status was suspect, add to suspect count or set Pj to failed
					// if p.TargetStatusCount < 200 {
					if p.TargetStatusCount <= uint(atomic.LoadUint32(&gm.FailureThreshold)) {
						p.TargetStatusCount++

						// put the inNoAck command into the channel and read it back (wait)
						gm.chin <- p
						p = <-gm.chout

					} else {
						// set failure status which will result in Pj's removal from the local
						// list as well.  set the failure on P for good measure.
						p.TargetStatus = gmcom.CStatusFailed
						p.TargetStatusCount = 1

						// put the inNoAck command into the channel and read it back (wait)
						gm.chin <- p
						p = <-gm.chout

						// was the failed process the leader?
						l := gmcom.GMMessage{
							Type: gmcom.CinGetLocalLeader,
						}
						gm.chin <- l
						l = <-gm.chout
						lw.Info("NOACK-sendPing() detected failed process-id: %v", p.TargetID)
						if l.TargetID == p.TargetID &&
							l.TargetIPAddress == p.TargetIPAddress {
							// START AN ELECTION?
							e := gmcom.GMMessage{
								Type: gmcom.CinStartElection,
							}
							lw.Info("Starting an election based on leader failure detection")
							gm.chin <- e
							e = <-gm.chout // temporary
						}
					}

				case gmcom.CStatusFailed:
					// should never get here as failed processes are not pinged

				case gmcom.CStatusDeparted:
					// should never get here as departed processes are not pinged

				default:
					// should never get here
				}

				// return a specific error if this is a join scenario, otherwise return general error (not evaluated)
				if l != nil {
					return &gmcom.JoinError{
						Err:     err.Error(),
						ErrCode: gmcom.CerrJoinNoContact,
					}
				}
				lw.ErrorWithPrefixString("ping transmission failure - got:", err)
				continue
			}

			// ACK message received - set process status
			if r.Type == gmcom.CAck {

				// in a JOINACK scenario, put the JOINACK command into the channel to add
				// the local process's new process-id into the local gm.memberMap and read
				// it back (wait).
				if l != nil {
					lw.Info("Received r: %v", r)
					if r.TargetID != 0 {
						r.Type = gmcom.CJoinAck
						gm.chin <- *r
						*r = <-gm.chout
						r.Type = gmcom.CAck
					} else {
						lw.ErrorWithPrefixString("JOIN failed - got:", fmt.Errorf("JOIN failure response received from leader"))
						return &gmcom.JoinError{
							Err:     "JOIN failure response received from leader",
							ErrCode: r.ErrCode,
						}
					}
				}

				// log.Printf("ACKproc: p.TargetStatus of proc-id %d is %s with incarnation %d\n", p.TargetID, p.TargetStatus, p.TargetIncarnationNumber)
				// log.Printf("ACKproc: r.SrcStatus of    proc-id %d is %s with incarnation %d\n", r.SrcID, r.SrcStatus, r.SrcIncarnationNumber)
				switch p.TargetStatus {
				case gmcom.CStatusAlive:
					// update alive ping count and incarnation number
					r.SrcStatus = gmcom.CStatusAlive
					r.SrcStatusCount = p.TargetStatusCount + 1

					// put the ACK command into the channel and read it back (wait)
					gm.chin <- *r
					*r = <-gm.chout

				case gmcom.CStatusSuspect:
					if r.SrcIncarnationNumber > p.TargetIncarnationNumber {
						r.SrcStatus = gmcom.CStatusAlive
						r.SrcStatusCount = 1

						// put the ACK command into the channel and read it back (wait)
						gm.chin <- *r
						*r = <-gm.chout

					} else {
						// old status was suspect and the ping has a stale incarnation number,
						// so leave the process as suspect and increment the TargetStatusCount.
						// if p.TargetStatusCount < 200 {
						if p.TargetStatusCount <= uint(atomic.LoadUint32(&gm.FailureThreshold)) {
							p.TargetStatusCount++
							gm.chin <- p
							p = <-gm.chout

						} else {
							// set failure status which will result in Pj's removal from the local
							// list as well.  set the failure on P for good measure.
							p.TargetStatus = gmcom.CStatusFailed
							p.TargetStatusCount = 1
							gm.chin <- p
							p = <-gm.chout

							// was the failed process the leader?
							l := gmcom.GMMessage{
								Type: gmcom.CinGetLocalLeader,
							}
							gm.chin <- l
							l = <-gm.chout
							lw.Info("ACK-sendPing() detected failed process-id: %v", p.TargetID)
							if l.TargetID == p.TargetID &&
								l.TargetIPAddress == p.TargetIPAddress {
								// START AN ELECTION?
								e := gmcom.GMMessage{
									Type: gmcom.CinStartElection,
								}
								lw.Info("Starting an election based on leader failure detection")
								gm.chin <- e
								e = <-gm.chout // temporary
							}
						}
					}

				case gmcom.CStatusFailed:
					// never reinstate once failed

				case gmcom.CStatusDeparted:
					// never reinstate once departed

				default:
					panic(fmt.Errorf("incoming ACK for process %d has unknown status %s", p.TargetID, p.TargetStatus))
				}

			} else {
				lw.Debug("MyID is: %d", gm.MyID)
				lw.Warning("Ping message received inappropriate message type as a response - got: %s", r.Type)
				lw.Warning("Failed Ping Message: %v", p)
				lw.Warning("Got Ping Response Message: %v", r)
				panic(fmt.Errorf("Ping message received inappropriate message type as a response - got: %s", r.Type))
			}
		} else {
			// update the failure status (dissemination) count in the local membermap
		}
	}
	return nil
}
