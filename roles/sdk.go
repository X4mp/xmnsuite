package roles

import (
	"github.com/XMNBlockchain/datamint/lists"
	crypto "github.com/tendermint/tendermint/crypto"
)

// Roles represents a role
type Roles interface {
	Lists() lists.Lists
	Add(key string, usrs ...crypto.PubKey) int
	Del(key string, usrs ...crypto.PubKey) int
	EnableWriteAccess(key string, keyPatterns ...string) int
	DisableWriteAccess(key string, keyPatterns ...string) int
	HasWriteAccess(key string, keys ...string) []string
}

// SDKFunc represents the Roles SDK func
var SDKFunc = struct {
	Create func() Roles
}{
	Create: func() Roles {
		return createConcreteRoles()
	},
}
