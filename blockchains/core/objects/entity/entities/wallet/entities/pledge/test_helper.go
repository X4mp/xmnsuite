package pledge

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/crypto"
)

// CreatePledgeForTests creates a pledge instance for tests
func CreatePledgeForTests() Pledge {
	id := uuid.NewV4()
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	from := withdrawal.CreateWithdrawalWithPublicKeyForTests(pubKey)
	to := wallet.CreateWalletWithPublicKeyForTests(pubKey)
	out := createPledge(&id, from, to)
	return out
}

// CreatePledgeWithWalletForTests creates a pledge with wallet instance for tests
func CreatePledgeWithWalletForTests(from wallet.Wallet, to wallet.Wallet, tok token.Token, amount int) Pledge {
	id := uuid.NewV4()
	fromWith := withdrawal.CreateWithdrawalWithTokenAndWalletForTests(tok, from, amount)
	out := createPledge(&id, fromWith, to)
	return out
}

// CreatePledgeWithAmountForTests creates a pledge instance with amount for tests
func CreatePledgeWithAmountForTests(amount int) Pledge {
	id := uuid.NewV4()
	pk := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := pk.PublicKey()
	from := withdrawal.CreateWithdrawalWithPublicKeyAndAmountForTests(pubKey, amount)
	to := wallet.CreateWalletWithPublicKeyForTests(pubKey)
	out := createPledge(&id, from, to)
	return out
}

// ComparePledgesForTests compares 2 Pledge instances for tests
func ComparePledgesForTests(t *testing.T, first Pledge, second Pledge) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid")
		return
	}

	withdrawal.CompareWithdrawalsForTests(t, first.From(), second.From())
	wallet.CompareWalletsForTests(t, first.To(), second.To())
}
