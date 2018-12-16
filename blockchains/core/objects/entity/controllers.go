package entity

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"

	uuid "github.com/satori/go.uuid"
	crypto "github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

type controllers struct {
	met                      MetaData
	rep                      Representation
	defaultAmountOfElements  int
	gazPricePerKb            int
	overwriteIfAlreadyExists bool
	routerRoleKey            string
}

func createControllers(
	met MetaData,
	rep Representation,
	defaultAmountOfElements int,
	gazPricePerKb int,
	overwriteIfAlreadyExists bool,
	routerRoleKey string,
) Controllers {
	out := controllers{
		met: met,
		rep: rep,
		defaultAmountOfElements:  defaultAmountOfElements,
		gazPricePerKb:            gazPricePerKb,
		overwriteIfAlreadyExists: overwriteIfAlreadyExists,
		routerRoleKey:            routerRoleKey,
	}

	return &out
}

// Save saves an entity
func (app *controllers) Save() routers.SaveTransactionFn {
	out := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {
		// create the repository:
		repository := createRepository(store)
		service := createService(store, repository)

		// convert the data to an entity:
		ins, insErr := app.met.ToEntity()(repository, data)
		if insErr != nil {
			return nil, insErr
		}

		// if we do not overwrite:
		if !app.overwriteIfAlreadyExists {
			// make sure the entity does not already exists:
			_, retErr := repository.RetrieveByID(app.met, ins.ID())
			if retErr == nil {
				str := fmt.Sprintf("the entity instance (Name: %s, ID: %s) already exists", app.met.Name(), ins.ID().String())
				return nil, errors.New(str)
			}
		}

		// save the entity instance:
		saveErr := service.Save(ins, app.rep)
		if saveErr != nil {
			return nil, saveErr
		}

		// convert the instance to json:
		jsData, jsDataErr := cdc.MarshalJSON(ins)
		if jsDataErr != nil {
			return nil, jsDataErr
		}

		// create the gaz price:
		gazUsed := int(unsafe.Sizeof(jsData)) * app.gazPricePerKb

		// create the element path:
		elementPath := fmt.Sprintf("%sid:%s", path, ins.ID().String())

		// add the owner of the resource in the role key, so that it can delete the resource later:
		if store.Roles().EnableWriteAccess(app.routerRoleKey, elementPath) != 1 {
			str := fmt.Sprintf("there was an error while enabling user (pubKey: %s) write access to the resource (path: %s) on the role (key: %s)", from.String(), elementPath, app.routerRoleKey)
			return nil, errors.New(str)
		}

		// return the response:
		resp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
			Code:    routers.IsSuccessful,
			Log:     "success",
			GazUsed: int64(gazUsed),
			Tags: map[string][]byte{
				elementPath: jsData,
			},
		})

		return resp, nil
	}

	return out
}

/*
 * Delete
 * Expected params:
 *      <id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>
 */

// Delete deletes an entity
func (app *controllers) Delete() routers.DeleteTransactionFn {
	out := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.TransactionResponse, error) {
		// create the repository:
		repository := createRepository(store)
		service := createService(store, repository)

		// fetch the ID:
		id, idErr := uuid.FromString(fetchFromParams(params, "id"))
		if idErr != nil {
			return nil, idErr
		}

		// retrieve the entity:
		retIns, retInsErr := repository.RetrieveByID(app.met, &id)
		if retInsErr != nil {
			return nil, retInsErr
		}

		// delete the entity instance:
		delErr := service.Delete(retIns, app.rep)
		if delErr != nil {
			return nil, delErr
		}

		// convert the instance to json:
		jsData, jsDataErr := cdc.MarshalJSON(retIns)
		if jsDataErr != nil {
			return nil, jsDataErr
		}

		// create the gaz price:
		gazUsed := int(unsafe.Sizeof(jsData)) * app.gazPricePerKb

		// disable the write access:
		if store.Roles().DisableWriteAccess(app.routerRoleKey, path) != 1 {
			str := fmt.Sprintf("there was an error while disabling user (pubKey: %s) write access to the resource (path: %s) on the role (key: %s)", from.String(), path, app.routerRoleKey)
			return nil, errors.New(str)
		}

		// return the response:
		resp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
			Code:    routers.IsSuccessful,
			Log:     "success",
			GazUsed: int64(gazUsed),
			Tags: map[string][]byte{
				path: jsData,
			},
		})

		return resp, nil
	}

	return out
}

