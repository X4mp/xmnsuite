package xmn

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func createPledgeForTests() Pledge {
	id := uuid.NewV4()
	from := createWalletForTests()
	to := createWalletForTests()
	amount := rand.Int()
	out := createPledge(&id, from, to, amount)
	return out
}

func comparePledgeForTests(t *testing.T, first Pledge, second Pledge) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the returned ID is invalid.  Expected: %s, Returned: %s", first.ID().String(), second.ID().String())
		return
	}

	if first.Amount() != second.Amount() {
		t.Errorf("the returned amount is invalid.  Expected: %d, Returned: %d", first.Amount(), second.Amount())
		return
	}

	compareWalletsForTests(t, first.From(), second.From())
	compareWalletsForTests(t, first.To(), second.To())
}

func TestPledge_Success(t *testing.T) {
	first := createPledgeForTests()

	// create services:
	store := datastore.SDKFunc.Create()
	walletService := createWalletService(store)
	pledgeService := createPledgeService(store, walletService)

	// save the from wallet:
	saveFromWalletErr := walletService.Save(first.From())
	if saveFromWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveFromWalletErr.Error())
		return
	}

	// save the to wallet:
	saveToWalletErr := walletService.Save(first.To())
	if saveToWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveToWalletErr.Error())
		return
	}

	// save the Pledge:
	saveFirstErr := pledgeService.Save(first)
	if saveFirstErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveFirstErr.Error())
		return
	}

	// retrieve the pledge by ID:
	retFirstPledge, retFirstPledgeErr := pledgeService.RetrieveByID(first.ID())
	if retFirstPledgeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retFirstPledgeErr.Error())
		return
	}

	// retrieve the pledge by fromWalletID:
	retFirstPledgeByFromWalletID, retFirstPledgeByFromWalletIDErr := pledgeService.RetrieveByFromWalletID(first.From().ID(), 0, -1)
	if retFirstPledgeByFromWalletIDErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retFirstPledgeByFromWalletIDErr.Error())
		return
	}

	if retFirstPledgeByFromWalletID.TotalAmount() != 1 {
		t.Errorf("the total amount was expected to be 1, %d returned", retFirstPledgeByFromWalletID.TotalAmount())
		return
	}

	// retrieve the pledge by toWalletID:
	retFirstPledgeByToWalletID, retFirstPledgeByToWalletIDErr := pledgeService.RetrieveByToWalletID(first.To().ID(), 0, -1)
	if retFirstPledgeByToWalletIDErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retFirstPledgeByToWalletIDErr.Error())
		return
	}

	if retFirstPledgeByToWalletID.TotalAmount() != 1 {
		t.Errorf("the total amount was expected to be 1, %d returned", retFirstPledgeByToWalletID.TotalAmount())
		return
	}

	// compare:
	comparePledgeForTests(t, first, retFirstPledge)

	pledgesPledgeByFromWalletID := retFirstPledgeByFromWalletID.Pledges()
	comparePledgeForTests(t, first, pledgesPledgeByFromWalletID[0])

	pledgesByToWalletID := retFirstPledgeByToWalletID.Pledges()
	comparePledgeForTests(t, first, pledgesByToWalletID[0])

	// convert:
	empty := new(pledge)
	tests.ConvertToBinary(t, first, empty, cdc)

	anotherEmpty := new(pledge)
	tests.ConvertToJSON(t, first, anotherEmpty, cdc)
}
