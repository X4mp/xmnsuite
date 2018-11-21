package vote

import (
	"reflect"
	"testing"
)

// CompareVoteForTests compares 2 Vote instances for tests
func CompareVoteForTests(t *testing.T, first Vote, second Vote) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid, expected: %s, returned: %s", first.ID().String(), second.ID().String())
		return
	}

	// compare the request IDs:
	if !reflect.DeepEqual(first.Request().ID(), second.Request().ID()) {
		t.Errorf("the request ID is invalid, expected: %s, returned: %s", first.Request().ID().String(), second.Request().ID().String())
		return
	}

	// compare the voter IDs:
	if !reflect.DeepEqual(first.Voter().ID(), second.Voter().ID()) {
		t.Errorf("the voter ID is invalid, expected: %s, returned: %s", first.Voter().ID().String(), second.Voter().ID().String())
		return
	}

	if first.IsApproved() != second.IsApproved() {
		t.Errorf("the isApproved is invalid")
	}
}
