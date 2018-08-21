package modules

import (
	"testing"
)

func TestXRoles_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	//create the module:
	CreateXRoles(l)
	CreateXCrypto(l)
	CreateXUsers(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/xroles_test.lua")
}
