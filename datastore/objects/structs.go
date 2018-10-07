package objects

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/helpers"
	"github.com/xmnservices/xmnsuite/datastore/keys"
)

/*
 * Concrete Objects
 */

type concreteObjects struct {
	K keys.Keys
}

func createObjects() Objects {
	out := concreteObjects{
		K: keys.SDKFunc.Create(),
	}

	return &out
}

// Keys returns the keys instance
func (app *concreteObjects) Keys() keys.Keys {
	return app.K
}

// Copy copies the objects instance
func (app *concreteObjects) Copy() Objects {
	out := concreteObjects{
		K: app.Keys().Copy(),
	}

	return &out
}

// Retrieve populates the Obj pointers in the passed ObjInKey instances.  Returns the amount of instances retrieved
func (app *concreteObjects) Retrieve(objs ...*ObjInKey) int {
	cpt := 0
	for index, oneObj := range objs {
		if app.K.Exists(oneObj.Key) == 1 {
			data := app.K.Retrieve(oneObj.Key)
			marErr := helpers.Marshal(data.([]byte), objs[index].Obj)
			if marErr != nil {
				str := fmt.Sprintf("there was an error while unmarshalling data to the given pointer (index: %d): %s", index, marErr.Error())
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
		bytes, bytesErr := helpers.GetBytes(oneObj.Obj)
		if bytesErr != nil {
			str := fmt.Sprintf("there was an error while converting an instance to []byte: %s", bytesErr.Error())
			panic(errors.New(str))
		}

		app.K.Save(oneObj.Key, bytes)
		cpt++
	}

	//returns the amount of saved keys:
	return cpt
}
