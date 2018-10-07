package roles

import (
	crypto "github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore/lists"
)

// Roles represents a role
type Roles interface {
	Lists() lists.Lists
	Copy() Roles
	Add(key string, usrs ...crypto.PublicKey) int
	Del(key string, usrs ...crypto.PublicKey) int
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
