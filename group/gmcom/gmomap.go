package gmcom

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:57 CDT
//=============================================================================================

import (
	"fmt"

	"github.com/1414C/lw"
)

type IDList []uint
type IDMap map[uint]GMMember

// OrderedMapI holds convenience accessors
type OrderedMapI interface {
	Add(id uint, g GMMember) error
	AddWithoutID(g *GMMember) error
	Read(id uint) (g GMMember, ok bool)
	ReadAll() (g []GMMember)
	ReadByIndex(idx int) (g GMMember, ok bool)
	ReadByAddress(ipAddress string) (g []GMMember)
	ReadCWNeighbour(id uint) (g GMMember, ok bool)
	ReadMaxActiveID() (g *GMMember)
	ReadMinActiveID() (g *GMMember)
	ReadActiveProcessList() (g []GMMember)
	UpdFromPing(localID uint, id uint, mm *OMap, failureThreshold uint, leaderID uint) (checkElect bool, err error)
	UpdFailedDepartedProcesses() error
	Delete(id uint) bool
	DeleteByIndex(idx int) bool
	Flush() bool
	Count() int
	sendFailureFlush(pid uint)
}

// OMap struct
type OMap struct {
	IDList IDList
	IDMap  IDMap
	OrderedMapI
}

// ensure consistency against interface
var _ OrderedMapI = &OMap{}

// Init creates the ordered map
func (o *OMap) Init() {
	o.IDMap = make(map[uint]GMMember)
}

// Add inserts or updates an entry to the ordered map
func (o *OMap) Add(id uint, g GMMember) error {

	// check for zero incarnation numbers (illegal)
	if g.IncarnationNumber == 0 {
		panic("OMap.Add called with incarnation number == 0")
	}

	// make sure the process has some sort of status
	if g.Status == "" {
		panic(fmt.Errorf("OMap.Add g.Status == space - got: %v", g))
	}

	// add or update?
	_, ok := o.IDMap[id]
	if !ok {
		if id == 0 {
			panic(id)
		}
		o.IDList = append(o.IDList, id)
	}
	o.IDMap[id] = g
	return nil
}

// AddWithoutID inserts an entry and assigns an initial id
// based on the map's max(id) + 1
func (o *OMap) AddWithoutID(g *GMMember) error {

	// check for zero incarnation numbers (illegal)
	if g.IncarnationNumber == 0 {
		panic("OMap.Add called with incarnation number == 0")
	}

	// prevent double-entry
	if g.ID != 0 {
		return fmt.Errorf("AddWithoutID called with GMMember.ID != 0 - got ID: %d", g.ID)
	}

	// make sure the process has some sort of status
	if g.Status == "" {
		panic(fmt.Errorf("g.Status == space - got: %v", g))
	}

	// get max(ID) in the most current known memberMap
	var maxID uint
	for _, v := range o.IDList {
		if v > maxID {
			maxID = v
		}
	}
	g.ID = maxID + 1
	o.IDMap[g.ID] = *g
	o.IDList = append(o.IDList, g.ID)
	lw.Info("AddWithoutID resulted in o.IDMap: &v", o.IDMap)
	// log.Println("AddWithoutID resulted in o.IDMap:", o.IDMap)
	return nil
}

// UpdFailedDepartedProcesses increments the process's status count.  when the status count
// reaches the threshold value, the failed or departed process is removed from the membermap.
func (o *OMap) UpdFailedDepartedProcesses() error {

	c := ((o.Count() * 2) - 1)
	for k, v := range o.IDMap {
		if v.Status != CStatusFailed && v.Status != CStatusDeparted {
			continue
		}

		v.StatusCount++
		if v.StatusCount < uint(c) { // update
			err := o.Add(k, v)
			if err != nil {
				return err //?
			}
			continue
		} else { // remove
			b := o.Delete(k)
			if !b {
				if v.Status == CStatusFailed {
					return fmt.Errorf("failed to remove failed process %d from local membermap", k)
				}
				return fmt.Errorf("failed to remove departed process %d from local membermap", k)
			}
		}
	}
	return nil
}

