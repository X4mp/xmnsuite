package tests

import (
	"reflect"
	"testing"

	hashtree "github.com/XMNBlockchain/xmnsuite/hashtree"
)

func TestAll_Success(t *testing.T) {
	ht, htErr := hashtree.SDKFunc.CreateHashTree(hashtree.CreateHashTreeParams{
		Blocks: [][]byte{
			[]byte("this"),
			[]byte("is"),
			[]byte("some"),
		},
	})

	if htErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", htErr.Error())
		return
	}

	jsCompact := hashtree.SDKFunc.CreateJSONCompact(ht.Compact())
	retCompact, retCompactErr := hashtree.SDKFunc.CreateCompactFromJSON(jsCompact)
	if retCompactErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retCompactErr.Error())
		return
	}

	retHT := retCompact.Leaves().HashTree()

	if !reflect.DeepEqual(ht, retHT) {
		t.Errorf("the returned hashtree is invalid")
		return
	}
}
