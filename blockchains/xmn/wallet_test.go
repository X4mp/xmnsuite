package xmn

import (
	"math/rand"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func createWalletForTests() Wallet {
	id := uuid.NewV4()
	concensusNeeded := rand.Float64()

	out := createWallet(&id, concensusNeeded)
	return out
}

func compareWalletsForTests(t *testing.T, first Wallet, second Wallet) {
	if !reflect.DeepEqual(first.ID(), second.ID()) {
		t.Errorf("the ID is invalid")
		return
	}

	if !reflect.DeepEqual(first.ConcensusNeeded(), second.ConcensusNeeded()) {
		t.Errorf("the concensusNeeded is invalid.  Expected: %f, Returned: %f", first.ConcensusNeeded(), second.ConcensusNeeded())
		return
	}
}

func TestWallet_Success(t *testing.T) {
	// variables:
	wal := createWalletForTests()
	anotherWal := createWalletForTests()

	// create services:
	store := datastore.SDKFunc.Create()
	walletService := createWalletService(store)

	// save the wallet:
	savedErr := walletService.Save(wal)
	if savedErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", savedErr.Error())
		return
	}

	// save again, returns error:
	saveAgainErr := walletService.Save(wal)
	if saveAgainErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned.")
		return
	}

	// save another instance:
	saveAnotherErr := walletService.Save(anotherWal)
	if saveAnotherErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveAnotherErr.Error())
		return
	}

	// retrieve by id:
	retWal, retWalErr := walletService.RetrieveByID(wal.ID())
	if retWalErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retWalErr.Error())
		return
	}

	// compare:
	compareWalletsForTests(t, wal, retWal)

	// retrieve, should have 2 wallets:
	retWals, retWalsErr := walletService.Retrieve(0, 20)
	if retWalsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retWalsErr.Error())
		return
	}

	if retWals.Amount() != 2 {
		t.Errorf("the was supposed to be %d wallets, %d returned", 2, retWals.Amount())
		return
	}

	empty := new(wallet)
	tests.ConvertToBinary(t, wal, empty, cdc)

	anotherEmpty := new(wallet)
	tests.ConvertToJSON(t, wal, anotherEmpty, cdc)
}
