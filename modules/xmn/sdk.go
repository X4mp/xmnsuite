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
	Client      applications.Client
}

// SDKFunc represents the public SDK func of the base module
var SDKFunc = struct {
	Execute func(params ExecuteParams) applications.Node
}{
	Execute: func(params ExecuteParams) applications.Node {

		createXMNFn := func(ds datastore.DataStore, cl applications.Client) *xmn {
			if cl != nil {
				out := createXMNWithClient(ds, cl)
				return out
			}

			out := createXMN(params.Store)
			return out
		}

		// create XMN:
		xmn := createXMNFn(params.Store, params.Client)

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
