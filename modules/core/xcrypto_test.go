package core

import (
	"testing"

	"github.com/XMNBlockchain/datamint/datastore"
)

func TestXCrypto_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	//create the module:
	createCore(l, datastore.SDKFunc.Create())

	//execute the chunk:
	executeChunkForTests(l, "lua/xcrypto_test.lua")
}
