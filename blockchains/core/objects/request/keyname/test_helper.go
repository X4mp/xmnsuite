package keyname

import (
	"reflect"
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
)

// CompareKeynameForTests compares 2 Keyname instances for tests
func CompareKeynameForTests(t *testing.T, first Keyname, second Keyname) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid, expected: %s, returned: %s", first.ID().String(), second.ID().String())
		return
	}

	if first.Name() != second.Name() {
		t.Errorf("the Name is invalid, expected: %s, returned: %s", first.Name(), second.Name())
		return
	}

	// compare the groups:
	group.CompareGroupForTests(t, first.Group(), second.Group())
}
