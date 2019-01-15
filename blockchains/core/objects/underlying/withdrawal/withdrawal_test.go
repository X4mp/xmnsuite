package withdrawal

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
	from := wallet.CreateWalletForTests()
	amount := rand.Int()

	// execute:
	with, withErr := createWithdrawal(&id, from, amount)
	if withErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", withErr.Error())
		return
	}

	// compare:
	if !reflect.DeepEqual(&id, with.ID()) {
		t.Errorf("the returned ID is invalid.  Expected: %s, Returned: %s", id.String(), with.ID().String())
		return
	}

	if !reflect.DeepEqual(from, with.From()) {
		t.Errorf("the returned From is invalid")
		return
	}

	if !reflect.DeepEqual(amount, with.Amount()) {
		t.Errorf("the returned Amount is invalid.  Expected: %d, Returned: %d", amount, with.Amount())
		return
	}

}

func TestCreate_withNegativeAmount_returnsError(t *testing.T) {
	// variables:
	id := uuid.NewV4()
	from := wallet.CreateWalletForTests()
	amount := rand.Int() * -1

	// execute:
	_, withErr := createWithdrawal(&id, from, amount)
	if withErr == nil {
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
	_, withErr := createWithdrawal(&id, from, amount)
	if withErr == nil {
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
	_, withErr := createWithdrawal(&id, from, amount)
	if withErr == nil {
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
	_, withErr := createWithdrawal(&id, from, amount)
	if withErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned")
		return
	}
}
