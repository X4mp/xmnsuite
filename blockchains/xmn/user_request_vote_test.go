package xmn

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
)

func createUserRequestVoteWithUserRequestForTests(req UserRequest, shares int, isApproved bool) UserRequestVote {
	id := uuid.NewV4()
	voter := createUserWithSharesForTests(shares)
	out := createUserRequestVote(&id, req, voter, isApproved)
	return out
}

func createUserRequestVoteWithUserRequestAndVoterForTests(req UserRequest, voter User, isApproved bool) UserRequestVote {
	id := uuid.NewV4()
	out := createUserRequestVote(&id, req, voter, isApproved)
	return out
}

func TestUserRequestVote_Success(t *testing.T) {
	req := createUserRequestForTests()
	firstVote := createUserRequestVoteWithUserRequestForTests(req, 9, true)
	secondVote := createUserRequestVoteWithUserRequestForTests(req, 8, false)
	thirdVote := createUserRequestVoteWithUserRequestForTests(req, 1, true)

	// create services:
	concensusNeeded := 10
	store := datastore.SDKFunc.Create()
	walletService := createWalletService(store)
	userReqService := createUserRequestService(store, walletService)
	userService := createUserService(store, walletService)
	userRequestVoteService := createUserRequestVoteService(concensusNeeded, store, userReqService, userService)

	// save the requester's wallet:
	saveReqWalletErr := walletService.Save(req.User().Wallet())
	if saveReqWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveReqWalletErr.Error())
		return
	}

	// save the first vote wallet:
	saveFirstVoteWalletErr := walletService.Save(firstVote.Voter().Wallet())
	if saveFirstVoteWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveFirstVoteWalletErr.Error())
		return
	}

	// save the second vote wallet:
	saveSecondVoteWalletErr := walletService.Save(secondVote.Voter().Wallet())
	if saveSecondVoteWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveSecondVoteWalletErr.Error())
		return
	}

	// save the third vote wallet:
	saveThirdVoteWalletErr := walletService.Save(thirdVote.Voter().Wallet())
	if saveThirdVoteWalletErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveThirdVoteWalletErr.Error())
		return
	}

	// save the user request:
	saveUserRequestErr := userReqService.Save(req)
	if saveUserRequestErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveUserRequestErr.Error())
		return
	}

	// save the first vote's user:
	saveFirstVoteUserErr := userService.Save(firstVote.Voter())
	if saveFirstVoteUserErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveFirstVoteUserErr.Error())
		return
	}

	// save the second vote's user:
	saveSecondVoteUserErr := userService.Save(secondVote.Voter())
	if saveSecondVoteUserErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveSecondVoteUserErr.Error())
		return
	}

	// save the third vote's user:
	saveThirdVoteUserErr := userService.Save(thirdVote.Voter())
	if saveThirdVoteUserErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveThirdVoteUserErr.Error())
		return
	}

	// save the first vote:
	saveFirstVoteErr := userRequestVoteService.Save(firstVote)
	if saveFirstVoteErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveFirstVoteErr.Error())
		return
	}

	// save the second vote:
	saveSecondVoteErr := userRequestVoteService.Save(secondVote)
	if saveSecondVoteErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveSecondVoteErr.Error())
		return
	}

	// save the third vote, we reached concensus, so it should fail:
	saveThirdVoteErr := userRequestVoteService.Save(thirdVote)
	if saveThirdVoteErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	// retrieve the user request, we reached concensus, so it should fail:
	_, retUserReqAfterConErr := userReqService.RetrieveByID(req.User().ID())
	if retUserReqAfterConErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	// retrieve the first vote, we reached concensus, so it should fail:
	_, retFirstVoteAfterConErr := userRequestVoteService.RetrieveByID(firstVote.ID())
	if retFirstVoteAfterConErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	// retrieve the second vote, we reached concensus, so it should fail:
	_, retSecondVoteAfterConErr := userRequestVoteService.RetrieveByID(secondVote.ID())
	if retSecondVoteAfterConErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	// retrieve the user:
	retUser, retUserErr := userService.RetrieveByID(req.User().ID())
	if retUserErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retUserErr.Error())
		return
	}

	// compare users:
	compareUserForTests(t, req.User(), retUser)
}
