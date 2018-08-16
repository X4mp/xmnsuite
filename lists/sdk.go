package lists

import "github.com/XMNBlockchain/redismint/hashtree"

// WalkFn represents the func called by walk
type WalkFn func(index int, value []byte) (interface{}, error)

// Lists represents the lists data store
type Lists interface {
	Head() hashtree.Hash
	HashTree(keys string) hashtree.HashTree
	HashTrees(keys ...string) []hashtree.HashTree
	Add(key string, values ...[]byte) int
	Retrieve(key string, index int, amount int) [][]byte
	Len(key string) int
	Union(key ...string) [][]byte
	UnionStore(destination string, key ...string) int
	Inter(key ...string) [][]byte
	InterStore(destination string, key ...string) int
	Delete(key ...string) int
	Push(key string, values ...[]byte) int
	PushX(key string, values ...[]byte) int
	Pop(key string) []byte
	Trim(key string, start int, stop int) error
	Walk(key string, fn WalkFn) []interface{}
}
