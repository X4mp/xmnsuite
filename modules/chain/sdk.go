package chain

import (
	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	applications "github.com/xmnservices/xmnsuite/applications"
	datastore_module "github.com/xmnservices/xmnsuite/modules/datastore"
	lua "github.com/yuin/gopher-lua"
)

// Chain represents the chain module
type Chain interface {
	Spawn() (applications.Node, error)
}

// CreateParams represents the create params
type CreateParams struct {
	Context     *lua.LState
	DBPath      string
	ID          *uuid.UUID
	RootPubKeys []crypto.PubKey
	NodePK      crypto.PrivKey
	Datastore   datastore_module.Datastore
}

// SDKFunc represents the chain module SDK func
var SDKFunc = struct {
	Create func(params CreateParams) Chain
}{
	Create: func(params CreateParams) Chain {
		out := createModule(params.Context, params.DBPath, params.ID, params.RootPubKeys, params.NodePK, params.Datastore)
		return out
	},
}