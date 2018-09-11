package xmn

import (
	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	applications "github.com/xmnservices/xmnsuite/applications"
	datastore "github.com/xmnservices/xmnsuite/datastore"
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

		// create XMN:
		xmn := createXMN(params.Store)

		// register:
		xmn.register(params.Context)

		// execute:
		node, nodeErr := xmn.execute(params.Context, params.DBPath, params.InstanceID, params.RootPubKeys, params.NodePK, params.ScriptPath)
		if nodeErr != nil {
			panic(nodeErr)
		}

		return node
	},
}
