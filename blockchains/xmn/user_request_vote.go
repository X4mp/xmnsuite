package xmn

import (
	"errors"
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/datastore/objects"
)

type storedUserRequestVote struct {
	ID      string `json:"id"`
	ReqID   string `json:"user_request_id"`
	VoterID string `json:"voter_id"`
	IsAppr  bool   `json:"is_approved"`
}

func createStoredUserRequestVote(vote UserRequestVote) *storedUserRequestVote {
	out := storedUserRequestVote{
		ID:      vote.ID().String(),
		ReqID:   vote.Request().User().ID().String(),
		VoterID: vote.Voter().ID().String(),
		IsAppr:  vote.IsApproved(),
	}

	return &out
}

type userRequestVote struct {
	UUID   *uuid.UUID  `json:"id"`
	Req    UserRequest `json:"user_request"`
	Vot    User        `json:"voter"`
	IsAppr bool        `json:"is_approved"`
}

func createUserRequestVote(id *uuid.UUID, req UserRequest, voter User, isAppr bool) UserRequestVote {
	out := userRequestVote{
		UUID:   id,
		Req:    req,
		Vot:    voter,
		IsAppr: isAppr,
	}

	return &out
}

// ID returns the ID
func (obj *userRequestVote) ID() *uuid.UUID {
	return obj.UUID
}

// Request returns the request
func (obj *userRequestVote) Request() UserRequest {
	return obj.Req
}

// Voter returns the voter
func (obj *userRequestVote) Voter() User {
	return obj.Vot
}

// IsApproved returns true if the request is approved, false otherwise
func (obj *userRequestVote) IsApproved() bool {
	return obj.IsAppr
}

type userRequestVotePartialSet struct {
	Votes  []UserRequestVote `json:"users"`
	Indx   int               `json:"index"`
	TotAmt int               `json:"total_amount"`
}

func createUserRequestVotePartialSet(votes []UserRequestVote, indx int, totAmt int) UserRequestVotePartialSet {
	out := userRequestVotePartialSet{
		Votes:  votes,
		Indx:   indx,
		TotAmt: totAmt,
	}

	return &out
}

// UserRequestVotes returns the []UserRequestVote
func (obj *userRequestVotePartialSet) UserRequestVotes() []UserRequestVote {
	return obj.Votes
}

// Index returns the index
func (obj *userRequestVotePartialSet) Index() int {
	return obj.Indx
}

// Amount returns the amount
func (obj *userRequestVotePartialSet) Amount() int {
	return len(obj.Votes)
}

// TotalAmount returns the total amount
func (obj *userRequestVotePartialSet) TotalAmount() int {
	return obj.TotAmt
}

type userRequestVoteService struct {
	keyname         string
	concensusNeeded int
	store           datastore.DataStore
	userReqService  UserRequestService
	userService     UserService
}

func createUserRequestVoteService(concensusNeeded int, store datastore.DataStore, userReqService UserRequestService, userService UserService) UserRequestVoteService {
	out := userRequestVoteService{
		keyname:         "user_request_vote",
		concensusNeeded: concensusNeeded,
		store:           store,
		userReqService:  userReqService,
		userService:     userService,
	}

	return &out
}

