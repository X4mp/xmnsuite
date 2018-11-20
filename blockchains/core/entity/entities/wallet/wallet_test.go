package wallet

import (
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func TestWallet_Success(t *testing.T) {
	// variables:
	privKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	pubKey := privKey.PublicKey()

	wal := CreateWalletWithPublicKeyForTests(pubKey)
	anotherWal := CreateWalletWithPublicKeyForTests(pubKey)

	// create repository and service:
	store := datastore.SDKFunc.Create()
	repository := entity.SDKFunc.CreateRepository(store)
	service := entity.SDKFunc.CreateService(store)

	// create the metadata and representation:
	metadata := SDKFunc.CreateMetaData()
	represenation := SDKFunc.CreateRepresentation()

	// save the wallet:
	savedErr := service.Save(wal, represenation)
	if savedErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", savedErr.Error())
		return
	}

	// save again, returns error:
	saveAgainErr := service.Save(wal, represenation)
	if saveAgainErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned.")
		return
	}

	// save another instance:
	saveAnotherErr := service.Save(anotherWal, represenation)
	if saveAnotherErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveAnotherErr.Error())
		return
	}

	// retrieve by id:
	retWal, retWalErr := repository.RetrieveByID(metadata, wal.ID())
	if retWalErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retWalErr.Error())
		return
	}

	// compare:
	CompareWalletsForTests(t, wal.(Wallet), retWal.(Wallet))

	// retrieve, should have 2 wallets:
	retWals, retWalsErr := repository.RetrieveSetByKeyname(metadata, retrieveAllWalletKeyname(), 0, 20)
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
