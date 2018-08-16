package lists

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/XMNBlockchain/datamint"
	"github.com/XMNBlockchain/datamint/objects"
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

	//if the list is not unique:
	if !app.isUnique {
		for _, oneValue := range values {
			list = append(list, oneValue)
		}

		retObj.Obj = list
		app.objs.Save(&retObj)
		return len(values)
	}

	uniqueValues := []interface{}{}
	for _, oneNewValue := range values {
		isUnique := true
		newBytes, newBytesErr := datamint.GetBytes(oneNewValue)
		if newBytesErr != nil {
			str := fmt.Sprintf("there was an error while converting a new value interface{} tp []byte: %s", newBytesErr.Error())
			panic(errors.New(str))
		}

		for _, oneExistingElement := range list {
			existingBytes, existingBytesErr := datamint.GetBytes(oneExistingElement)
			if existingBytesErr != nil {
				str := fmt.Sprintf("there was an error while converting an existing value interface{} tp []byte: %s", existingBytesErr.Error())
				panic(errors.New(str))
			}

			if bytes.Compare(existingBytes, newBytes) == 0 {
				isUnique = false
				break
			}
		}

		if isUnique {
			uniqueValues = append(uniqueValues, oneNewValue)
			continue
		}
	}

	//add the unique values to the list:
	for _, oneUniqueValues := range uniqueValues {
		list = append(list, oneUniqueValues)
	}

	//save:
	retObj.Obj = list
	app.objs.Save(&retObj)

	//return the amount of element saved:
	return len(uniqueValues)
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

	unique := []interface{}{}
	for _, onElement := range out {
		isUnique := true
		oneElementAsBytes, oneElementAsBytesErr := datamint.GetBytes(onElement)
		if oneElementAsBytesErr != nil {
			str := fmt.Sprintf("there was an error while converting an existing element to []byte: %s", oneElementAsBytesErr.Error())
			panic(errors.New(str))
		}

		for _, oneUnique := range unique {
			oneUniqueAsBytes, oneUniqueAsBytesErr := datamint.GetBytes(oneUnique)
			if oneUniqueAsBytesErr != nil {
				str := fmt.Sprintf("there was an error while converting a unique element to []byte: %s", oneUniqueAsBytesErr.Error())
				panic(errors.New(str))
			}

			if bytes.Compare(oneElementAsBytes, oneUniqueAsBytes) == 0 {
				isUnique = false
			}
		}

		if isUnique {
			unique = append(unique, onElement)
		}
	}

	return unique
}

// UnionStore executes a Union, then store the results in the destination key and return the amount of elements the key holds
func (app *concreteLists) UnionStore(destination string, key ...string) int {
	elements := app.Union(key...)
	return app.Add(destination, elements...)
}