// UpdFromDeparting updates the status of the departing process if the departing
// process has provided an incarnation number greater that that contained in the
// member-map.  (Sets DEPARTED status on departing process's entry in the local
// process's memberMap.)
func (o *OMap) UpdFromDeparting(g GMMember, leaderID uint) (checkElect bool, err error) {

	// does the departing process exist in the local map?
	w, ok := o.Read(g.ID)

	if ok {
		if w.IncarnationNumber < g.IncarnationNumber {
			lw.Info("DEPARTING: g.IncarnationNumber (%d) > w.IncarnationNumber (%d)\n", g.IncarnationNumber, w.IncarnationNumber)
			return false, fmt.Errorf("UpdFromDeparting() did not process departure request: incarnation number %d is not greater than %d", g.IncarnationNumber, w.IncarnationNumber)
		}
	}

	if ok && w.Status != CStatusFailed && w.Status != CStatusDeparted {
		w.StatusCount = 1
		w.Status = CStatusDeparted
		err := o.Add(g.ID, w)
		lw.Info("DEPARTING: set DEPARTED status for process %v in local membermap.\n", g.ID)
		if err != nil {
			// TODO: log
		}

		// START AN ELECTION?
		if leaderID == g.ID {
			checkElect = true
		}
	}
	return checkElect, err
}

