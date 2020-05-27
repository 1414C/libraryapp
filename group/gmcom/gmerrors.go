package gmcom

// Join errors
//
//  err = "No contact with leader process"
//  return &joinError{err, CerrJoinNoContact}

const (
	CerrJoinNoContact         = "errJoinNoContact"
	CerrJoinAddrInUse         = "errJoinAddrInUse"
	CerrJoinNotLeader         = "errJoinNotLeader"
	CerrJoinUnknownErr        = "errJoinUnknownErr"
	CerrPingIncorrectReceiver = "errPingIncorrectReceiver"
	CerrMemberMapFlushFailed  = "errMemberMapFlushFailed"
)

// JoinError is used to differentiate between join failures
type JoinError struct {
	Err     string //error description
	ErrCode string
}

// Error returns the error string
func (e *JoinError) Error() string {
	return e.Err
}

// ErrorCode returns the error code
func (e *JoinError) ErrorCode() string {
	return e.ErrCode
}

// ErrorSummary returns the error information as a single string
func (e *JoinError) ErrorSummary() string {
	return "join error code: " + e.ErrCode + " join error text: " + e.Err
}

// Clear cleans up the error structure
func (e *JoinError) Clear() {
	e.Err = ""
	e.ErrCode = ""
}
