package users

import (
	"github.com/XMNBlockchain/xmnsuite/objects"
	crypto "github.com/tendermint/tendermint/crypto"
)

// Users represents the users access control
type Users interface {
	Objects() objects.Objects
	Copy() Users
	Key(pubKey crypto.PubKey) string
	Exists(pubKey crypto.PubKey) bool
	Insert(pubKey crypto.PubKey) bool
	Delete(pubKey crypto.PubKey) bool
}

// SDKFunc represents the users SDK func
var SDKFunc = struct {
	Create func() Users
}{
	Create: func() Users {
		return createConcreteUsers()
	},
}
