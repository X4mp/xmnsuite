package tendermint

import (
	"encoding/binary"
	"encoding/json"

	datastore "github.com/XMNBlockchain/datamint/datastore"
	hashtree "github.com/XMNBlockchain/datamint/hashtree"
	keys "github.com/XMNBlockchain/datamint/keys"
	router "github.com/XMNBlockchain/datamint/router"
	code "github.com/tendermint/tendermint/abci/example/code"
	types "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

/*
 * State
 */

type state struct {
	Hash          []byte `json:"hash"`
	KeynamePrefix string `json:"keyname_prefix"`
	Height        int    `json:"height"`
	Size          int64  `json:"size"`
}

func createEmptyState(appHash []byte, keynamePrefix string) (*state, error) {
	out := createState(appHash, keynamePrefix, 0, 0)
	return out, nil
}

func createState(hash []byte, keynamePrefix string, height int, size int64) *state {
	out := state{
		Hash:          hash,
		KeynamePrefix: keynamePrefix,
		Height:        height,
		Size:          size,
	}

	return &out
}

// GetHash returns the hash
func (obj *state) GetHash() []byte {
	return obj.Hash
}

// GetHeight returns the height
func (obj *state) GetHeight() int {
	return obj.Height
}

// GetSize returns the size
func (obj *state) GetSize() int64 {
	return obj.Size
}

// GetKeynamePrefix returns the keyname prefix
func (obj *state) GetKeynamePrefix() string {
	return obj.KeynamePrefix
}

// Increment increments the database size
func (obj *state) Increment() int64 {
	obj.Size++
	return obj.Size
}

/*
 * StoredState
 */

type storedState struct {
	st *state
	k  keys.Keys
}

func createStoredState(st *state, k keys.Keys) *storedState {
	out := storedState{
		st: st,
		k:  k,
	}

	return &out
}

// GetState returns the state
func (obj *storedState) GetState() *state {
	return obj.st
}

// Set sets data to the keyname in the database
func (obj *storedState) Set(keyname string, data []byte) {
	obj.k.Save(keyname, data)
}

/*
 * Application
 */

type abciApplication struct {
	types.BaseApplication
	rt       router.Router
	store    datastore.DataStore
	state    *storedState
	stateKey string
}

func createABCIApplication(stateKey string, keynamePrefix string, appHash []byte, k keys.Keys, store datastore.DataStore, rt router.Router) (*abciApplication, error) {
	storedState, storedStateErr := func(stateKey string, k keys.Keys) (*storedState, error) {
		retState := k.Retrieve(stateKey)
		if stateBytes, ok := retState.([]byte); ok {
			if len(stateBytes) > 0 {
				st := new(state)
				jsErr := json.Unmarshal(stateBytes, st)
				if jsErr != nil {
					return nil, jsErr
				}

				storedState := createStoredState(st, k)
				return storedState, nil
			}
		}

		st, stErr := createEmptyState(appHash, keynamePrefix)
		if stErr != nil {
			return nil, stErr
		}

		storedState := createStoredState(st, k)
		return storedState, nil
	}(stateKey, k)

	if storedStateErr != nil {
		return nil, storedStateErr
	}

	out := abciApplication{
		state:    storedState,
		rt:       rt,
		store:    store,
		stateKey: stateKey,
	}

	return &out, nil
}

// Info outputs information related to the abciApplication state
func (app *abciApplication) Info(req types.RequestInfo) types.ResponseInfo {
	out := struct {
		Size int64 `json:"size"`
	}{
		Size: app.state.GetState().GetSize(),
	}

	js, jsErr := cdc.MarshalJSON(out)
	if jsErr != nil {
		panic(js)
	}

	return types.ResponseInfo{Data: string(js)}
}

// DeliverTx delivers a transaction to the abciApplication
func (app *abciApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {

	//execute the transaction:
	resp := app.rt.Transact(router.SDKFunc.CreateRequest(router.CreateRequestParams{
		JSData: tx,
	}))

	//if the transaction is un-authorized:
	if !resp.IsAuthorized() {
		return types.ResponseDeliverTx{Code: code.CodeTypeUnauthorized, Log: resp.Log(), GasUsed: resp.GazUsed()}
	}

	//if the transaction was not successful, flag the encoding error:
	if !resp.IsSuccess() {
		return types.ResponseDeliverTx{Code: code.CodeTypeEncodingError, Log: resp.Log(), GasUsed: resp.GazUsed()}
	}

	//create the tags:
	tags := resp.Tags()
	tagPairs := []cmn.KVPair{}
	for key, value := range tags {
		tagPairs = append(tagPairs, cmn.KVPair{
			Key:   []byte(key),
			Value: value,
		})
	}

	//the transaction was successful:
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Log: resp.Log(), GasUsed: resp.GazUsed(), Tags: tagPairs}
}

// CheckTx verifies that a transaction is valid before it gets executed
func (app *abciApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	//execute the transaction:
	resp := app.rt.CheckTrx(router.SDKFunc.CreateTrxChkRequest(router.CreateTrxChkRequestParams{
		JSData: tx,
	}))

	if !resp.CanBeAuthorized() {
		return types.ResponseCheckTx{Code: code.CodeTypeUnauthorized, Log: resp.Log(), GasWanted: resp.GazWanted()}
	}

	if !resp.CanBeExecuted() {
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError, Log: resp.Log(), GasWanted: resp.GazWanted()}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, Log: resp.Log(), GasWanted: resp.GazWanted()}
}

// Commit commits the blockchain
func (app *abciApplication) Commit() types.ResponseCommit {

	//get the current state:
	st := app.state.GetState()
	size := app.state.GetState().GetSize()

	//generate an app hash:
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, size)

	//if the size is bigger than 0, add the store head:
	if size > 0 {
		head := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
			Blocks: [][]byte{
				appHash,
				app.store.Head().Head().Get(),
			},
		})

		appHash = head.Head().Get()
	}

	//create the updated state:
	newSt := createState(appHash, st.GetKeynamePrefix(), st.GetHeight()+1, st.GetSize())
	stJS, stJSErr := json.Marshal(newSt)
	if stJSErr != nil {
		panic(stJSErr)
	}

	//set the updated state:
	app.state.Set(app.stateKey, stJS)

	//return the response:
	return types.ResponseCommit{Data: appHash}
}

// Query executes a query on the abciApplication
func (app *abciApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	resp := app.rt.Query(router.SDKFunc.CreateRequest(router.CreateRequestParams{
		JSData: reqQuery.GetData(),
	}))

	js, jsErr := cdc.MarshalJSON(resp)
	if jsErr != nil {
		panic(jsErr)
	}

	return types.ResponseQuery{
		Value: js,
		Log:   "success",
	}
}