// Save saves a UserRequestVote instance
func (app *userRequestVoteService) Save(vote UserRequestVote) error {

	// fetch some data:
	voteID := vote.ID()
	req := vote.Request()
	requesterUsr := req.User()
	voter := vote.Voter()

	// make sure the vote does not already exists:
	_, retVoteErr := app.RetrieveByID(vote.ID())
	if retVoteErr == nil {
		str := fmt.Sprintf("the UserRequestVote (ID: %s) already exists", vote.ID().String())
		return errors.New(str)
	}

	// make sure the voter never voted on this request before:
	_, retVoteByVoterAndUserReqIDErr := app.RetrieveByVoterIDAndUserRequestID(voter.ID(), requesterUsr.ID())
	if retVoteByVoterAndUserReqIDErr == nil {
		str := fmt.Sprintf("the User (ID: %s) already voted on this UserRequest (ID: %s)", voter.ID().String(), requesterUsr.ID().String())
		return errors.New(str)
	}

	// create the set keys:
	keys := []string{
		app.keynameByRequesterWalletID(requesterUsr.Wallet().ID()),
		app.keynameByUserRequestID(requesterUsr.ID()),
		app.keynameByVoterID(voter.ID()),
		app.keynameByVoterWalletID(voter.Wallet().ID()),
	}

	// add the ID to the set keynames:
	amountAddedToSets := app.store.Sets().AddMul(keys, voteID.String())

	if amountAddedToSets != 1 {
		// revert:
		app.store.Sets().DelMul(keys, voteID.String())

		// returns error:
		str := fmt.Sprintf("there was an error while adding the UserRequestVote (ID: %s) to the sets... reverting", voteID.String())
		return errors.New(str)
	}

	// save the object:
	keyname := app.keynameByID(voteID)
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: keyname,
		Obj: createStoredUserRequestVote(vote),
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the UserRequestVote instance")
	}

	// retrieve the votes:
	votes, votesErr := app.RetrieveByUserRequestID(req.User().ID(), 0, -1)
	if votesErr != nil {
		return votesErr
	}

	// if there is now a concensus:
	approvedVotes, disApprovedVotes, weightErr := app.retrieveVoteWeight(votes)
	if weightErr != nil {
		return weightErr
	}

	weight := approvedVotes + disApprovedVotes
	if app.concensusNeeded <= weight {
		// delete the votes:
		rqVotes := votes.UserRequestVotes()
		for _, oneVote := range rqVotes {
			delErr := app.Delete(oneVote)
			if delErr != nil {
				log.Printf("there was an error while deleting the UserRequestVote (ID: %s): %s", oneVote.ID().String(), delErr.Error())
			}
		}

		// delete the request:
		delReqErr := app.userReqService.Delete(req)
		if delReqErr != nil {
			log.Printf("there was an error while deleting the UserRequest (ID: %s): %s", req.User().ID().String(), delReqErr.Error())
		}

		// if the vote passed, save the user:
		if approvedVotes >= disApprovedVotes {
			usr := req.User()
			saveUsrErr := app.userService.Save(usr)
			if saveUsrErr != nil {
				str := fmt.Sprintf("there was an error while saving the User (ID: %s) even if it reaches concensus (obtained: %d, needed: %d, approved: %d, disapproved: %d): %s", usr.ID().String(), weight, app.concensusNeeded, approvedVotes, disApprovedVotes, saveUsrErr.Error())
				return errors.New(str)
			}
		}
	}

	return nil
}

func (app *userRequestVoteService) retrieveVoteWeight(votes UserRequestVotePartialSet) (int, int, error) {
	approved := 0
	disApproved := 0
	for _, oneVote := range votes.UserRequestVotes() {
		if oneVote.IsApproved() {
			approved += oneVote.Voter().Shares()
			continue
		}

		disApproved += oneVote.Voter().Shares()
	}

	return approved, disApproved, nil
}

