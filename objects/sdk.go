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
	Exists(key string) bool
	Retrieve(objs ...ObjInKey) int
	Save(objs ...ObjInKey) int
	Delete(key ...string) int
}
