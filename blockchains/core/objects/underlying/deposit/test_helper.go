package deposit

import (
	"math/rand"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/crypto"
)

// CreateDepositWithPubKeyForTests creates an Deposit for tests
func CreateDepositWithPubKeyForTests(pubKey crypto.PublicKey) Deposit {
	id := uuid.NewV4()
	wal := wallet.CreateWalletWithPublicKeyForTests(pubKey)
	tok := token.CreateTokenForTests()
	amount := rand.Int() % 200
	out, _ := createDeposit(&id, wal, tok, amount)
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
