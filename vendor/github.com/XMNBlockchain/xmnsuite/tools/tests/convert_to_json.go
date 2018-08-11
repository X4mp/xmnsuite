package tests

import (
	"reflect"
	"testing"

	amino "github.com/tendermint/go-amino"
)

// ConvertToJSON converts an instance to JSON, then converts it back to instance.  Will fail if an error occurs
func ConvertToJSON(t *testing.T, v interface{}, empty interface{}, cdc *amino.Codec) {
	if v == nil {
		t.Errorf("the returned instance was expected to be valid, nil returned")
	}

	js, jsErr := cdc.MarshalJSON(v)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
	}

	newErr := cdc.UnmarshalJSON(js, empty)
	if newErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", newErr.Error())
	}

	backJS, backJSErr := cdc.MarshalJSON(empty)
	if backJSErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", backJSErr.Error())
	}

	if !reflect.DeepEqual(js, backJS) {
		t.Errorf("the json conversion (back and forth) did not succeed.  \n Expected: %v, \n Returned: %v\n", js, backJS)
	}
}
