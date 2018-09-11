package lists

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/XMNBlockchain/xmnsuite/helpers"
	"github.com/XMNBlockchain/xmnsuite/objects"
)

type concreteLists struct {
	isUnique bool
	objs     objects.Objects
}

func createConcreteLists(isUnique bool) Lists {
	out := concreteLists{
		isUnique: isUnique,
		objs:     objects.SDKFunc.Create(),
	}

	return &out
}

// Objects returns the objects
func (app *concreteLists) Objects() objects.Objects {
	return app.objs
}

// Copy copies the lists object
func (app *concreteLists) Copy() Lists {
	out := concreteLists{
		isUnique: app.isUnique,
		objs:     app.objs.Copy(),
	}

	return &out
}

// Add add values to a key, and returns the amount of elements added
func (app *concreteLists) Add(key string, values ...interface{}) int {
	//if the key is new:
	if app.objs.Keys().Exists(key) == 0 {
		app.objs.Save(&objects.ObjInKey{
			Key: key,
			Obj: values,
		})

		return len(values)
	}

	//retrieve the elements:
	retObj := objects.ObjInKey{
		Key: key,
		Obj: new([]interface{}),
	}

	app.objs.Retrieve(&retObj)
	ptrList := retObj.Obj.(*[]interface{})
	list := *ptrList

	//calculate the begin length:
	beginLength := len(list)

	//add the elements to the list:
	for _, oneValue := range values {
		list = append(list, oneValue)
	}

	//if the list is unique:
	if app.isUnique {
		list = helpers.MakeUnique(list...)
	}

	//save:
	retObj.Obj = list
	app.objs.Save(&retObj)
	return len(list) - beginLength
}

// Del deletes the passed values from the list, then return the amount of deleted elements
func (app *concreteLists) Del(key string, values ...interface{}) int {
	elements := app.Retrieve(key, 0, -1)
	beginLength := len(elements)
	for _, oneValue := range values {
		valueAsBytes, valueAsBytesErr := helpers.GetHash(oneValue)
		if valueAsBytesErr != nil {
			str := fmt.Sprintf("there was an error while converting a value to []byte: %s", valueAsBytesErr.Error())
			panic(errors.New(str))
		}

		for index, oneElement := range elements {
			elementAsBytes, elementAsBytesErr := helpers.GetHash(oneElement)
			if elementAsBytesErr != nil {
				str := fmt.Sprintf("there was an error while converting an element to []byte: %s", elementAsBytesErr.Error())
				panic(errors.New(str))
			}

			if bytes.Compare(valueAsBytes, elementAsBytes) == 0 {
				elements = append(elements[:index], elements[index+1:]...)
			}
		}
	}

	//replace the elements:
	app.objs.Keys().Delete(key)
	return beginLength - app.Add(key, elements...)
}

// Retrieve retrieves a subset of the stored list
func (app *concreteLists) Retrieve(key string, index int, amount int) []interface{} {
	retObj := objects.ObjInKey{
		Key: key,
		Obj: new([]interface{}),
	}

	amountRet := app.objs.Retrieve(&retObj)
	if amountRet != 1 {
		return nil
	}

	ptrList := retObj.Obj.(*[]interface{})
	list := *ptrList
	length := len(list)
	if length <= 0 {
		return []interface{}{}
	}

	if index < 0 {
		return nil
	}

	if amount == -1 {
		amount = length
	}

	if amount <= 0 {
		return nil
	}

	from := index
	if from >= length {
		return []interface{}{}
	}

	to := from + amount
	if to >= length {
		to = length
	}

	return list[from:to]
}

// Len returns the amount of elements inside the key
func (app *concreteLists) Len(key string) int {
	retObjs := objects.ObjInKey{
		Key: key,
		Obj: new([]interface{}),
	}

	amountRet := app.objs.Retrieve(&retObjs)
	if amountRet != 1 {
		return 0
	}

	list := retObjs.Obj.(*[]interface{})
	return len(*list)
}

// Union merges the elements of all the passed keys and returned them
func (app *concreteLists) Union(key ...string) []interface{} {
	out := []interface{}{}
	for _, oneKey := range key {
		elements := app.Retrieve(oneKey, 0, -1)
		if elements == nil {
			continue
		}

		for _, oneElement := range elements {
			out = append(out, oneElement)
		}
	}

	if !app.isUnique {
		return out
	}

	return helpers.MakeUnique(out...)
}

// UnionStore executes a Union, then store the results in the destination key and return the amount of elements the key holds
func (app *concreteLists) UnionStore(destination string, key ...string) int {
	elements := app.Union(key...)
	return app.Add(destination, elements...)
}

// Inter intersects the elements of all the passed keys and returned the ones that are contained in all keys
func (app *concreteLists) Inter(key ...string) []interface{} {

	type el struct {
		obj       interface{}
		occurence int
	}

	all := map[string]*el{}
	for _, oneKey := range key {
		elements := app.Retrieve(oneKey, 0, -1)
		if elements == nil {
			continue
		}

		for _, oneElement := range elements {
			elementAsBytes, elementAsBytesErr := helpers.GetBytes(oneElement)
			if elementAsBytesErr != nil {
				str := fmt.Sprintf("there was an error while converting an instance to []byte: %s", elementAsBytesErr.Error())
				panic(errors.New(str))
			}

			ha := sha256.New()
			_, haErr := ha.Write(elementAsBytes)
			if haErr != nil {
				str := fmt.Sprintf("there was an error while []byte: %s", haErr.Error())
				panic(errors.New(str))
			}

			key := hex.EncodeToString(ha.Sum(nil))
			if _, ok := all[key]; ok {
				all[key].occurence++
				continue
			}

			all[key] = &el{
				obj:       oneElement,
				occurence: 1,
			}
		}
	}

	out := []interface{}{}
	amountKeys := len(key)
	for _, oneEl := range all {
		if oneEl.occurence < amountKeys {
			continue
		}

		out = append(out, oneEl.obj)
	}

	return out
}

// InterStore executes an Inter, then store the results in the destination key and return the amount of elements the key holds
func (app *concreteLists) InterStore(destination string, key ...string) int {
	elements := app.Inter(key...)
	return app.Add(destination, elements...)
}

// Trim only keeps the elements of the list between the index and amount.  It returns the amount of elements remaining in the list
func (app *concreteLists) Trim(key string, index int, amount int) int {
	length := app.Len(key)
	if index > length {
		app.Objects().Keys().Delete(key)
		app.Add(key, []interface{}{}...)
		return 0
	}

	if index < 0 {
		index = 0
	}

	if amount == -1 {
		amount = length
	}

	stop := index + amount
	if stop > length {
		stop = length
	}

	elements := app.Retrieve(key, index, stop-index)
	app.Objects().Keys().Delete(key)
	app.Add(key, elements...)
	return len(elements)
}

// Walk will execute the WalkFn func to every element of the keys and return the list of elements that the called WalkFn calls returned
func (app *concreteLists) Walk(key string, fn WalkFn) []interface{} {
	if app.Objects().Keys().Exists(key) != 1 {
		return nil
	}

	out := []interface{}{}
	elements := app.Retrieve(key, 0, -1)
	for index, oneElement := range elements {
		ret, retErr := fn(index, oneElement)
		if retErr != nil {
			continue
		}

		out = append(out, ret)
	}

	return out
}

// WalkStore executes a Walk, then store the results in the destination key and return the amount of elements the key holds
func (app *concreteLists) WalkStore(destination string, key string, fn WalkFn) int {
	elements := app.Walk(key, fn)
	return app.Add(destination, elements...)
}
