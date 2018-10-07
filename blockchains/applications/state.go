package applications

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	objects "github.com/xmnservices/xmnsuite/datastore/objects"
)

/*
 * State
 */

type state struct {
	Hsh  []byte `json:"hash"`
	Hght int64  `json:"height"`
	Siz  int64  `json:"size"`
}

func createEmptyState() (*state, error) {
	//generate an app hash:
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, 0)

	//create the state:
	out := createState(appHash, 0, 0)
	return out, nil
}

func createState(hash []byte, height int64, size int64) *state {
	out := state{
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

/*
 * StoredState
 */

type storedState struct {
	states map[string]*state
	objs   objects.Objects
}

func retrieveOrCreateState(currVersion string, stateKey string, objs objects.Objects) (*storedState, error) {
	// create the retrieval params:
	statesVersionInKey := objects.ObjInKey{
		Key: stateKey,
		Obj: map[string][]byte{},
	}

	// retrieve the stored state:
	amount := objs.Retrieve(&statesVersionInKey)

	//if it the stored state cannot be retrieved.  Probably menas that it has never been saved:
	if amount != 1 {
		st, stErr := createEmptyState()
		if stErr != nil {
			return nil, stErr
		}

		mapStatesVersion := map[string]*state{
			currVersion: st,
		}

		storedState := createStoredState(mapStatesVersion, objs)
		return storedState, nil
	}

	// cast the returned element to a version state map:
	if mapStatesAsBytesVersion, ok := statesVersionInKey.Obj.(map[string][]byte); ok {
		mapStatesVersion := map[string]*state{}
		for version, stateBytes := range mapStatesAsBytesVersion {
			if len(stateBytes) > 0 {
				st := new(state)
				jsErr := json.Unmarshal(stateBytes, st)
				if jsErr != nil {
					return nil, jsErr
				}

				mapStatesVersion[version] = st
				continue
			}

			st, stErr := createEmptyState()
			if stErr != nil {
				return nil, stErr
			}

			mapStatesVersion[version] = st
		}

		storedState := createStoredState(mapStatesVersion, objs)
		return storedState, nil
	}

	//the returned element is invalid:
	str := fmt.Sprintf("the element at the given stateKey (%s) is not valid", stateKey)
	return nil, errors.New(str)
}

func createStoredState(states map[string]*state, objs objects.Objects) *storedState {
	out := storedState{
		states: states,
		objs:   objs,
	}

	return &out
}

// State returns the state
func (app *storedState) State(version string) *state {
	if st, ok := app.states[version]; ok {
		return st
	}

	return nil
}

// Set sets data to the keyname in the database, for a given version
func (app *storedState) Set(keyname string, obj interface{}) int {
	return app.objs.Save(&objects.ObjInKey{
		Key: keyname,
		Obj: obj,
	})
}
