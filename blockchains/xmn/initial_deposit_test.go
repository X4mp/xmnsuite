package xmn

import (
	"math/rand"
	"testing"

	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func createInitialDepositForTests() InitialDeposit {
	wallet := createWalletForTests()
	amount := rand.Int()
	out := createInitialDeposit(wallet, amount)
	return out
}

func compareInitialDepositForTests(t *testing.T, first InitialDeposit, second InitialDeposit) {
	if first.Amount() != second.Amount() {
		t.Errorf("the returned amount is invalid.  Expected: %d, Returned: %d", first.Amount(), second.Amount())
		return
	}

	compareWalletsForTests(t, first.To(), second.To())
}

func TestInitialDeposit_Success(t *testing.T) {
	initialDep := createInitialDepositForTests()

	// create service:
	store := datastore.SDKFunc.Create()
	walletService := createWalletService(store)
	serv := createInitialDepositService(store, walletService)

	// save:
	saveErr := serv.Save(initialDep)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	// save again, returns error:
	saveAgainErr := serv.Save(createInitialDepositForTests())
	if saveAgainErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	// retrieve the token:
	retInitialDep, retInitialDepErr := serv.Retrieve()
	if retInitialDepErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retInitialDepErr.Error())
		return
	}

	// compare the elements:
	compareInitialDepositForTests(t, initialDep, retInitialDep)

	empty := new(initialDeposit)
	tests.ConvertToBinary(t, initialDep, empty, cdc)

	anotherEmpty := new(initialDeposit)
	tests.ConvertToJSON(t, initialDep, anotherEmpty, cdc)
}
