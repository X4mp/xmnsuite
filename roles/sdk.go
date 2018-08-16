package roles

import (
	"github.com/XMNBlockchain/datamint/keys"
	"github.com/XMNBlockchain/datamint/users"
)

// Roles represents a role
type Roles interface {
	Keys() keys.Keys
	Add(key string, usrs ...users.Users) error
	Del(key string, usrs ...users.Users) error
	EnableWriteAccess(key string, keyPatterns ...string) error
	DisableWriteAccess(key string, keyPatterns ...string) error
	AddControl(fromKey string, toKey string) error
}
