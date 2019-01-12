package keys

import (
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/xmnservices/xmnsuite/hashtree"
	"github.com/xmnservices/xmnsuite/helpers"
)

/*
 * Stored Instance
 */

type storedInstance struct {
	HT   hashtree.HashTree
	Data interface{}
}

func createStoredInstance(data interface{}) *storedInstance {
	blocks, blocksErr := helpers.GetBytes(data)
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
	IsUpdated bool
	HD        hashtree.HashTree
	Dat       map[string]*storedInstance
}

func createConcreteKeys() Keys {
	out := concreteKeys{
		HD:        nil,
		IsUpdated: true,
		Dat:       map[string]*storedInstance{},
	}

	return &out
}

// Head returns the head hashtree
func (app *concreteKeys) Head() hashtree.HashTree {
	app.rebuildHead()
	return app.HD
}

// Copy copies the Keys instance
func (app *concreteKeys) Copy() Keys {

	data := map[string]*storedInstance{}
	for keyname, oneData := range app.Dat {
		data[keyname] = createStoredInstance(oneData.Data)
	}

	out := concreteKeys{
		HD:        nil,
		IsUpdated: true,
		Dat:       data,
	}

	return &out
}

// HashTree returns the hashtree of the object at key
func (app *concreteKeys) HashTree(key string) hashtree.HashTree {
	if app.Exists(key) == 1 {
		return app.Dat[key].HT
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
	return len(app.Dat)
}

// Exists returns the amount of keys passed to Exists that exists
func (app *concreteKeys) Exists(key ...string) int {
	cpt := 0
	for _, oneKey := range key {
		if _, ok := app.Dat[oneKey]; ok {
			cpt++
		}
	}
	return cpt
}

// Retrieve retrieves data at key
func (app *concreteKeys) Retrieve(key string) interface{} {
	if app.Exists(key) == 1 {
		return app.Dat[key].Data
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
	for oneKeyname := range app.Dat {
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
	app.Dat[key] = createStoredInstance(data)

	// the data is now updated:
	app.IsUpdated = true
}

// Delete deletes the passed keys
func (app *concreteKeys) Delete(key ...string) int {
	cpt := 0
	for _, oneKey := range key {
		if _, ok := app.Dat[oneKey]; ok {
			delete(app.Dat, oneKey)
			cpt++
		}
	}

	// the data is now updated:
	app.IsUpdated = true

	//returns the amount of deleted keys:
	return cpt
}

func (app *concreteKeys) rebuildHead() {
	if !app.IsUpdated {
		return
	}

	blocks := [][]byte{
		[]byte("root"),
	}

	for keyname, ins := range app.Dat {
		blks, blksErr := helpers.GetBytes(ins.Data)
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

	app.HD = ht
	app.IsUpdated = false
}
