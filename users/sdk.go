package users

import (
	"github.com/XMNBlockchain/datamint/hashtree"
	crypto "github.com/tendermint/tendermint/crypto"
)

// Users represents the users access control
type Users interface {
	Head() hashtree.Hash
	Add(pubKey crypto.PubKey) error
	Delete(pubKey crypto.PubKey) error
}
