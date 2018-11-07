package genesis

import (
	"math/rand"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
)

// CreateGenesisForTests creates a Genesis for tests
func CreateGenesisForTests() Genesis {
	id := uuid.NewV4()
	gazPricePerKb := rand.Int() % 30
	maxAmountOfValidators := rand.Int() % 20
	dep := deposit.CreateDepositForTests()
	out := createGenesis(&id, gazPricePerKb, maxAmountOfValidators, dep)
	return out
}

// CompareGenesisForTests compares Genesis instances for tests
func CompareGenesisForTests(t *testing.T, first Genesis, second Genesis) {
	if first.GazPricePerKb() != second.GazPricePerKb() {
		t.Errorf("the returned gaz price is invalid.  Expected: %d, Returned: %d", first.GazPricePerKb(), second.GazPricePerKb())
		return
	}

	if first.MaxAmountOfValidators() != second.MaxAmountOfValidators() {
		t.Errorf("the returned maximum amount of validatoirs is invalid.  Expected: %d, Returned: %d", first.MaxAmountOfValidators(), second.MaxAmountOfValidators())
		return
	}

	deposit.CompareDepositForTests(t, first.Deposit(), second.Deposit())
}
