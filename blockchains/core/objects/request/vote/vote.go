package vote

import (
	"bytes"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
)

type vote struct {
	UUID   *uuid.UUID      `json:"id"`
	Req    request.Request `json:"request"`
	Votr   user.User       `json:"voter"`
	IsAppr bool            `json:"is_approved"`
}

func createVote(id *uuid.UUID, req request.Request, votr user.User, isApproved bool) (Vote, error) {
	// make sure the voter has the same walletID as the requester:
	requesterWalletID := req.From().Wallet().ID()
	voterWalletID := votr.Wallet().ID()
	if bytes.Compare(requesterWalletID.Bytes(), voterWalletID.Bytes()) != 0 {
		str := fmt.Sprintf("the requester is binded to a wallet (ID: %s) that is different from the voter's wallet (ID: %s)", requesterWalletID.String(), voterWalletID.String())
		return nil, errors.New(str)
	}

	out := vote{
		UUID:   id,
		Req:    req,
		Votr:   votr,
		IsAppr: isApproved,
	}

	return &out, nil
}

func createVoteFromNormalized(normalized *normalizedVote) (Vote, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	reqIns, reqInsErr := request.SDKFunc.CreateMetaData().Denormalize()(normalized.Request)
	if reqInsErr != nil {
		return nil, reqInsErr
	}

	voterIns, voterInsErr := user.SDKFunc.CreateMetaData().Denormalize()(normalized.Voter)
	if voterInsErr != nil {
		return nil, voterInsErr
	}

	if req, ok := reqIns.(request.Request); ok {
		if voter, ok := voterIns.(user.User); ok {
			return createVote(&id, req, voter, normalized.IsApproved)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", voterIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Request instance", reqIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *vote) ID() *uuid.UUID {
	return obj.UUID
}

// Request returns the request
func (obj *vote) Request() request.Request {
	return obj.Req
}

// Voter returns the voter
func (obj *vote) Voter() user.User {
	return obj.Votr
}

// IsApproved returns true if the vote is approved, false otherwise
func (obj *vote) IsApproved() bool {
	return obj.IsAppr
}
