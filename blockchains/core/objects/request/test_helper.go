package request

import (
	"reflect"
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

// CompareRequestForTests compares 2 Request instances for tests
func CompareRequestForTests(t *testing.T, first Request, second Request) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid, expected: %s, returned: %s", first.ID().String(), second.ID().String())
		return
	}

	// compare the new entity ids:
	if !reflect.DeepEqual(first.Save().ID(), second.Save().ID()) {
		t.Errorf("the new entity ID is invalid, expected: %s, returned: %s", first.Save().ID().String(), second.Save().ID().String())
		return
	}

	// compare the user IDS:
	if !reflect.DeepEqual(first.From().ID(), second.From().ID()) {
		t.Errorf("the new from user ID is invalid, expected: %s, returned: %s", first.From().ID().String(), second.From().ID().String())
		return
	}

	// compare the keynames:
	keyname.CompareKeynameForTests(t, first.Keyname(), second.Keyname())
}
