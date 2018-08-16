package lists

import "github.com/XMNBlockchain/datamint/objects"

// WalkFn represents the func called by walk
type WalkFn func(index int, value []byte) (interface{}, error)

// Lists represents the lists data store
type Lists interface {
	Objects() objects.Objects
	Add(key string, values ...interface{}) int
	Retrieve(key string, index int, amount int) []interface{}
	Len(key string) int
	Union(key ...string) []interface{}
	UnionStore(destination string, key ...string) int
	/*Inter(key ...string) []interface{}
	InterStore(destination string, key ...string) int
	Push(key string, values ...interface{}) int
	PushX(key string, values ...interface{}) int
	Pop(key string) interface{}
	Trim(key string, start int, stop int) error
	Walk(key string, fn WalkFn) []interface{}*/
}
