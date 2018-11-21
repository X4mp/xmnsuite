package withdrawal

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	"github.com/xmnservices/xmnsuite/crypto"
)

// CreateWithdrawalWithPublicKeyForTests creates a withdrawal instance for tests
func CreateWithdrawalWithPublicKeyForTests(pubKey crypto.PublicKey) Withdrawal {
	id := uuid.NewV4()
	fromWallet := wallet.CreateWalletWithPublicKeyForTests(pubKey)
	tok := token.CreateTokenForTests()
	amount := rand.Int() % 200
	out, _ := createWithdrawal(&id, fromWallet, tok, amount)
	return out
}

// CompareWithdrawalsForTests compares 2 Withdrawals instances for tests
func CompareWithdrawalsForTests(t *testing.T, first Withdrawal, second Withdrawal) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid")
		return
	}

	if !reflect.DeepEqual(first.Amount(), second.Amount()) {
		t.Errorf("the amount is invalid.  Expected: %d, Returned: %d", first.Amount(), second.Amount())
		return
	}

	wallet.CompareWalletsForTests(t, first.From(), second.From())
	token.CompareTokensForTests(t, first.Token(), second.Token())
}
