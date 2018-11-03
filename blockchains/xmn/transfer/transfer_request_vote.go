package xmn

import uuid "github.com/satori/go.uuid"

type transferRequestVote struct {
	UUID   *uuid.UUID            `json:"id"`
	Req    SignedTransferRequest `json:"signed_transfer_request"`
	Vot    User                  `json:"voter"`
	IsAppr bool                  `json:"is_approved"`
}

// ID returns the ID
func (obj *transferRequestVote) ID() *uuid.UUID {
	return obj.UUID
}

// Request returns the SignedTransferRequest
func (obj *transferRequestVote) Request() SignedTransferRequest {
	return obj.Req
}

// Voter returns the voter
func (obj *transferRequestVote) Voter() User {
	return obj.Vot
}

// IsApproved returns true if the voter approved the transfer, false otherwise
func (obj *transferRequestVote) IsApproved() bool {
	return obj.IsAppr
}

type transferRequestVotePartialSet struct {
	Vots  []TransferRequestVote `json:"votes"`
	Idx   int                   `json:"index"`
	TotAm int                   `json:"total_amount"`
}

func createTransferRequestVotePartialSet(vots []TransferRequestVote, index int, totalAmount int) TransferRequestVotePartialSet {
	out := transferRequestVotePartialSet{
		Vots:  vots,
		Idx:   index,
		TotAm: totalAmount,
	}

	return &out
}

// Votes returns the votes
func (obj *transferRequestVotePartialSet) Votes() []TransferRequestVote {
	return obj.Vots
}

// Index returns the index
func (obj *transferRequestVotePartialSet) Index() int {
	return obj.Idx
}

// Amount returns the amount
func (obj *transferRequestVotePartialSet) Amount() int {
	return len(obj.Vots)
}

// TotalAmount returns the totalAmount
func (obj *transferRequestVotePartialSet) TotalAmount() int {
	return obj.TotAm
}
