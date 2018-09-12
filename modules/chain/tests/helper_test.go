package tests

import (
	"os"
	"testing"

	crypto "github.com/tendermint/tendermint/crypto"
	lua "github.com/yuin/gopher-lua"
)

func TestExecuteForTests_Success(t *testing.T) {
	// variables:
	dbPath := "./test_files"
	scriptPath := "lua/chain.lua"
	rootPubKeys := []crypto.PubKey{}
	defer func() {
		os.RemoveAll(dbPath)
	}()

	//create lua state:
	context := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120 * 20,
	})
	defer context.Close()

	// execute:
	node, nodeErr := ExecuteForTests(context, scriptPath, dbPath, rootPubKeys)
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return
	}
	defer node.Stop()
	node.Start()
}
