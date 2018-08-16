package tests

import (
	"testing"

	"github.com/XMNBlockchain/datamint/keys"
)

func TestCreate_Success(t *testing.T) {
	obj := keys.SDKFunc.Create()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}
}
