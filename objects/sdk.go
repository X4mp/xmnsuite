package objects

import "github.com/XMNBlockchain/redismint/hashtree"

// ObjInKey represents an object in key
type ObjInKey struct {
	Key string
	Obj interface{}
}

// Objects represents the object data store
type Objects interface {
	Head() hashtree.Hash
	HashTree(keys string) hashtree.HashTree
	HashTrees(keys ...string) []hashtree.HashTree
	Retrieve(objs ...ObjInKey) int
	Save(objs ...ObjInKey) int
	Delete(key ...string) int
}
