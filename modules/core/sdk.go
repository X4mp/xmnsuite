package core

import (
	applications "github.com/XMNBlockchain/datamint/applications"
	datastore "github.com/XMNBlockchain/datamint/datastore"
	crypto "github.com/tendermint/tendermint/crypto"
	lua "github.com/yuin/gopher-lua"
)

// ExecuteParams represents the execute func params
type ExecuteParams struct {
	DBPath     string
	NodePK     crypto.PrivKey
	Store      datastore.DataStore
	Context    *lua.LState
	ScriptPath string
}

// SDKFunc represents the public SDK func of the base module
var SDKFunc = struct {
	Execute func(params ExecuteParams) applications.Node
}{
	Execute: func(params ExecuteParams) applications.Node {
		node, nodeErr := execute(params.DBPath, params.NodePK, params.Store, params.Context, params.ScriptPath)
		if nodeErr != nil {
			panic(nodeErr)
		}

		return node
	},
}
