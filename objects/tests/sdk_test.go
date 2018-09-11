package tests

import (
	"testing"

	"github.com/XMNBlockchain/xmnsuite/objects"
)

func TestCreate_Success(t *testing.T) {
	obj := objects.SDKFunc.Create()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}
}
