package entity

import (
	"errors"
	"fmt"
	"log"

	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type service struct {
	store      datastore.DataStore
	repository Repository
}

func createService(store datastore.DataStore, repository Repository) Service {
	out := service{
		store:      store,
		repository: repository,
	}

	return &out
}

// Save saves an Entity instance
func (app *service) Save(ins Entity, rep Representation) error {
	// variables:
	met := rep.MetaData()
	name := met.Name()
	toStorableFunc := rep.ToStorable()

	// make sure the entity does not exists already:
	_, retErr := app.repository.RetrieveByID(met, ins.ID())
	if retErr == nil {
		str := fmt.Sprintf("the %s instance (ID: %s) already exists", name, ins.ID().String())
		return errors.New(str)
	}

	// sync the entity:
	if rep.HasSync() {
		syncErr := rep.Sync()(app.repository, app, ins)
		if syncErr != nil {
			str := fmt.Sprintf("there was an error while syncing the instance (Name: %s, ID: %s): %s", name, ins.ID().String(), syncErr.Error())
			return errors.New(str)
		}
	}

	// convert the instance to a storable instance:
	storable, storableErr := toStorableFunc(ins)
	if storableErr != nil {
		str := fmt.Sprintf("there was an error while converting a %s instance to a storable instance: %s", name, storableErr.Error())
		return errors.New(str)
	}

	// save the object:
	key := keynameByID(met.Keyname(), ins.ID())
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: key,
		Obj: storable,
	})

	if amountSaved != 1 {
		str := fmt.Sprintf("there was an error while saving the %s instance", name)
		return errors.New(str)
	}

	// if there is no keynames, return:
	if !rep.HasKeynames() {
		return nil
	}

	// add the instance to the sets, if any:
	keynames, keynamesErr := rep.Keynames()(ins)
	if keynamesErr != nil {
		return keynamesErr
	}

	amountAdded := app.store.Sets().AddMul(keynames, ins.ID().String())
	if amountAdded != 1 {
		// revert:
		app.store.Sets().DelMul(keynames, ins.ID().String())

		str := fmt.Sprintf("there was an error while saving the %s ID (%s) to the sets... reverting", name, ins.ID().String())
		return errors.New(str)
	}

	return nil
}

// Delete deletes an Entity instance
func (app *service) Delete(ins Entity, rep Representation) error {
	met := rep.MetaData()
	_, retEntErr := app.repository.RetrieveByID(met, ins.ID())
	if retEntErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the given Entity instance (ID: %s, Name: %s): %s", ins.ID().String(), met.Name(), retEntErr.Error())
		return errors.New(str)
	}

	// delete the instance:
	keynameByID := keynameByID(met.Keyname(), ins.ID())
	amountDel := app.store.Objects().Keys().Delete(keynameByID)
	if amountDel != 1 {
		str := fmt.Sprintf("there was an error while deleting the entity keyname (keyname: %s)", keynameByID)
		return errors.New(str)
	}

	// if there are no keynames, return
	if !rep.HasKeynames() {
		return nil
	}

	// delete the set keynames:
	setKeys, setKeysErr := rep.Keynames()(ins)
	if setKeysErr != nil {
		return setKeysErr
	}

	amountSetDel := app.store.Sets().Objects().Keys().Delete(setKeys...)
	if amountSetDel != len(setKeys) {
		log.Printf("some set keynames were not properly deleted while deleting the given Entity (Keyname: %s)", keynameByID)
	}

	return nil
}
