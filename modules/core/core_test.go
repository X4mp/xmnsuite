package core

import (
	"os"
	"testing"

	datastore "github.com/XMNBlockchain/datamint/datastore"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

func TestCore_Success(t *testing.T) {

	// variables:
	dbPath := "./test_files"
	defer func() {
		os.RemoveAll(dbPath)
	}()

	nodePK := ed25519.GenPrivKey()

	//create lua state:
	l := createLuaState()
	defer l.Close()

	// create the datastore:
	ds := datastore.SDKFunc.Create()

	// execute:
	node, nodeErr := execute(dbPath, nodePK, ds, l, "lua/core_test.lua")
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return
	}

	// retrieve the client:
	_, clientErr := node.GetClient()
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}
}