// UpdFromPing updates the memberMap based on the memberMap
// included in the incoming PING message.  Pj(memberMap) is
// used to update Pi(memberMap) where Pj is the sender of
// the PING message and Pi is the receiver of the PING message.
func (o *OMap) UpdFromPing(localID uint, id uint, mm *OMap, failureThreshold uint, leaderID uint) (checkElect bool, err error) {

	lw.Info("UpdFromPing() got the following from: %v", id)
	for _, v := range mm.IDMap {
		lw.Info("UpdFromPing() %v %v", id, v)
	}

	// process the incoming dissemination map in any order
	for _, v := range mm.IDMap {

		// already operated on the pinger / the local process-id
		if v.ID == id || v.ID == localID {
			continue
		}

		// if the current loop process-id (v.ID) does not exist in the local
		// memberMap it should be added according to the following rules:
		// 1. do not add processes that are in a failed state.
		// 2. do not add processes that are in a departed state.
		// 3. reset the incoming status count to 1 before adding the new process.
		w, ok := o.Read(v.ID)
		if !ok && v.Status != CStatusFailed && v.Status != CStatusDeparted {
			v.StatusCount = 1
			err := o.Add(v.ID, v)
			if err != nil {
				// TODO: log
			}
			continue
		}

		// if the current loop process-id (v.ID) exists in the local memberMap
		// with a failed or departed status no updates should be made.
		if w.Status == CStatusFailed || w.Status == CStatusDeparted {
			continue
		}

		// if the current loop process-id has an incoming status of failed, and
		// the local memberMap has any other status then update the local memberMap
		// accordingly.
		if ok && v.Status == CStatusFailed {
			w.Status = CStatusFailed
			w.StatusCount = 1
			err := o.Add(w.ID, w)
			if err != nil {
				// TODO: log
			}
			// START AN ELECTION?
			lw.Info("ELECTION: SET checkElect BASED ON INCOMING MEMBERMAP FAILURE INDICATION")
			if leaderID == v.ID {
				checkElect = true
			}
			continue
		}

		// if the current loop process-id has an incoming status of departed, and
		// the local memberMap has any other status then update the local memberMap
		// accordingly.
		if ok && v.Status == CStatusDeparted {
			w.Status = CStatusDeparted
			w.StatusCount = 1
			err := o.Add(w.ID, w)
			if err != nil {
				// TODO: log
			}
			// start an election?
			lw.Info("ELECTION: SET checkElect BASED ON INCOMING MEMBERMAP DEPARTED INDICATION")
			if leaderID == v.ID {
				checkElect = true
			}
			continue
		}

		// if the current loop process-id has a suspect status, update the
		// local memberMap according to the following rules:
		// 1. if the current loop process-id has an incarnation number less
		//    than that of the same process-id's locally stored incarnation
		//    number then perform no updates.
		// 2. if the current loop process-id has an incarnation number >=
		//    that of the same process-id's locally stored incarnation
		//    number then perform an update.
		if ok && v.Status == CStatusSuspect {
			if v.IncarnationNumber < w.IncarnationNumber {
				lw.Info("INCARNATION: SUSPECT: v.IncarnationNumber (%d) < w.IncarnationNumber (%d)", v.IncarnationNumber, w.IncarnationNumber)
				continue
			}

			// take action based on local understanding (w) of the process-id's state
			switch w.Status {
			case CStatusAlive:
				w.Status = CStatusSuspect
				w.StatusCount = 1
				if v.IncarnationNumber > w.IncarnationNumber {
					w.IncarnationNumber = v.IncarnationNumber
				}
				err := o.Add(w.ID, w)
				if err != nil {
					// TODO: log
				}

			case CStatusSuspect:
				w.StatusCount++
				if w.StatusCount < failureThreshold { // 200
					if v.IncarnationNumber > w.IncarnationNumber {
						w.IncarnationNumber = v.IncarnationNumber
						err := o.Add(w.ID, w)
						if err != nil {
							// TODO: log
						}
					}
				} else {
					w.Status = CStatusFailed
					w.StatusCount = 1
					err := o.Add(w.ID, w)
					if err != nil {
						// TODO: log
					}

					// START AN ELECTION?
					lw.Info("ELECTION: SET checkElect BASED ON SUPECT -> FAILURE TRANSITION")
					if leaderID == v.ID {
						checkElect = true
					}
					continue
				}

			case CStatusFailed:
				// do nothing - the status count of failed process is incremented
				// in sendPing() in order to track the dissemination count.

			case CStatusDeparted:
				// do nothing - the status count of departed process will be incremented
				// in sendPing() in order to track the dissemination count.

			default:
				// do nothing

			}
			continue
		}

		// if the current loop process-id has an alive status, update the
		// local memberMap according to the following rules:
		// 1. if the current loop process-id has an incarnation number less
		//    than that of the same process-id's locally stored incarnation
		//    number then perform no updates.
		// 2. if the current loop process-id has an incarnation number >=
		//    that of the same process-id's locally stored incarnation
		//    number then perform an update.
		if ok && v.Status == CStatusAlive {
			if v.IncarnationNumber < w.IncarnationNumber {
				lw.Info("INCARNATION: ALIVE: v.IncarnationNumber (%d) < w.IncarnationNumber (%d)", v.IncarnationNumber, w.IncarnationNumber)
				continue
			}

			// take action based on local understanding of the process-id's state
			switch w.Status {
			case CStatusAlive:
				if v.IncarnationNumber >= w.IncarnationNumber {
					w.StatusCount++
					w.IncarnationNumber = v.IncarnationNumber
					err := o.Add(w.ID, w)
					if err != nil {
						// TODO: log
					}
				}

			case CStatusSuspect:
				if v.IncarnationNumber > w.IncarnationNumber {
					w.Status = CStatusAlive
					w.StatusCount = 1
					err := o.Add(w.ID, w)
					if err != nil {
						// TODO: log
					}
				}

			case CStatusFailed:
				// do nothing - failed statuses are immutable once set

			default:
				// do nothing
			}
		}
	}
	return checkElect, nil
}

// Read reads the specified entry from the ordered map by id
func (o *OMap) Read(id uint) (g GMMember, ok bool) {

	g, ok = o.IDMap[id]
	return g, ok
}

// ReadAll reads all entries from the process map irrespective of process-status
func (o *OMap) ReadAll() (g []GMMember) {

	for _, v := range o.IDMap {
		g = append(g, v)
	}
	return g
}

// ReadByIndex reads the specified entry from the ordered map by index value
func (o *OMap) ReadByIndex(idx int) (g GMMember, ok bool) {

	if (idx - 1) > len(o.IDList) {
		return g, false
	}

	id := o.IDList[idx]
	return o.Read(id)
}