/*
 * RetrieveByID
 * Expected params:
 *      <id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>
 */
// RetrieveByID retrieves an entity by its ID
func (app *controllers) RetrieveByID() routers.QueryFn {
	out := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
		// create the repository:
		repository := createRepository(store)

		// retrieve the entity partial set:
		idAsString := fetchFromParams(params, "id")
		id, idErr := uuid.FromString(idAsString)
		if idErr != nil {
			return nil, idErr
		}

		retIns, retInsErr := repository.RetrieveByID(app.met, &id)
		if retInsErr != nil {
			return nil, retInsErr
		}

		// convert the entity to data:
		jsData, jsDataErr := cdc.MarshalJSON(retIns)
		if jsDataErr != nil {
			return nil, jsDataErr
		}

		// return the response:
		resp := routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
			Code:  routers.IsSuccessful,
			Log:   "success",
			Key:   path,
			Value: jsData,
		})

		return resp, nil
	}

	return out
}

/*
 * RetrieveByIntersectKeynames
 * Expected params:
 *      <keynames|[a-z,]+>
 */
// RetrieveByIntersectKeynames retrieves an entity by keynames intersect
func (app *controllers) RetrieveByIntersectKeynames() routers.QueryFn {
	out := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
		// create the repository:
		repository := createRepository(store)

		// retrieve the entity partial set:
		keynames := strings.Split(fetchFromParams(params, "keynames"), ",")
		retIns, retInsErr := repository.RetrieveByIntersectKeynames(app.met, keynames)
		if retInsErr != nil {
			return nil, retInsErr
		}

		// convert the entity to data:
		jsData, jsDataErr := cdc.MarshalJSON(retIns)
		if jsDataErr != nil {
			return nil, jsDataErr
		}

		// return the response:
		resp := routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
			Code:  routers.IsSuccessful,
			Log:   "success",
			Key:   path,
			Value: jsData,
		})

		return resp, nil
	}

	return out
}

/*
 * RetrieveSetByIntersectKeynames
 * Expected params:
 *      <keynames|[a-z-,]+>
 */
// RetrieveSetByIntersectKeynames retrieves an entity partial set by keynames intersect
func (app *controllers) RetrieveSetByIntersectKeynames() routers.QueryFn {
	out := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
		// create the repository:
		repository := createRepository(store)

		// retrieve the entity partial set:
		index := fetchIntFromParams(params, "index", 0)
		amount := fetchIntFromParams(params, "amount", app.defaultAmountOfElements)
		keynames := strings.Split(fetchFromParams(params, "keynames"), ",")
		retIns, retInsErr := repository.RetrieveSetByIntersectKeynames(app.met, keynames, index, amount)
		if retInsErr != nil {
			return nil, retInsErr
		}

		// convert the entity to data:
		jsData, jsDataErr := cdc.MarshalJSON(retIns)
		if jsDataErr != nil {
			return nil, jsDataErr
		}

		// return the response:
		resp := routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
			Code:  routers.IsSuccessful,
			Log:   "success",
			Key:   path,
			Value: jsData,
		})

		return resp, nil
	}

	return out
}

/*
 * RetrieveSetByKeyname
 * Expected params:
 *      <keyname|[a-z-]+>
 */
// RetrieveSetByKeyname retrieves an entity partial set by keyname
func (app *controllers) RetrieveSetByKeyname() routers.QueryFn {
	out := func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
		// create the repository:
		repository := createRepository(store)

		// retrieve the entity partial set:
		index := fetchIntFromParams(params, "index", 0)
		amount := fetchIntFromParams(params, "amount", app.defaultAmountOfElements)
		keyname := fetchFromParams(params, "keyname")
		retIns, retInsErr := repository.RetrieveSetByKeyname(app.met, keyname, index, amount)
		if retInsErr != nil {
			return nil, retInsErr
		}

		// convert the entity to data:
		jsData, jsDataErr := cdc.MarshalJSON(retIns)
		if jsDataErr != nil {
			return nil, jsDataErr
		}

		// return the response:
		resp := routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
			Code:  routers.IsSuccessful,
			Log:   "success",
			Key:   path,
			Value: jsData,
		})

		return resp, nil
	}

	return out
}
