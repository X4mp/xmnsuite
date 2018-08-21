package tests

import (
	"testing"

	"github.com/XMNBlockchain/datamint/lists"
)

func TestCreateList_Success(t *testing.T) {
	obj := lists.SDKFunc.CreateList()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}
}

func TestCreateSet_Success(t *testing.T) {
	obj := lists.SDKFunc.CreateSet()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}
}