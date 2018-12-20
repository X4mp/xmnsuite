package pledge

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
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

// ComparePledgesForTests compares 2 Pledge instances for tests
func ComparePledgesForTests(t *testing.T, first Pledge, second Pledge) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid")
		return
	}

	withdrawal.CompareWithdrawalsForTests(t, first.From(), second.From())
	wallet.CompareWalletsForTests(t, first.To(), second.To())
}
