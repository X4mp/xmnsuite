package tests

import (
	"reflect"
	"testing"

	amino "github.com/tendermint/go-amino"
)

// ConvertToBinary converts an instance to data, then converts it back to instance.  Will fail if an error occurs
func ConvertToBinary(t *testing.T, v interface{}, empty interface{}, cdc *amino.Codec) {
	if v == nil {
		t.Errorf("the returned instance was expected to be valid, nil returned")
	}

	data, dataErr := cdc.MarshalBinary(v)
	if dataErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", dataErr.Error())
	}

	newErr := cdc.UnmarshalBinary(data, empty)
	if newErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", newErr.Error())
	}

	backData, backDataErr := cdc.MarshalBinary(empty)
	if backDataErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", backDataErr.Error())
	}

	if !reflect.DeepEqual(data, backData) {
		t.Errorf("the binary conversion (back and forth) did not succeed.  \n Expected: %v, \n Returned: %v\n", data, backData)
	}
}
