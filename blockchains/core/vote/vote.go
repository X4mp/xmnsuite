package vote

import (
	"bytes"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
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
