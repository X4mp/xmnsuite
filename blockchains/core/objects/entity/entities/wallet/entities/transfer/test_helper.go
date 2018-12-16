package transfer

import (
	"reflect"
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

// CompareTransfersForTests compares 2 Transfer instances for tests
func CompareTransfersForTests(t *testing.T, first Transfer, second Transfer) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid")
		return
	}

	withdrawal.CompareWithdrawalsForTests(t, first.Withdrawal(), second.Withdrawal())
	deposit.CompareDepositForTests(t, first.Deposit(), second.Deposit())
}