// ReadCWNeighbour returns the Pid and IPAddress of the process clockwise to
// process id.
func (o *OMap) ReadCWNeighbour(id uint) (g GMMember, ok bool) {

	var ti int
	c := o.Count()

	// if the calling process is the sole process in the
	// group, set ok to true and return the calling
	// GMMember struct
	if c == 1 {
		return g, true
	}

	idx := -1
	for i, v := range o.IDList {
		if v == id {
			idx = i
			break
		}
	}

	// the calling process was not found in the list
	// TODO
	if idx == -1 {
		return g, false
	}

	// check if the calling process is at the end of the list
	if (idx + 1) == c {
		ti = 0
	} else {
		ti = idx + 1
	}
	return o.ReadByIndex(ti)
}

// ReadMaxActiveID returns the details of the highest Active (or Suspect) ID in the ordered map.
func (o *OMap) ReadMaxActiveID() (g *GMMember) {

	if o.Count() == 0 {
		return nil
	}

	var k uint
	var maxG GMMember
	var tMember GMMember
	var ok bool

	for _, v := range o.IDList {
		if v > k {
			maxG, ok = o.Read(v)
			if ok && (maxG.Status != CStatusFailed && maxG.Status != CStatusDeparted) {
				k = v
				tMember = maxG
			}
		}
	}
	if k == 0 {
		return nil
	}
	return &tMember
}

// ReadMinActiveID returns the details of the lowest Active (or Suspect) ID in the ordered map.
func (o *OMap) ReadMinActiveID() (g *GMMember) {

	if o.Count() == 0 {
		return nil
	}

	var k uint
	var minG GMMember
	var bMember GMMember
	var ok bool

	k = ^uint(0) // max-value
	for _, v := range o.IDList {
		if v < k {
			minG, ok = o.Read(v)
			if ok && (minG.Status != CStatusFailed && minG.Status != CStatusDeparted) {
				k = v
				bMember = minG
			}
		}
	}
	if k == 0 {
		return nil
	}
	return &bMember
}

// ReadByAddress returns all entries in the local membermap containing an
// ipAddress (<addresss:port>) matching parameter ipAddress.  this is useful
// when the leader is deciding whether to accept or reject an incoming JOIN
// request.  This is intended for initialization / addition use only.
func (o *OMap) ReadByAddress(ipAddress string) (g []GMMember) {

	if o.IDMap == nil {
		return nil
	}

	for _, v := range o.IDMap {
		if v.IPAddress == ipAddress {
			g = append(g, v)
		}
	}
	if len(g) > 0 {
		return g
	}
	return nil
}

// ReadActiveProcessList returns a list of active processes in the GMMember format.
func (o *OMap) ReadActiveProcessList() (g []GMMember) {

	if o.Count() == 0 {
		return nil
	}

	var m GMMember
	var ok bool

	for _, v := range o.IDList {
		m, ok = o.Read(v)
		if ok && m.Status != CStatusFailed && m.Status != CStatusDeparted {
			g = append(g, m)
		}
	}
	return g
}

// Delete deletes the specified entry from the ordered map by id
func (o *OMap) Delete(id uint) bool {

	// deleting from an empty list is deemed to be successful
	if len(o.IDList) == 0 {
		return true
	}

	for n, v := range o.IDList {
		if v == id {
			ok := o.DeleteByIndex(n)
			return ok
		}
	}
	return false
}

// DeleteByIndex deletes the specified entry from the ordered map by index value
func (o *OMap) DeleteByIndex(idx int) bool {

	if len(o.IDList) == 1 {
		o.IDList = o.IDList[:0]
		return true
	}

	v := o.IDList[idx]
	o.IDList = append(o.IDList[:idx], o.IDList[idx+1:]...)
	delete(o.IDMap, v)
	return true
}

// Flush flushes the ordered map with no checks
func (o *OMap) Flush() bool {

	// deleting from an empty list is deemed to be successful
	if len(o.IDList) == 0 {
		return true
	}

	for n := range o.IDList {
		ok := o.DeleteByIndex(n)
		if !ok {
			return ok
		}
		return ok
	}
	return true
}

// Count returns the current count of the ordered map
func (o *OMap) Count() int {

	return len(o.IDList)
}
