package objects

import (
	"github.com/xmnservices/xmnsuite/datastore/keys"
)

// ObjInKey represents an object in key
type ObjInKey struct {
	Key string
	Obj interface{}
}

// Objects represents the object data store
type Objects interface {
	Keys() keys.Keys
	Copy() Objects
	Retrieve(objs ...*ObjInKey) int
	Save(objs ...*ObjInKey) int
}

// SDKFunc represents the objects SDK func
var SDKFunc = struct {
	Create func() Objects
}{
	Create: func() Objects {
		return createObjects()
	},
}
