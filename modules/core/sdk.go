package core

import (
	applications "github.com/XMNBlockchain/datamint/applications"
	datastore "github.com/XMNBlockchain/datamint/datastore"
	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	lua "github.com/yuin/gopher-lua"
)

// ExecuteParams represents the execute func params
type ExecuteParams struct {
	DBPath      string
	NodePK      crypto.PrivKey
	InstanceID  *uuid.UUID
	RootPubKeys []crypto.PubKey
	Store       datastore.DataStore
	Context     *lua.LState
	ScriptPath  string
}

// SDKFunc represents the public SDK func of the base module
var SDKFunc = struct {
	Execute func(params ExecuteParams) applications.Node
}{
	Execute: func(params ExecuteParams) applications.Node {
		node, nodeErr := execute(params.DBPath, params.InstanceID, params.RootPubKeys, params.NodePK, params.Store, params.Context, params.ScriptPath)
		if nodeErr != nil {
			panic(nodeErr)
		}

		return node
	},
}
