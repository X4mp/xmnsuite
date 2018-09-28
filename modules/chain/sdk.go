package chain

import (
	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
	crypto "github.com/xmnservices/xmnsuite/crypto"
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
	Port        int
	ID          *uuid.UUID
	RootPubKeys []crypto.PublicKey
	NodePK      tcrypto.PrivKey
	Datastore   datastore_module.Datastore
}

// SDKFunc represents the chain module SDK func
var SDKFunc = struct {
	Create func(params CreateParams) Chain
}{
	Create: func(params CreateParams) Chain {
		out := createModule(params.Context, params.DBPath, params.Port, params.ID, params.RootPubKeys, params.NodePK, params.Datastore)
		return out
	},
}
