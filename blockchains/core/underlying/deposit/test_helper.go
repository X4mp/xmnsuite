package deposit

import (
	"math/rand"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

// CreateDepositForTests creates an Deposit for tests
func CreateDepositForTests() Deposit {
	id := uuid.NewV4()
	wal := wallet.CreateWalletForTests()
	tok := token.CreateTokenForTests()
	amount := rand.Int()
	out := createDeposit(&id, wal, tok, amount)
	return out
}

// CompareDepositForTests compares Deposit instances for tests
func CompareDepositForTests(t *testing.T, first Deposit, second Deposit) {
	if first.Amount() != second.Amount() {
		t.Errorf("the returned amount is invalid.  Expected: %d, Returned: %d", first.Amount(), second.Amount())
		return
	}

	wallet.CompareWalletsForTests(t, first.To(), second.To())
}
