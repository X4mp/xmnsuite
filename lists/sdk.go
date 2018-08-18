package lists

import "github.com/XMNBlockchain/datamint/objects"

// WalkFn represents the func called by walk
type WalkFn func(index int, value interface{}) (interface{}, error)

// Lists represents the lists data store
type Lists interface {
	Objects() objects.Objects
	Add(key string, values ...interface{}) int
	Del(key string, values ...interface{}) int
	Retrieve(key string, index int, amount int) []interface{}
	Len(key string) int
	Union(key ...string) []interface{}
	UnionStore(destination string, key ...string) int
	Inter(key ...string) []interface{}
	InterStore(destination string, key ...string) int
	Trim(key string, index int, amount int) int
	Walk(key string, fn WalkFn) []interface{}
	WalkStore(destination string, key string, fn WalkFn) int
}

// SDKFunc represents the objects SDK func
var SDKFunc = struct {
	CreateList func() Lists
	CreateSet  func() Lists
}{
	CreateList: func() Lists {
		return createConcreteLists(false)
	},
	CreateSet: func() Lists {
		return createConcreteLists(true)
	},
}
