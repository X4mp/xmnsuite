package keys

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/XMNBlockchain/datamint/hashtree"
)

/*
 * Helper func
 */

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

/*
 * Stored Instance
 */

type storedInstance struct {
	HT   hashtree.HashTree
	Data interface{}
}

func createStoredInstance(data interface{}) *storedInstance {
	blocks, blocksErr := getBytes(data)
	if blocksErr != nil {
		str := fmt.Sprintf("the data could not be converted to []byte: %s", blocksErr.Error())
		panic(errors.New(str))
	}

	ht := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: [][]byte{
			blocks,
		},
	})

	out := storedInstance{
		HT:   ht,
		Data: data,
	}

	return &out
}

/*
 * Concrete Keys
 */

type concreteKeys struct {
	head hashtree.HashTree
	data map[string]*storedInstance
}

func createConcreteKeys() Keys {
	out := concreteKeys{
		head: nil,
		data: map[string]*storedInstance{},
	}

	out.rebuildHead()
	return &out
}

// Head returns the hash head
func (app *concreteKeys) Head() hashtree.HashTree {
	return app.head
}

// HashTree returns the hashtree of the object at key
func (app *concreteKeys) HashTree(key string) hashtree.HashTree {
	if app.Exists(key) == 1 {
		return app.data[key].HT
	}

	return nil
}

// HashTrees returns the hashtrees of the objects at keys
func (app *concreteKeys) HashTrees(keys ...string) []hashtree.HashTree {
	out := []hashtree.HashTree{}
	for _, oneKey := range keys {
		out = append(out, app.HashTree(oneKey))
	}

	return out
}

// Len returns the amount of objects stored
func (app *concreteKeys) Len() int {
	return len(app.data)
}

// Exists returns the amount of keys passed to Exists that exists
func (app *concreteKeys) Exists(key ...string) int {
	cpt := 0
	for _, oneKey := range key {
		if _, ok := app.data[oneKey]; ok {
			cpt++
		}
	}
	return cpt
}

// Retrieve retrieves data at key
func (app *concreteKeys) Retrieve(key string) interface{} {
	if app.Exists(key) == 1 {
		return app.data[key].Data
	}

	return nil
}

// Search searches the keys using a pattern, and returns the keys that matches
func (app *concreteKeys) Search(pattern string) []string {
	reg, regErr := regexp.Compile(pattern)
	if regErr != nil {
		str := fmt.Sprintf("the given pattern is invalid: %s", regErr.Error())
		panic(errors.New(str))
	}

	out := []string{}
	for oneKeyname := range app.data {
		if !reg.MatchString(oneKeyname) {
			continue
		}

		out = append(out, oneKeyname)
	}

	//sort then return:
	sort.Strings(out)
	return out
}

// Save saves data at key
func (app *concreteKeys) Save(key string, data interface{}) {
	//add the data:
	app.data[key] = createStoredInstance(data)

	//rebuild the head:
	app.rebuildHead()
}

// Delete deletes the passed keys
func (app *concreteKeys) Delete(key ...string) int {
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

func (app *concreteKeys) rebuildHead() {
	blocks := [][]byte{
		[]byte("root"),
	}

	for keyname, ins := range app.data {
		blks, blksErr := getBytes(ins.Data)
		if blksErr != nil {
			str := fmt.Sprintf("the data could not be converted to []byte: %s", blksErr.Error())
			panic(errors.New(str))
		}

		blocks = append(blocks, []byte(keyname))
		blocks = append(blocks, blks)
	}

	ht := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: blocks,
	})

	app.head = ht
}
