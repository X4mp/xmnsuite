package user

import (
	"testing"

	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func TestUser_Success(t *testing.T) {
	// variables:
	usr := CreateUserForTests()
	anotherUsr := CreateUserWithWalletForTests(usr.Wallet())

	// create repository:
	store := datastore.SDKFunc.Create()
	repository := entity.SDKFunc.CreateRepository(entity.CreateRepositoryParams{
		Store: store,
	})

	// create the service:
	service := entity.SDKFunc.CreateService(entity.CreateServiceParams{
		Store: store,
	})

	// create the metadata and representation:
	metadata := SDKFunc.CreateMetaData()
	represenation := SDKFunc.CreateRepresentation()

	// save the wallet:
	walletRepresentation := wallet.SDKFunc.CreateRepresentation()
	saveWalletErr := service.Save(usr.Wallet(), walletRepresentation)
	if saveWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveWalletErr.Error())
		return
	}

	// save the user:
	savedErr := service.Save(usr, represenation)
	if savedErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", savedErr.Error())
		return
	}

	// save again, returns error:
	saveAgainErr := service.Save(usr, represenation)
	if saveAgainErr == nil {
		t.Errorf("the returned error was expected to be valid, nil returned.")
		return
	}

	// save another instance:
	saveAnotherErr := service.Save(anotherUsr, represenation)
	if saveAnotherErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveAnotherErr.Error())
		return
	}

	// retrieve by id:
	retUsr, retUsrErr := repository.RetrieveByID(metadata, usr.ID())
	if retUsrErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retUsrErr.Error())
		return
	}

	// compare:
	CompareUserForTests(t, usr.(User), retUsr.(User))

	// retrieve, should have 2 usrlets:
	retUsrs, retUsrsErr := repository.RetrieveSetByKeyname(metadata, retrieveAllUserKeyname(), 0, 20)
	if retUsrsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retUsrsErr.Error())
		return
	}

	if retUsrs.Amount() != 2 {
		t.Errorf("the was supposed to be %d usrlets, %d returned", 2, retUsrs.Amount())
		return
	}

	empty := new(user)
	tests.ConvertToBinary(t, usr, empty, cdc)

	anotherEmpty := new(user)
	tests.ConvertToJSON(t, usr, anotherEmpty, cdc)
}
