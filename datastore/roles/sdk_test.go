package roles

import (
	"testing"
)

func TestCreate_Success(t *testing.T) {
	obj := SDKFunc.Create()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}
}
