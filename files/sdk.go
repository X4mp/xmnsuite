package files

import (
	"net/url"

	"github.com/XMNBlockchain/redismint/hashtree"
)

// Files represents the files data store
type Files interface {
	HashTree(keys string) hashtree.HashTree
	HashTrees(keys ...string) []hashtree.HashTree
	Init(key string, ht hashtree.HashTree) bool
	Add(key string, values ...[]byte) int
	Delete(keys ...string) int
	Retrieve(key string, index int, amount int) ([]*url.URL, error)
	RetrieveByHashes(key string, hashes []byte) ([]*url.URL, error)
	Len(key string) int
}
