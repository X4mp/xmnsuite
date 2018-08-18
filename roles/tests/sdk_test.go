package tests

import (
	"testing"

	"github.com/XMNBlockchain/datamint/roles"
)

func TestCreate_Success(t *testing.T) {
	obj := roles.SDKFunc.Create()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}
}
