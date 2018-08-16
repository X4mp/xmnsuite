package objects

import "github.com/XMNBlockchain/datamint/hashtree"

// ObjInKey represents an object in key
type ObjInKey struct {
	Key string
	Obj interface{}
}

// Objects represents the object data store
type Objects interface {
	Head() hashtree.HashTree
	HashTree(key string) hashtree.HashTree
	HashTrees(keys ...string) []hashtree.HashTree
	Len() int
	Exists(key ...string) int
	Retrieve(objs ...*ObjInKey) int
	Save(objs ...*ObjInKey) int
	Delete(key ...string) int
}

// SDKFunc represents the objects SDK func
var SDKFunc = struct {
	Create func() Objects
}{
	Create: func() Objects {
		return createObjects()
	},
}
