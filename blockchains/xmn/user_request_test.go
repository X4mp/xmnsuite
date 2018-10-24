package xmn

import (
	"bytes"
	"testing"

	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/tests"
)

func createUserRequestForTests() UserRequest {
	usr := createUserForTests()
	req := createUserRequest(usr)
	return req
}

func createUserRequestWithWalletForTests(wal Wallet) UserRequest {
	usr := createUserWithWalletForTests(wal)
	req := createUserRequest(usr)
	return req
}

func createUserRequestWithPubKeyForTests(pubKey crypto.PublicKey) UserRequest {
	usr := createUserWithPublicKeyForTests(pubKey)
	req := createUserRequest(usr)
	return req
}

func createUserRequestWithWalletAndPubKeyForTests(wal Wallet, pubKey crypto.PublicKey) UserRequest {
	usr := createUserWithWalletAndPublicKeyForTests(wal, pubKey)
	req := createUserRequest(usr)
	return req
}

func TestUserRequest_Success(t *testing.T) {
	firstReq := createUserRequestForTests()
	secondReq := createUserRequestForTests()

	// create services:
	store := datastore.SDKFunc.Create()
	walletService := createWalletService(store)
	userRequestService := createUserRequestService(store, walletService)

	// try to delete a UserRequest not saved yet, returns error:
	firstDelBeforeSavingErr := userRequestService.Delete(firstReq)
	if firstDelBeforeSavingErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned: %s", firstDelBeforeSavingErr.Error())
		return
	}

	// try to save the first user request before saving the UserRequest wallet, returns an error:
	firstSaveBeforeWalletErr := userRequestService.Save(firstReq)
	if firstSaveBeforeWalletErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	// save the wallet:
	firstSaveWalletErr := walletService.Save(firstReq.User().Wallet())
	if firstSaveWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveWalletErr)
		return
	}

	// save the first user request:
	firstSaveErr := userRequestService.Save(firstReq)
	if firstSaveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveErr)
		return
	}

	// save again, should return an error:
	firstSaveAgainErr := userRequestService.Save(firstReq)
	if firstSaveAgainErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	// retrieve by ID:
	firstRet, firstRetErr := userRequestService.RetrieveByID(firstReq.User().ID())
	if firstRetErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstRetErr)
		return
	}

	// compare:
	compareUserForTests(t, firstReq.User(), firstRet.User())

	// save the second wallet:
	secondSaveWalletErr := walletService.Save(secondReq.User().Wallet())
	if secondSaveWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondSaveWalletErr)
		return
	}

	// save the second UserRequest instance:
	secondSaveErr := userRequestService.Save(secondReq)
	if secondSaveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondSaveErr)
		return
	}

	// retrieve by walletID:
	retByWalletIDPartialSet, retByWalletIDPartialSetErr := userRequestService.RetrieveByWalletID(firstReq.User().Wallet().ID(), 0, -1)
	if retByWalletIDPartialSetErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retByWalletIDPartialSetErr)
		return
	}

	if retByWalletIDPartialSet.TotalAmount() != 1 {
		t.Errorf("the total amount was expected to be 1, %d returned", retByWalletIDPartialSet.TotalAmount())
		return
	}

	// compare:
	reqsByWalletIDs := retByWalletIDPartialSet.Requests()
	compareUserForTests(t, firstReq.User(), reqsByWalletIDs[0].User())

	// retrieve by pubkey:
	retByPubKeyPartialSet, retByPubKeyPartialSetErr := userRequestService.RetrieveByPubKey(firstReq.User().PubKey(), 0, -1)
	if retByPubKeyPartialSetErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retByPubKeyPartialSetErr)
		return
	}

	if retByPubKeyPartialSet.TotalAmount() != 1 {
		t.Errorf("the total amount was expected to be 1, %d returned", retByPubKeyPartialSet.TotalAmount())
		return
	}

	// compare:
	reqsByPubKeys := retByPubKeyPartialSet.Requests()
	compareUserForTests(t, firstReq.User(), reqsByPubKeys[0].User())

	// first delete:
	firstDelErr := userRequestService.Delete(firstReq)
	if firstDelErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstDelErr.Error())
		return
	}

	// second delete:
	secondDelErr := userRequestService.Delete(secondReq)
	if secondDelErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondDelErr.Error())
		return
	}

	// convert back and forth:
	empty := new(userRequest)
	tests.ConvertToBinary(t, firstReq, empty, cdc)

	anotherEmpty := new(userRequest)
	tests.ConvertToJSON(t, firstReq, anotherEmpty, cdc)
}

