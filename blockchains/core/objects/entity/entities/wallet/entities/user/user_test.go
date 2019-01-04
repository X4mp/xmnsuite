package user

import (
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func TestUser_Success(t *testing.T) {
	// variables:
	usr := CreateUserForTests()
	anotherUsr := CreateUserWithWalletForTests(usr.Wallet())

	// create the metadata and representation:
	metadata := SDKFunc.CreateMetaData()
	represenation := SDKFunc.CreateRepresentation()

	// create repository and services:
	store := datastore.SDKFunc.Create()
	entityRepository := entity.SDKFunc.CreateRepository(store)
	entityService := entity.SDKFunc.CreateService(store)
	repository := createRepository(metadata, entityRepository)

	// save the wallet:
	walletRepresentation := wallet.SDKFunc.CreateRepresentation()
	saveWalletErr := entityService.Save(usr.Wallet(), walletRepresentation)
	if saveWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveWalletErr.Error())
		return
	}

	// save the user:
	savedErr := entityService.Save(usr, represenation)
	if savedErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", savedErr.Error())
		return
	}

	// save again, returns error:
	saveAgainErr := entityService.Save(usr, represenation)
	if saveAgainErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned.")
		return
	}

	// save another instance:
	saveAnotherErr := entityService.Save(anotherUsr, represenation)
	if saveAnotherErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveAnotherErr.Error())
		return
	}

	// retrieve by id:
	retUsr, retUsrErr := entityRepository.RetrieveByID(metadata, usr.ID())
	if retUsrErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retUsrErr.Error())
		return
	}

	// compare:
	CompareUserForTests(t, usr.(User), retUsr.(User))

	// retrieve, should have 2 users:
	retUsrs, retUsrsErr := entityRepository.RetrieveSetByKeyname(metadata, retrieveAllUserKeyname(), 0, 20)
	if retUsrsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retUsrsErr.Error())
		return
	}

	if retUsrs.Amount() != 2 {
		t.Errorf("the was supposed to be %d usrlets, %d returned", 2, retUsrs.Amount())
		return
	}

	// retrieve, should have 1 user with that public key:
	retPS, retPSErr := repository.RetrieveSetByPubKey(usr.PubKey(), 0, 20)
	if retPSErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retPSErr.Error())
		return
	}

	if retPS.Amount() != 1 {
		t.Errorf("the was supposed to be %d usrlets, %d returned", 1, retPS.Amount())
		return
	}

	empty := new(user)
	tests.ConvertToBinary(t, usr, empty, cdc)

	anotherEmpty := new(user)
	tests.ConvertToJSON(t, usr, anotherEmpty, cdc)
}
