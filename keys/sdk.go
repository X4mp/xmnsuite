package keys

import "github.com/XMNBlockchain/datamint/hashtree"

// Keys represents the keys datastore
type Keys interface {
	Head() hashtree.HashTree
	HashTree(key string) hashtree.HashTree
	HashTrees(keys ...string) []hashtree.HashTree
	Len() int
	Exists(key ...string) int
	Retrieve(key string) interface{}
	Search(pattern string) []string
	Save(key string, data interface{})
	Delete(key ...string) int
}

// SDKFunc represents the Keys SDK func
var SDKFunc = struct {
	Create func() Keys
}{
	Create: func() Keys {
		return createConcreteKeys()
	},
}
