package deposit

import (
	"math"
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

func TestCreate_Success(t *testing.T) {
	// variables:
	id := uuid.NewV4()
	to := wallet.CreateWalletForTests()
	amount := rand.Int()

	// execute:
	dep, depErr := createDeposit(&id, to, amount)
	if depErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", depErr.Error())
		return
	}

	// compare:
	if !reflect.DeepEqual(&id, dep.ID()) {
		t.Errorf("the returned ID is invalid.  Expected: %s, Returned: %s", id.String(), dep.ID().String())
		return
	}

	if !reflect.DeepEqual(to, dep.To()) {
		t.Errorf("the returned To is invalid")
		return
	}

	if !reflect.DeepEqual(amount, dep.Amount()) {
		t.Errorf("the returned Amount is invalid.  Expected: %d, Returned: %d", amount, dep.Amount())
		return
	}

}

func TestCreate_withNegativeAmount_returnsError(t *testing.T) {
	// variables:
	id := uuid.NewV4()
	from := wallet.CreateWalletForTests()
	amount := rand.Int() * -1

	// execute:
	_, depErr := createDeposit(&id, from, amount)
	if depErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return
	}
}

func TestCreate_withZeroAmount_returnsError(t *testing.T) {
	// variables:
	id := uuid.NewV4()
	from := wallet.CreateWalletForTests()
	amount := 0

	// execute:
	_, depErr := createDeposit(&id, from, amount)
	if depErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return
	}
}

func TestCreate_withOverflowAmount_returnsError(t *testing.T) {
	// variables:
	id := uuid.NewV4()
	from := wallet.CreateWalletForTests()
	amount := math.MaxInt64 + rand.Int()

	// execute:
	_, depErr := createDeposit(&id, from, amount)
	if depErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return
	}
}

func TestCreate_withTooBigAmount_returnsError(t *testing.T) {
	// variables:
	id := uuid.NewV4()
	from := wallet.CreateWalletForTests()
	amount := math.MaxInt64

	// execute:
	_, depErr := createDeposit(&id, from, amount)
	if depErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return
	}
}
