package objects

import (
	"encoding/json"

	"github.com/XMNBlockchain/datamint/hashtree"
)

/*
 * Stored Instance
 */

type storedInstance struct {
	HT  hashtree.HashTree
	Obj interface{}
}

func createStoredInstance(obj interface{}) *storedInstance {
	js, jsErr := json.Marshal(obj)
	if jsErr != nil {
		panic(jsErr)
	}

	ht, htErr := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: [][]byte{
			js,
		},
	})

	if htErr != nil {
		panic(htErr)
	}

	out := storedInstance{
		HT:  ht,
		Obj: obj,
	}

	return &out
}

/*
 * Concrete Objects
 */

type concreteObjects struct {
	head hashtree.HashTree
	data map[string]*storedInstance
}

// Head returns the hash head
func (app *concreteObjects) Head() hashtree.HashTree {
	return app.head
}

// HashTree returns the hashtree of the object at key
func (app *concreteObjects) HashTree(key string) hashtree.HashTree {
	if app.Exists(key) == 1 {
		return app.data[key].HT
	}

	return nil
}

// HashTrees returns the hashtrees of the objects at keys
func (app *concreteObjects) HashTrees(keys ...string) []hashtree.HashTree {
	out := []hashtree.HashTree{}
	for _, oneKey := range keys {
		out = append(out, app.HashTree(oneKey))
	}

	return out
}

// Exists returns the amount of keys passed to Exists that exists
func (app *concreteObjects) Exists(key ...string) int {
	cpt := 0
	for _, oneKey := range key {
		if _, ok := app.data[oneKey]; ok {
			cpt++
		}
	}
	return cpt
}

// Retrieve populates the Obj pointers in the passed ObjInKey instances.  Returns the amount of instances retrieved
func (app *concreteObjects) Retrieve(objs ...ObjInKey) int {
	cpt := 0
	for _, oneObj := range objs {
		if app.Exists(oneObj.Key) == 1 {
			oneObj.Obj = app.data[oneObj.Key].Obj
			cpt++
		}
	}

	return cpt
}

// Save saves the Obj at key as explained in the passed ObjInKey instances
func (app *concreteObjects) Save(objs ...ObjInKey) int {
	cpt := 0
	for _, oneObj := range objs {
		app.data[oneObj.Key] = createStoredInstance(oneObj.Obj)
		cpt++
	}

	//rebuild the head:
	app.rebuildHead()

	//returns the amount of saved keys:
	return cpt
}

// Delete deletes the passed keys
func (app *concreteObjects) Delete(key ...string) int {
	cpt := 0
	for _, oneKey := range key {
		if _, ok := app.data[oneKey]; ok {
			delete(app.data, oneKey)
			cpt++
		}
	}

	//rebuild the head:
	app.rebuildHead()

	//returns the amount of deleted keys:
	return cpt
}

func (app *concreteObjects) rebuildHead() {
	blocks := [][]byte{}
	for keyname, oneStoredIns := range app.data {
		blocks = append(blocks, []byte(keyname))
		blocks = append(blocks, oneStoredIns.HT.Head().Get())
	}

	ht, htErr := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: blocks,
	})

	if htErr != nil {
		panic(htErr)
	}

	app.head = ht
}
