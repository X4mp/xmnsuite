package group

import (
	"reflect"
	"testing"
)

// CompareGroupForTests compares 2 Group instances for tests
func CompareGroupForTests(t *testing.T, first Group, second Group) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid, expected: %s, returned: %s", first.ID().String(), second.ID().String())
		return
	}

	if first.Name() != second.Name() {
		t.Errorf("the Name is invalid, expected: %s, returned: %s", first.Name(), second.Name())
		return
	}
}
