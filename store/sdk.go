package store

import (
	"github.com/XMNBlockchain/redismint/hashes"
	"github.com/XMNBlockchain/redismint/hashtree"
)

// Keystore represents the keystore
type Keystore interface {
	GetHead() hashtree.HashTree
	GetHashes() hashes.Hashes
}
