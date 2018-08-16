package objects

// ObjInKey represents an object in key
type ObjInKey struct {
	Key string
	Obj interface{}
}

// Objects represents the object data store
type Objects interface {
	Retrieve(objs ...ObjInKey) int
	Save(objs ...ObjInKey) int
	Delete(key ...string) int
}
