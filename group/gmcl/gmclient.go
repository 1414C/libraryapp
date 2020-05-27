package gmcl

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/1414C/libraryapp/group/gmcom"
	"golang.org/x/net/websocket"

	"github.com/1414C/lw"
)

// Cache access errors
//
//  err = "Requested key was not found"
//  return &cacheError{err, CerrCacheKeyNotFoundError}

const (
	CerrCacheTechError        = "errCacheTechError"
	CerrCacheKeyNotFoundError = "errCacheKeyNotFoundError"
)

// CacheError is used to differentiate between types of cache-access failures
type CacheError struct {
	Err     string //error description
	ErrCode string
}

// Error returns the error string
func (e *CacheError) Error() string {
	return e.Err
}

// ErrorCode returns the error code
func (e *CacheError) ErrorCode() string {
	return e.ErrCode
}

// ErrorSummary returns the error information as a single string
func (e *CacheError) ErrorSummary() string {
	return "cache error code: " + e.ErrCode + " join error text: " + e.Err
}

// Clear cleans up the error structure
func (e *CacheError) Clear() {
	e.Err = ""
	e.ErrCode = ""
}

// GetErrorTest is used for testing error codes.
func GetErrorTest() error {
	err := fmt.Errorf("error message")
	return &CacheError{
		Err:     err.Error(),
		ErrCode: CerrCacheKeyNotFoundError,
	}
}

// AddUpdUsrCache adds or updates an entry in the local usr cache, resulting
// in a cascaded dispatch of the same call to all non-failed group-members.
// err := AddUpdUsrCache(gmcom.ActUsrD{Email:test@test.com, Active:true}, "192.168.1.66:4444")
func AddUpdUsrCache(u gmcom.ActUsrD, address string) error {

	// gob encode the Usr data
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(u)
	if err != nil {
		lw.ErrorWithPrefixString("failed to gob-encode Usr - got:", err)
		return err
	}
	encUsr := encBuf.Bytes()

	// connect to remote cache server
	origin := "http://localhost/"
	url := "ws://" + address + "/updateusrcache"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdUsrCache() ws connection failed - got:", err)
		return err
	}
	defer ws.Close()

	// push the encoded (reduced) Usr
	_, err = ws.Write(encUsr)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdUsrCache() ws.Write error - got:", err)
		return err
	}

	var msg = make([]byte, 64)

	// single read from the ws is okay here
	n, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdUsrCache() ws.Read error - got:", err)
		return err
	}

	// if update is confirmed do a little dance =)
	if string(msg[:n]) == "true" {
		// cw <- na
	} else {
		e := fmt.Errorf("AddUpdUsrCache() appeared to fail - got %v(raw),%v(string)", msg[:n], string(msg[:n]))
		lw.Error(e)
		return e
	}
	return nil
}

// AddUpdGroupAuthCache adds or updates an entry in the local group-auth cache, resulting
// in a cascaded dispatch of the same call to all non-failed group-members.
// err := AddUpdGroupAuthCache(gmcom.ActUsr{Email:test@test.com, Active:true}, "192.168.1.66:4444")
func AddUpdGroupAuthCache(g gmcom.GroupAuthD, address string) error {

	// gob encode the GroupAuth data
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(g)
	if err != nil {
		lw.ErrorWithPrefixString("failed to gob-encode GroupAuth - got:", err)
		return err
	}
	encGA := encBuf.Bytes()

	// connect to remote cache server
	origin := "http://localhost/"
	url := "ws://" + address + "/updategroupauthcache"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdGroupAuthCache() ws connection failed - got:", err)
		return err
	}
	defer ws.Close()

	// push the encoded (reduced) GroupAuth
	_, err = ws.Write(encGA)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdGroupAuthCache() ws.Write error - got:", err)
		return err
	}

	var msg = make([]byte, 64)

	// single read from the ws is okay here
	n, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdGroupAuthCache() ws.Read error - got:", err)
		return err
	}

	// if update is confirmed do a little dance =)
	if string(msg[:n]) == "true" {
		// cw <- na
	} else {
		e := fmt.Errorf("AddUpdGroupAuthCache() appeared to fail - got %v(raw),%v(string)", msg[:n], string(msg[:n]))
		lw.Error(e)
		return e
	}
	return nil
}

// AddUpdAuthCache adds or updates an entry in the local auth cache, resulting
// in a cascaded dispatch of the same call to all non-failed group-members.
// err := AddUpdUsrCache(gmcom.ActUsr{Email:test@test.com, Active:true}, "192.168.1.66:4444")
func AddUpdAuthCache(a gmcom.AuthD, address string) error {

	// gob encode the Usr data
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(a)
	if err != nil {
		lw.ErrorWithPrefixString("failed to gob-encode Usr - got:", err)
		return err
	}
	encAuth := encBuf.Bytes()

	// connect to remote cache server
	origin := "http://localhost/"
	url := "ws://" + address + "/updateauthcache"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdAuthCache() ws connection failed - got:", err)
		return err
	}
	defer ws.Close()

	// push the encoded Auth
	_, err = ws.Write(encAuth)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdAuthCache() ws.Write error - got:", err)
		return err
	}

	var msg = make([]byte, 64)

	// single read from the ws is okay here
	n, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdAuthCache() ws.Read error - got:", err)
		return err
	}

	// if update is confirmed do a little dance =)
	if string(msg[:n]) == "true" {
		// cw <- na
	} else {
		e := fmt.Errorf("AddUpdAuthCache() appeared to fail - got %v(raw),%v(string)", msg[:n], string(msg[:n]))
		lw.Error(e)
		return e
	}
	return nil
}

// AddUpdUsrGroupCache adds or updates an entry in the local usrgroup cache, resulting
// in a cascaded dispatch of the same call to all non-failed group-members.
// err := AddUpdUsrGroupCache(gmcom.ActUsr{Email:test@test.com, Active:true}, "192.168.1.66:4444")
func AddUpdUsrGroupCache(a gmcom.UsrGroupD, address string) error {

	// gob encode the UsrGroup data
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(a)
	if err != nil {
		lw.ErrorWithPrefixString("failed to gob-encode Usr - got:", err)
		return err
	}
	encAuth := encBuf.Bytes()

	// connect to remote cache server
	origin := "http://localhost/"
	url := "ws://" + address + "/updateusrgroupcache"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdUsrGroupCache() ws connection failed - got:", err)
		return err
	}
	defer ws.Close()

	// push the encoded UsrGroup
	_, err = ws.Write(encAuth)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdUsrGroupdCache() ws.Write error - got:", err)
		return err
	}

	var msg = make([]byte, 64)

	// single read from the ws is okay here
	n, err := ws.Read(msg)
	if err != nil {
		lw.ErrorWithPrefixString("AddUpdUsrGroupCache() ws.Read error - got:", err)
		return err
	}

	// if update is confirmed do a little dance =)
	if string(msg[:n]) == "true" {
		// cw <- na
	} else {
		e := fmt.Errorf("AddUpdUsrGroupCache() appeared to fail - got %v(raw),%v(string)", msg[:n], string(msg[:n]))
		lw.Error(e)
		return e
	}
	return nil
}
