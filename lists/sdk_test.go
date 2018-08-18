package lists

import (
	"testing"
)

func TestCreateList_Success(t *testing.T) {

	//variables:
	element := "this-is-an-element"
	key := "this-is-a-key"

	obj := SDKFunc.CreateList()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}

	retAmount := obj.Add(key, element, element, element, element, element)
	if retAmount != 5 {
		t.Errorf("the returned amount was expected to be 5, %d returned", retAmount)
		return
	}
}

func TestCreateSet_Success(t *testing.T) {

	//variables:
	element := "this-is-an-element"
	key := "this-is-a-key"

	obj := SDKFunc.CreateSet()
	if obj == nil {
		t.Errorf("the created object was not expected to be nil")
		return
	}

	retAmount := obj.Add(key, element, element, element, element, element)
	if retAmount != 5 {
		t.Errorf("the returned amount was expected to be 5, %d returned", retAmount)
		return
	}
}