func TestUserRequest_onSameWalletID_Success(t *testing.T) {
	firstReq := createUserRequestForTests()
	secondReq := createUserRequestWithWalletForTests(firstReq.User().Wallet())

	// create services:
	store := datastore.SDKFunc.Create()
	walletService := createWalletService(store)
	userRequestService := createUserRequestService(store, walletService)

	// save the wallet:
	firstSaveWalletErr := walletService.Save(firstReq.User().Wallet())
	if firstSaveWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveWalletErr)
		return
	}

	// save the first user request:
	firstSaveErr := userRequestService.Save(firstReq)
	if firstSaveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveErr)
		return
	}

	// save the second user request:
	secondeSaveErr := userRequestService.Save(secondReq)
	if secondeSaveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondeSaveErr)
		return
	}

	// retrieve by walletID:
	retByWalletIDPartialSet, retByWalletIDPartialSetErr := userRequestService.RetrieveByWalletID(firstReq.User().Wallet().ID(), 0, -1)
	if retByWalletIDPartialSetErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retByWalletIDPartialSetErr)
		return
	}

	if retByWalletIDPartialSet.TotalAmount() != 2 {
		t.Errorf("the total amount was expected to be 2, %d returned", retByWalletIDPartialSet.TotalAmount())
		return
	}

	// compare:
	reqsByWalletIDs := retByWalletIDPartialSet.Requests()
	for _, oneReq := range reqsByWalletIDs {
		if bytes.Compare(oneReq.User().ID().Bytes(), firstReq.User().ID().Bytes()) == 0 {
			compareUserForTests(t, oneReq.User(), firstReq.User())
			return
		}

		if bytes.Compare(oneReq.User().ID().Bytes(), secondReq.User().ID().Bytes()) == 0 {
			compareUserForTests(t, oneReq.User(), secondReq.User())
			return
		}

		t.Errorf("the given UserRequest (ID: %s) is invalid", oneReq.User().ID().String())
	}
}

func TestUserRequest_onSamePubKey_Success(t *testing.T) {
	firstReq := createUserRequestForTests()
	secondReq := createUserRequestWithPubKeyForTests(firstReq.User().PubKey())

	// create services:
	store := datastore.SDKFunc.Create()
	walletService := createWalletService(store)
	userRequestService := createUserRequestService(store, walletService)

	// save the first wallet:
	firstSaveWalletErr := walletService.Save(firstReq.User().Wallet())
	if firstSaveWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveWalletErr)
		return
	}

	// save the first user request:
	firstSaveErr := userRequestService.Save(firstReq)
	if firstSaveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveErr)
		return
	}

	// save the second wallet:
	secondSaveWalletErr := walletService.Save(secondReq.User().Wallet())
	if secondSaveWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondSaveWalletErr)
		return
	}

	// save the second user request:
	secondeSaveErr := userRequestService.Save(secondReq)
	if secondeSaveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", secondeSaveErr)
		return
	}

	// retrieve by pubKey:
	retByWalletIDPartialSet, retByWalletIDPartialSetErr := userRequestService.RetrieveByPubKey(firstReq.User().PubKey(), 0, -1)
	if retByWalletIDPartialSetErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retByWalletIDPartialSetErr)
		return
	}

	if retByWalletIDPartialSet.TotalAmount() != 2 {
		t.Errorf("the total amount was expected to be 2, %d returned", retByWalletIDPartialSet.TotalAmount())
		return
	}

	// compare:
	reqsByPubKey := retByWalletIDPartialSet.Requests()
	for _, oneReq := range reqsByPubKey {
		if bytes.Compare(oneReq.User().ID().Bytes(), firstReq.User().ID().Bytes()) == 0 {
			compareUserForTests(t, oneReq.User(), firstReq.User())
			return
		}

		if bytes.Compare(oneReq.User().ID().Bytes(), secondReq.User().ID().Bytes()) == 0 {
			compareUserForTests(t, oneReq.User(), secondReq.User())
			return
		}

		t.Errorf("the given UserRequest (ID: %s) is invalid", oneReq.User().ID().String())
	}
}

func TestUserRequest_oneSameWalletID_onSamePubKey_ReturnsError(t *testing.T) {
	firstReq := createUserRequestForTests()
	secondReq := createUserRequestWithWalletAndPubKeyForTests(firstReq.User().Wallet(), firstReq.User().PubKey())

	// create services:
	store := datastore.SDKFunc.Create()
	walletService := createWalletService(store)
	userRequestService := createUserRequestService(store, walletService)

	// save the first wallet:
	firstSaveWalletErr := walletService.Save(firstReq.User().Wallet())
	if firstSaveWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveWalletErr)
		return
	}

	// save the first user request:
	firstSaveErr := userRequestService.Save(firstReq)
	if firstSaveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSaveErr)
		return
	}

	// save the second user request:
	secondSaveErr := userRequestService.Save(secondReq)
	if secondSaveErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}
}
