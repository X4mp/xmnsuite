package entity

import (
	"errors"
	"fmt"
	"log"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type repository struct {
	store datastore.DataStore
}

func createRepository(store datastore.DataStore) Repository {
	out := repository{
		store: store,
	}

	return &out
}

// RetrieveByIntersectKeynames retrieves an Entity instance by intersecting keynames
func (app *repository) RetrieveByIntersectKeynames(met MetaData, keynames []string) (Entity, error) {
	ps, psErr := app.RetrieveSetByIntersectKeynames(met, keynames, 0, 1)
	if psErr != nil {
		str := fmt.Sprintf("there was an error while retrieving an EntityPartialSet instance from keynames: %s", psErr.Error())
		return nil, errors.New(str)
	}

	if ps.TotalAmount() != 1 {
		str := fmt.Sprintf("the totalAmount of instances was expected to be 1, %d returned", ps.TotalAmount())
		return nil, errors.New(str)
	}

	ins := ps.Instances()
	return ins[0], nil
}

// RetrieveByID retrieves an Entity instance by ID
func (app *repository) RetrieveByID(met MetaData, id *uuid.UUID) (Entity, error) {
	// create the retriever criteria:
	obj := objects.ObjInKey{
		Key: keynameByID(met.Name(), id),
		Obj: met.CopyStorable(),
	}

	// retrieve the instance:
	amountRet := app.store.Objects().Retrieve(&obj)
	if amountRet != 1 {
		str := fmt.Sprintf("there was an error while retrieving the %s instance", met.Name())
		return nil, errors.New(str)
	}

	// cast the instance:
	return met.ToEntity()(app, obj.Obj)
}

// RetrieveSetByKeyname retrieves an EntityPartialSet instance by keyname
func (app *repository) RetrieveSetByKeyname(met MetaData, keyname string, index int, amount int) (PartialSet, error) {
	// retrieve the ids:
	ids := app.store.Sets().Retrieve(keyname, index, amount)

	// loop:
	entities := []Entity{}
	for _, uncastedID := range ids {
		// cast the ID:
		id, idErr := uuid.FromString(uncastedID.(string))
		if idErr != nil {
			str := fmt.Sprintf("there is at least 1 element (%s) that is not a valid UUID inside the given set keyname (%s): %s", uncastedID.(string), keyname, idErr.Error())
			return nil, errors.New(str)
		}

		// retrieve the entity:
		ins, insErr := app.RetrieveByID(met, &id)
		if insErr != nil {
			str := fmt.Sprintf("there was an error while retrieving the entity (Name: %s, ID: %s): %s", met.Name(), id.String(), insErr.Error())
			return nil, errors.New(str)
		}

		// append to list:
		entities = append(entities, ins)
	}

	// retrieve the totalAmount:
	totalAmount := app.store.Sets().Len(keyname)

	// create the partial set and return:
	out, outErr := createEntityPartialSet(entities, index, totalAmount)
	if outErr != nil {
		return nil, outErr
	}

	return out, nil
}

// RetrieveSetByIntersectKeynames retrieves an EntityPartialSet by intersecting keynames
func (app *repository) RetrieveSetByIntersectKeynames(met MetaData, keynames []string, index int, amount int) (PartialSet, error) {
	// create the destination and intersect:
	destination := fmt.Sprintf("inter:%s", strings.Join(keynames, "|"))
	amountInDest := app.store.Sets().InterStore(destination, keynames...)
	if amountInDest <= 0 {
		str := fmt.Sprintf("there is no elements that intersect the given keynames (%s)", strings.Join(keynames, ","))
		return nil, errors.New(str)
	}

	// retrieve set by keyname:
	entityPartialSet, entityPartialSetErr := app.RetrieveSetByKeyname(met, destination, index, amount)
	if entityPartialSetErr != nil {
		str := fmt.Sprintf("there was an error while retrieving an EntityPartialSet using an intersect of keynames (%s): %s", strings.Join(keynames, ","), entityPartialSetErr)
		return nil, errors.New(str)
	}

	// delete the destination key:
	amountDel := app.store.Sets().Objects().Keys().Delete(destination)
	if amountDel != 1 {
		str := fmt.Sprintf("there was a problem while deleting the interstore key: %s", destination)
		log.Printf(str)
	}

	// returns:
	return entityPartialSet, nil
}
