package core

import (
	"testing"
)

func TestXTables_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	//create the module:
	Register(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/xtables_test.lua")
}
