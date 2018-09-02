package core

import (
	"testing"
)

func TestXUsers_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	//create the modules:
	Register(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/xusers_test.lua")
}
