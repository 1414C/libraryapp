package gmcom

import (
	"sync"
)

type OpType string

const (
	COpCreate OpType = "C"
	COpRead   OpType = "R"
	COpUpdate OpType = "U"
	COpDelete OpType = "D"
)

// ActUsrsH is used as the runtime-type of the active-usrs cache on the Usr controller.
type ActUsrsH struct {
	sync.RWMutex
	ActiveUsrs map[uint64]bool
}

// ActUsrD is a carrier structure for disseminating USRUPDATE messages to group-members.
type ActUsrD struct {
	Forward bool
	ID      uint64
	Active  bool
}

// GroupAuthNames is used in the GroupAuthsID map (enable deletion)
type GroupAuthNames struct {
	GroupName string
	AuthName  string
}

// GroupAuthsH is used as the runtime-type of the group-authorization cache on the GroupAuth controller.
type GroupAuthsH struct {
	sync.RWMutex
	GroupAuths   map[string]map[string]bool
	GroupAuthsID map[uint64]GroupAuthNames
}

// GroupAuthD is a carrier structure for disseminating GROUPAUTHUPDATE messages to group-members.
type GroupAuthD struct {
	Forward   bool
	ID        uint64 // delete-only!
	GroupName string
	AuthName  string
	GroupID   uint64
	AuthID    uint64
	Op        OpType
}

// AuthsH is used as the runtime-type of the authorization-description cache on the Auth controller.
type AuthsH struct {
	sync.RWMutex
	Auths map[uint64]string
}

// AuthD is a carrier structure for disseminating AUTHUPDATE messages to group-members.
type AuthD struct {
	Forward  bool
	ID       uint64
	AuthName string
	Op       OpType
}

// UsrGroupsH is used as the runtime-type of the usrgroup-description cache on the UsrGroup controller.
type UsrGroupsH struct {
	sync.RWMutex
	GroupNames map[uint64]string
}

// UsrGroupD is a carrier structure for disseminating USRGROUPUPDATE messages to group-members.
type UsrGroupD struct {
	Forward   bool
	ID        uint64
	GroupName string
	Op        OpType
}
