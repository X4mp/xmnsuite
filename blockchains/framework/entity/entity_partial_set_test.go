package entity

import (
	"math/rand"
	"testing"
)

func TestEntityPartialSet_Success(t *testing.T) {
	index := 0
	set := []Entity{
		createTestEntityForTests(),
		createTestEntityForTests(),
	}

	totalAmount := (rand.Int() % 20) + len(set) + index

	// execute:
	ps, psErr := createEntityPartialSet(set, index, totalAmount)
	if psErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", psErr.Error())
		return
	}

	// compare:
	CompareEntityPartialSetForTests(t, ps, set, index, totalAmount)
}

func TestEntityPartialSet_withEmptySet_Success(t *testing.T) {
	index := 0
	set := []Entity{}

	totalAmount := (rand.Int() % 20) + len(set) + index

	// execute:
	ps, psErr := createEntityPartialSet(set, index, totalAmount)
	if psErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", psErr.Error())
		return
	}

	// compare:
	CompareEntityPartialSetForTests(t, ps, set, index, totalAmount)
}

func TestEntityPartialSet_withIndexLowerThanZero_returnsError(t *testing.T) {
	set := []Entity{
		createTestEntityForTests(),
		createTestEntityForTests(),
	}

	// execute:
	_, psErr := createEntityPartialSet(set, -1, 200)
	if psErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}
}

func TestEntityPartialSet_WithTotalAmountTooLow_returnsError(t *testing.T) {
	set := []Entity{
		createTestEntityForTests(),
		createTestEntityForTests(),
	}

	// execute:
	_, psErr := createEntityPartialSet(set, 1, len(set))
	if psErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}
}
