package information

import (
	"math/rand"
	"testing"

	uuid "github.com/satori/go.uuid"
)

// CreateInformationWithConcensusNeededForTests creates a Information for tests
func CreateInformationWithConcensusNeededForTests(concensusNeeded int) Information {
	id := uuid.NewV4()
	gazPricePerKb := rand.Int() % 30
	maxAmountOfValidators := rand.Int() % 20
	out, outErr := createInformation(&id, concensusNeeded, gazPricePerKb, maxAmountOfValidators)
	if outErr != nil {
		panic(outErr)
	}

	return out
}

// CompareInformationForTests compares Information instances for tests
func CompareInformationForTests(t *testing.T, first Information, second Information) {
	if first.GazPricePerKb() != second.GazPricePerKb() {
		t.Errorf("the returned gaz price is invalid.  Expected: %d, Returned: %d", first.GazPricePerKb(), second.GazPricePerKb())
		return
	}

	if first.MaxAmountOfValidators() != second.MaxAmountOfValidators() {
		t.Errorf("the returned maximum amount of validatoirs is invalid.  Expected: %d, Returned: %d", first.MaxAmountOfValidators(), second.MaxAmountOfValidators())
		return
	}
}
