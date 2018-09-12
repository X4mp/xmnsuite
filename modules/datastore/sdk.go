package datastore

import (
	datastore "github.com/xmnservices/xmnsuite/datastore"
	lua "github.com/yuin/gopher-lua"
)

// Datastore represents the datastore module
type Datastore interface {
	Get() datastore.DataStore
	Replace(newDS datastore.DataStore)
}

// CreateParams represents the create params
type CreateParams struct {
	Context   *lua.LState
	Datastore datastore.DataStore
}

// SDKFunc represents the datastore module SDK func
var SDKFunc = struct {
	Create func(params CreateParams) Datastore
}{
	Create: func(params CreateParams) Datastore {
		out := createModule(params.Context, params.Datastore)
		return out
	},
}
