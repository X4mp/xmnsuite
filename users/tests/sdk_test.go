package tests

import (
	"testing"

	"github.com/XMNBlockchain/xmnsuite/users"
)

func TestCreate_Success(t *testing.T) {
	obj := users.SDKFunc.Create()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}
}
