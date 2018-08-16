package objects

import (
	"errors"
	"fmt"

	"github.com/XMNBlockchain/datamint/keys"
)

/*
 * Concrete Objects
 */

type concreteObjects struct {
	keys keys.Keys
}

func createObjects() Objects {
	out := concreteObjects{
		keys: keys.SDKFunc.Create(),
	}

	return &out
}

// Keys returns the keys instance
func (app *concreteObjects) Keys() keys.Keys {
	return app.keys
}

// Retrieve populates the Obj pointers in the passed ObjInKey instances.  Returns the amount of instances retrieved
func (app *concreteObjects) Retrieve(objs ...*ObjInKey) int {
	cpt := 0
	for index, oneObj := range objs {
		if app.keys.Exists(oneObj.Key) == 1 {
			js := app.keys.Retrieve(oneObj.Key)
			jsErr := cdc.UnmarshalJSON(js, objs[index].Obj)
			if jsErr != nil {
				str := fmt.Sprintf("there was an error while unmarshalling json data to the given pointer (index: %d): %s", index, jsErr.Error())
				panic(errors.New(str))
			}

			cpt++
		}
	}

	return cpt
}

// Save saves the Obj at key as explained in the passed ObjInKey instances
func (app *concreteObjects) Save(objs ...*ObjInKey) int {
	cpt := 0
	for _, oneObj := range objs {
		js, jsErr := cdc.MarshalJSON(oneObj.Obj)
		if jsErr != nil {
			str := fmt.Sprintf("there was an error while marshalling an instance to json: %s", jsErr.Error())
			panic(errors.New(str))
		}

		app.keys.Save(oneObj.Key, js)
		cpt++
	}

	//returns the amount of saved keys:
	return cpt
}