// Delete deletes a UserRequestVote instance
func (app *userRequestVoteService) Delete(vote UserRequestVote) error {
	// fetch some data:
	voteID := vote.ID()
	req := vote.Request()
	requesterUsr := req.User()
	voter := vote.Voter()

	// make sure the vote exists:
	_, retVoteErr := app.RetrieveByID(vote.ID())
	if retVoteErr != nil {
		str := fmt.Sprintf("the UserRequestVote (ID: %s) does not exists: %s", vote.ID().String(), retVoteErr.Error())
		return errors.New(str)
	}

	// create the set keys:
	keys := []string{
		app.keynameByRequesterWalletID(requesterUsr.Wallet().ID()),
		app.keynameByUserRequestID(requesterUsr.ID()),
		app.keynameByVoterID(voter.ID()),
		app.keynameByVoterWalletID(voter.Wallet().ID()),
	}

	// delete the ID from the sets:
	amountDeletedFromSets := app.store.Sets().DelMul(keys, voteID.String())

	if amountDeletedFromSets != 1 {
		// revert:
		app.store.Sets().DelMul(keys, voteID.String())

		// returns error:
		str := fmt.Sprintf("there was an error while deleting the UserRequestVote (ID: %s) from the sets... reverting", voteID.String())
		return errors.New(str)
	}

	// save the object:
	keyname := app.keynameByID(voteID)
	amountSaved := app.store.Objects().Save(&objects.ObjInKey{
		Key: keyname,
		Obj: storedUserRequestVote{
			ID:      voteID.String(),
			ReqID:   req.User().ID().String(),
			VoterID: voter.ID().String(),
			IsAppr:  vote.IsApproved(),
		},
	})

	if amountSaved != 1 {
		return errors.New("there was an error while saving the UserRequestVote instance")
	}

	return nil
}

// RetrieveByID retrieves a UserRequestVote instance by ID
func (app *userRequestVoteService) RetrieveByID(id *uuid.UUID) (UserRequestVote, error) {
	keyname := app.keynameByID(id)
	obj := objects.ObjInKey{
		Key: keyname,
		Obj: new(storedUserRequestVote),
	}

	amount := app.store.Objects().Retrieve(&obj)
	if amount != 1 {
		str := fmt.Sprintf("there was an error while retrieving the UserRequestVote (ID: %s)", id.String())
		return nil, errors.New(str)
	}

	if vot, ok := obj.Obj.(*storedUserRequestVote); ok {
		return app.FromStoredToUserRequestVote(vot)
	}

	return nil, errors.New("the retrieved data cannot be casted to a UserRequestVote instance")
}

// RetrieveByVoterIDAndUserRequestID retrieves a UserRequestVote instance by its voterID and UserRequestID
func (app *userRequestVoteService) RetrieveByVoterIDAndUserRequestID(voterID *uuid.UUID, requestID *uuid.UUID) (UserRequestVote, error) {

	// create keys:
	keys := []string{
		app.keynameByVoterID(voterID),
		app.keynameByUserRequestID(requestID),
	}

	// intersect:
	ids := app.store.Sets().Inter(keys...)

	// if there is no values:
	if len(ids) <= 0 {
		str := fmt.Sprintf("there is no UserRequestVote that contains both that Voter (ID: %s) and that UserRequest (ID: %s)", voterID.String(), requestID.String())
		return nil, errors.New(str)
	}

	if len(ids) == 1 {
		// cast the ID:
		id, idErr := uuid.FromString(ids[0].(string))
		if idErr != nil {
			str := fmt.Sprintf("the element stored in the set is not a valid UUID: %s", idErr.Error())
			return nil, errors.New(str)
		}

		// retrieve the instance, then return it:
		return app.RetrieveByID(&id)
	}

	str := fmt.Sprintf("there is %d UserRequestVote instances that contains both that Voter (ID: %s) and that UserRequest (ID: %s), this should never happen", len(ids), voterID.String(), requestID.String())
	return nil, errors.New(str)
}

// RetrieveByWalletID retrieves []UserRequestVote instances by walletID
func (app *userRequestVoteService) RetrieveByRequesterWalletID(walletID *uuid.UUID, index int, amount int) (UserRequestVotePartialSet, error) {
	keyname := app.keynameByRequesterWalletID(walletID)
	return app.retrieveByKeyname(keyname, index, amount)
}

// RetrieveByUserRequestID retrieves []UserRequestVote instances by userRequestID
func (app *userRequestVoteService) RetrieveByUserRequestID(requestID *uuid.UUID, index int, amount int) (UserRequestVotePartialSet, error) {
	keyname := app.keynameByUserRequestID(requestID)
	return app.retrieveByKeyname(keyname, index, amount)
}

