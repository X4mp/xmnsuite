package modules

import (
	"testing"
)

func TestXCrypto_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	//create the module:
	CreateXCrypto(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/xcrypto_test.lua")
}
