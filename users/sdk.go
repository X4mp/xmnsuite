package users

import (
	crypto "github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/objects"
)

// Users represents the users access control
type Users interface {
	Objects() objects.Objects
	Copy() Users
	Key(pubKey crypto.PublicKey) string
	Exists(pubKey crypto.PublicKey) bool
	Insert(pubKey crypto.PublicKey) bool
	Delete(pubKey crypto.PublicKey) bool
}

// SDKFunc represents the users SDK func
var SDKFunc = struct {
	Create func() Users
}{
	Create: func() Users {
		return createConcreteUsers()
	},
}