// RetrieveByVoterID retrieves []UserRequestVote instances by voterID
func (app *userRequestVoteService) RetrieveByVoterID(voterID *uuid.UUID, index int, amount int) (UserRequestVotePartialSet, error) {
	keyname := app.keynameByVoterID(voterID)
	return app.retrieveByKeyname(keyname, index, amount)
}

// RetrieveByVoterWalletID retrieves []UserRequestVote instances by voter walletID
func (app *userRequestVoteService) RetrieveByVoterWalletID(walletID *uuid.UUID, index int, amount int) (UserRequestVotePartialSet, error) {
	keyname := app.keynameByVoterWalletID(walletID)
	return app.retrieveByKeyname(keyname, index, amount)
}

// FromStoredToUserRequestVote converts a StoredUserRequestVote to a UserRequestVote instance
func (app *userRequestVoteService) FromStoredToUserRequestVote(vote *storedUserRequestVote) (UserRequestVote, error) {

	// cast the requestID:
	reqID, reqIDErr := uuid.FromString(vote.ReqID)
	if reqIDErr != nil {
		return nil, reqIDErr
	}

	// retrieve the request using its ID:
	retRequest, retRequestErr := app.userReqService.RetrieveByID(&reqID)
	if retRequestErr != nil {
		return nil, retRequestErr
	}

	// cast the voterID:
	voterID, voterIDErr := uuid.FromString(vote.VoterID)
	if voterIDErr != nil {
		return nil, voterIDErr
	}

	// retrieve the user voter by its ID:
	retVoter, retVoterErr := app.userService.RetrieveByID(&voterID)
	if retVoterErr != nil {
		return nil, retVoterErr
	}

	// cast the voteID:
	voteID, voteIDErr := uuid.FromString(vote.ID)
	if voteIDErr != nil {
		return nil, voteIDErr
	}

	return createUserRequestVote(&voteID, retRequest, retVoter, vote.IsAppr), nil
}

func (app *userRequestVoteService) retrieveByKeyname(keyname string, index int, amount int) (UserRequestVotePartialSet, error) {
	votes := []UserRequestVote{}
	uncastedVoteIDs := app.store.Sets().Retrieve(keyname, index, amount)
	for _, oneUncastedVoteID := range uncastedVoteIDs {
		voteID, voteIDErr := uuid.FromString(oneUncastedVoteID.(string))
		if voteIDErr != nil {
			str := fmt.Sprintf("one of the elements in the set (key: %s) is not a valid ID (element: %s): %s", keyname, oneUncastedVoteID.(string), voteIDErr.Error())
			return nil, errors.New(str)
		}

		oneVote, oneVoteErr := app.RetrieveByID(&voteID)
		if oneVoteErr != nil {
			return nil, oneVoteErr
		}

		votes = append(votes, oneVote)
	}

	// retireve the total amount:
	totAmount := app.store.Sets().Len(keyname)
	ps := createUserRequestVotePartialSet(votes, index, totAmount)
	return ps, nil
}

func (app *userRequestVoteService) keynameByID(id *uuid.UUID) string {
	return fmt.Sprintf("%s:by_id:%s", app.keyname, id.String())
}

func (app *userRequestVoteService) keynameByRequesterWalletID(walletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_requester_wallet_id:%s", app.keyname, walletID.String())
}

func (app *userRequestVoteService) keynameByUserRequestID(requestID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_user_request_id:%s", app.keyname, requestID.String())
}

func (app *userRequestVoteService) keynameByVoterID(voterID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_voter_id:%s", app.keyname, voterID.String())
}

func (app *userRequestVoteService) keynameByVoterWalletID(walletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_voter_wallet_id:%s", app.keyname, walletID.String())
}
