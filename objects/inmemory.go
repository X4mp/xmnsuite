package objects

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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
	js, jsErr := cdc.MarshalJSON(obj)
	if jsErr != nil {
		str := fmt.Sprintf("the object cannot be stored because it cannot be converted to JSON: %s", jsErr.Error())
		panic(errors.New(str))
	}

	ht := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: [][]byte{
			js,
		},
	})

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

func createObjects() Objects {

	ht := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: [][]byte{
			[]byte(strconv.Itoa(int(time.Now().UTC().UnixNano()))),
		},
	})

	out := concreteObjects{
		head: ht,
		data: map[string]*storedInstance{},
	}

	return &out
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

// Len returns the amount of objects stored
func (app *concreteObjects) Len() int {
	return len(app.data)
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
func (app *concreteObjects) Retrieve(objs ...*ObjInKey) int {
	cpt := 0
	for index, oneObj := range objs {
		if app.Exists(oneObj.Key) == 1 {
			objs[index].Obj = app.data[oneObj.Key].Obj
			cpt++
		}
	}

	return cpt
}

// Save saves the Obj at key as explained in the passed ObjInKey instances
func (app *concreteObjects) Save(objs ...*ObjInKey) int {
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
	blocks := [][]byte{
		[]byte(strconv.Itoa(int(time.Now().UTC().UnixNano()))),
	}

	for keyname, oneStoredIns := range app.data {
		blocks = append(blocks, []byte(keyname))
		blocks = append(blocks, oneStoredIns.HT.Head().Get())
	}

	ht := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: blocks,
	})

	app.head = ht
}
