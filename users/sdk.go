package users

import (
	"github.com/XMNBlockchain/datamint/hashtree"
	crypto "github.com/tendermint/tendermint/crypto"
)

// Users represents the users access control
type Users interface {
	Head() hashtree.HashTree
	Len() int
	Exists(pubKey crypto.PubKey) bool
	Insert(pubKey crypto.PubKey) error
	Delete(pubKey crypto.PubKey) error
}

// SDKFunc represents the users SDK func
var SDKFunc = struct {
	Create func() Users
}{
	Create: func() Users {
		return createConcreteUsers()
	},
}
