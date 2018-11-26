package applications

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/datastore"
	objects "github.com/xmnservices/xmnsuite/datastore/objects"
)

/*
 * State
 */

type state struct {
	Hsh  []byte `json:"hash"`
	Hght int64  `json:"height"`
	Siz  int64  `json:"size"`
	Ver  string `json:"version"`
}

func createEmptyState(version string) (State, error) {
	//generate an app hash:
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, 0)

	//create the state:
	out := createState(version, appHash, 0, 0)
	return out, nil
}

func createState(version string, hash []byte, height int64, size int64) State {
	out := state{
		Ver:  version,
		Hsh:  hash,
		Hght: height,
		Siz:  size,
	}

	return &out
}

// Hash returns the hash
func (obj *state) Hash() []byte {
	return obj.Hsh
}

// Height returns the blockchain height
func (obj *state) Height() int64 {
	return obj.Hght
}

// Size returns the size
func (obj *state) Size() int64 {
	return obj.Siz
}

// Increment increments the database size
func (obj *state) Increment() int64 {
	obj.Siz++
	return obj.Siz
}

// Version returns the version
func (obj *state) Version() string {
	return obj.Ver
}

/*
 * Database
 */

type database struct {
	stateKey string
	states   map[string]State
	ds       datastore.StoredDataStore
}

func createDatabase(states map[string]State, ds datastore.StoredDataStore, stateKey string) Database {
	out := database{
		states:   states,
		ds:       ds,
		stateKey: stateKey,
	}

	return &out
}

func retrieveOrCreateState(currVersion string, stateKey string, ds datastore.StoredDataStore) (Database, error) {
	// retrieve the elements:
	states := ds.DataStore().Sets().Retrieve(stateKey, 0, -1)

	// if there is no state, create the first one:
	if len(states) <= 0 || states == nil {
		st, stErr := createEmptyState(currVersion)
		if stErr != nil {
			return nil, stErr
		}

		mapStatesVersion := map[string]State{
			currVersion: st,
		}

		storedState := createDatabase(mapStatesVersion, ds, stateKey)
		return storedState, nil
	}

	mapStatesVersion := map[string]State{}
	for _, oneStateHash := range states {
		stateHashAsString := hex.EncodeToString(oneStateHash.([]byte))
		stKey := fmt.Sprintf("%s:%s", stateKey, stateHashAsString)
		stRetParams := objects.ObjInKey{
			Key: stKey,
			Obj: new(state),
		}

		// retrieve the stored state:
		amount := ds.DataStore().Objects().Retrieve(&stRetParams)
		if amount != 1 {
			str := fmt.Sprintf("the state (key: %s) could not be retrieved, but is listed in the %s set", stKey, stateKey)
			return nil, errors.New(str)
		}

		newSt := stRetParams.Obj.(State)
		mapStatesVersion[newSt.Version()] = newSt
	}

	storedState := createDatabase(mapStatesVersion, ds, stateKey)
	return storedState, nil
}

// State returns the state
func (app *database) State(version string) State {
	if st, ok := app.states[version]; ok {
		return st
	}

	return nil
}

// Update updates the state
func (app *database) Update(version string) (State, error) {
	//get the current state:
	st := app.State(version)
	size := st.Size()

	// get the hash from state:
	appHash := st.Hash()

	//if the size is bigger than 0, use the store head hash:
	if size > 0 {
		appHash = app.ds.DataStore().Head().Head().Get()
	}

	//create the updated state:
	app.states[version] = createState(version, appHash, st.Height()+1, st.Size())

	// create the hash as string:
	hashAsString := hex.EncodeToString(app.states[version].Hash())

	// add the new state in the list:
	amountAdded := app.DataStore().DataStore().Sets().Add(app.stateKey, app.states[version].Hash())
	if amountAdded != 1 {
		fmt.Printf("the hash: %s, already exists in the state key: %s.  Skipping...\n", hashAsString, app.stateKey)
		return app.states[version], nil
	}

	//save the updated state:
	stKey := fmt.Sprintf("%s:%s", app.stateKey, hashAsString)
	amount := app.ds.DataStore().Objects().Save(&objects.ObjInKey{
		Key: stKey,
		Obj: app.states[version],
	})

	if amount != 1 {
		str := fmt.Sprintf("there was a problem while saving the state in the key: %s", stKey)
		return nil, errors.New(str)
	}

	// save the datastore on disk:
	saveErr := app.ds.Save()
	if saveErr != nil {
		return nil, saveErr
	}

	return app.states[version], nil
}

// DataStore returns the datastore
func (app *database) DataStore() datastore.StoredDataStore {
	return app.ds
}
