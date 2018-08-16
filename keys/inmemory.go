package keys

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/XMNBlockchain/datamint/hashtree"
)

/*
 * Stored Instance
 */

type storedInstance struct {
	HT   hashtree.HashTree
	Data []byte
}

func createStoredInstance(data []byte) *storedInstance {
	ht := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: [][]byte{
			data,
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
func (app *concreteKeys) Retrieve(key string) []byte {
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
func (app *concreteKeys) Save(key string, data []byte) bool {
	//add the data:
	app.data[key] = createStoredInstance(data)

	//rebuild the head:
	app.rebuildHead()

	//returns:
	return true
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
		[]byte(strconv.Itoa(int(time.Now().UTC().UnixNano()))),
	}

	for keyname, ins := range app.data {
		blocks = append(blocks, []byte(keyname))
		blocks = append(blocks, ins.Data)
	}

	ht := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: blocks,
	})

	app.head = ht
}
