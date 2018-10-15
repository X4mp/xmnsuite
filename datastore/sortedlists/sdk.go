package sortedlists

import (
	"github.com/xmnservices/xmnsuite/datastore/keys"
	"github.com/xmnservices/xmnsuite/datastore/lists"
	"github.com/xmnservices/xmnsuite/hashtree"
)

// WalkFn represents the func called by walk
type WalkFn func(index int, score int, value []interface{}) (interface{}, error)

// ValueScore represents a value with a score
type ValueScore struct {
	Value []interface{}
	Score int
}

// SortedLists represents the lists data store
type SortedLists interface {
	Lists() lists.Lists
	Scores() keys.Keys
	Head() hashtree.Hash
	HashTree(keys string) hashtree.HashTree
	HashTrees(keys ...string) []hashtree.HashTree
	Add(key string, values ...*ValueScore) int
	Del(key string, values ...interface{}) int
	Retrieve(key string, index int, amount int) []*ValueScore
	Union(key ...string) []*ValueScore
	UnionStore(destination string, key ...string) int
	Inter(key ...string) []*ValueScore
	InterStore(destination string, key ...string) int
	Trim(key string, start int, stop int) error
	Walk(key string, fn WalkFn) []interface{}
	WalkStore(destination string, key string, fn WalkFn) int
}
