package pledge

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
)

// CreatePledgeForTests creates a pledge instance for tests
func CreatePledgeForTests() Pledge {
	id := uuid.NewV4()
	from := withdrawal.CreateWithdrawalForTests()
	to := wallet.CreateWalletForTests()
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
