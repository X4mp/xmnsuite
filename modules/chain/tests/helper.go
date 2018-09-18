package tests

import (
	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
	crypto "github.com/xmnservices/xmnsuite/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	chain_module "github.com/xmnservices/xmnsuite/modules/chain"
	datastore_module "github.com/xmnservices/xmnsuite/modules/datastore"
	json_module "github.com/xmnservices/xmnsuite/modules/json"
	lua "github.com/yuin/gopher-lua"
)

// ExecuteForTests executes the chain for tests
func ExecuteForTests(context *lua.LState, scriptPath string, dbPath string, rootPubKeys []crypto.PublicKey) (applications.Node, error) {
	// variables:
	instanceID := uuid.NewV4()
	nodePK := ed25519.GenPrivKey()
	ds := datastore.SDKFunc.Create()

	// preload JSON:
	json_module.SDKFunc.Create(json_module.CreateParams{
		Context: context,
	})

	// preload datastore:
	dsMod := datastore_module.SDKFunc.Create(datastore_module.CreateParams{
		Context:   context,
		Datastore: ds,
	})

	// create module:
	module := chain_module.SDKFunc.Create(chain_module.CreateParams{
		Context:     context,
		DBPath:      dbPath,
		ID:          &instanceID,
		RootPubKeys: rootPubKeys,
		NodePK:      nodePK,
		Datastore:   dsMod,
	})

	//execute the script:
	doFileErr := context.DoFile(scriptPath)
	if doFileErr != nil {
		return nil, doFileErr
	}

	// spawn:
	node, nodeErr := module.Spawn()
	if nodeErr != nil {
		return nil, nodeErr
	}

	return node, nil
}
