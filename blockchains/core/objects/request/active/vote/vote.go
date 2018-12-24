package vote

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
)

type vote struct {
	UUID     *uuid.UUID      `json:"id"`
	Req      request.Request `json:"request"`
	Votr     user.User       `json:"voter"`
	Rson     string          `json:"reason"`
	IsNeutrl bool            `json:"is_neutral"`
	IsAppr   bool            `json:"is_approved"`
}

func createVote(id *uuid.UUID, req request.Request, votr user.User, reason string, isNeutrl bool, isApproved bool) (Vote, error) {
	if isNeutrl == isApproved {
		str := fmt.Sprintf("the vote cannot have the same value for neutral and approved")
		return nil, errors.New(str)
	}

	out := vote{
		UUID:     id,
		Req:      req,
		Votr:     votr,
		Rson:     reason,
		IsNeutrl: isNeutrl,
		IsAppr:   isApproved,
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
			return createVote(&id, req, voter, normalized.Reason, normalized.IsNeutral, normalized.IsApproved)
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
func (obj *vote) Reason() string {
	return obj.Rson
}

// IsNeutral returns true if the vote is neutral, false otherwise
func (obj *vote) IsNeutral() bool {
	return obj.IsNeutrl
}

// IsApproved returns true if the vote is approved, false otherwise
func (obj *vote) IsApproved() bool {
	return obj.IsAppr
}
