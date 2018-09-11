package sortedlists

import "github.com/XMNBlockchain/xmnsuite/hashtree"

// WalkFn represents the func called by walk
type WalkFn func(index int, score int, value []byte) (interface{}, error)

// ValueScore represents a value with a score
type ValueScore struct {
	Value []byte
	Score int
}

// SortedLists represents the lists data store
type SortedLists interface {
	Head() hashtree.Hash
	HashTree(keys string) hashtree.HashTree
	HashTrees(keys ...string) []hashtree.HashTree
	Add(key string, values ...*ValueScore) int
	Retrieve(key string, index int, amount int) []*ValueScore
	Len(key string) int
	Union(key ...string) []*ValueScore
	UnionStore(destination string, key ...string) int
	Inter(key ...string) []*ValueScore
	InterStore(destination string, key ...string) int
	Delete(key ...string) int
	Push(key string, values ...*ValueScore) int
	PushX(key string, values ...*ValueScore) int
	Pop(key string) []byte
	Trim(key string, start int, stop int) error
	Walk(key string, fn WalkFn) []interface{}
}
