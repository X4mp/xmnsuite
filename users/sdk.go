package users

import (
	"github.com/XMNBlockchain/datamint/hashtree"
	crypto "github.com/tendermint/tendermint/crypto"
)

// Users represents the users access control
type Users interface {
	Head() hashtree.Hash
	Add(pubKey crypto.PubKey) error
	Delete(fromPubKey crypto.PubKey, toDelPubKey crypto.PubKey) error
	AddWriteAccess(fromPubKey crypto.PubKey, toPubKey crypto.PubKey, keyPatterns ...string) error
	AddControl(fromPubKey crypto.PubKey, toPubKey crypto.PubKey, affectedPubKeys crypto.PubKey) error
}
