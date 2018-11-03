package xmn

import (
	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/crypto"
)

/*
 * TransferRequest
 */

// TransferRequest represents a transfer request of token
type TransferRequest interface {
	From() User
	Amount() string
	PubKey() crypto.PublicKey
	Reason() string
}

// SignedTransferRequest represents a signed transfer request
type SignedTransferRequest interface {
	ID() *uuid.UUID
	Request() TransferRequest
	Signature() crypto.RingSignature
}

// SignedTransferRequestPartialSet represents the signed transfer partial set
type SignedTransferRequestPartialSet interface {
	Requests() []SignedTransferRequest
	Index() int
	Amount() int
	TotalAmount() int
}

// SignedTransferRequestService represents the signed transfer request service
type SignedTransferRequestService interface {
	Save(req SignedTransferRequest) error
	RetrieveByID(id *uuid.UUID) (SignedTransferRequest, error)
	RetrieveByFromWalletID(fromWalletID *uuid.UUID, index int, amount int) (SignedTransferRequestPartialSet, error)
}

/*
 * TransferRequestVote
 */

// TransferRequestVote represents a transfer request vote
type TransferRequestVote interface {
	ID() *uuid.UUID
	Request() SignedTransferRequest
	Voter() User
	IsApproved() bool
}

// TransferRequestVotePartialSet represents the transfer request vote partial set
type TransferRequestVotePartialSet interface {
	Votes() []TransferRequestVote
	Index() int
	Amount() int
	TotalAmount() int
}

// TransferRequestVoteService represents a transfer request vote service
type TransferRequestVoteService interface {
	Save(vote TransferRequestVote) error
	Delete(vote TransferRequestVote) error
	RetrieveByID(id *uuid.UUID) (TransferRequestVote, error)
	//FromStoredToUserRequestVote(vote *storedUserRequestVote) (UserRequestVote, error)
	RetrieveByVoterIDAndTransferRequestID(voterID *uuid.UUID, requestID *uuid.UUID) (TransferRequestVote, error)
	RetrieveByTransferRequestID(requestID *uuid.UUID, index int, amount int) (TransferRequestVotePartialSet, error)
	RetrieveByRequesterWalletID(walletID *uuid.UUID, index int, amount int) (TransferRequestVotePartialSet, error)
	RetrieveByVoterID(voterID *uuid.UUID, index int, amount int) (TransferRequestVotePartialSet, error)
	RetrieveByVoterWalletID(walletID *uuid.UUID, index int, amount int) (TransferRequestVotePartialSet, error)
}

/*
 * Transfer
 */

// Transfer represents a transfer of token that can be claimed
type Transfer interface {
	ID() *uuid.UUID
	Amount() int
	Content() string
	PubKey() crypto.PublicKey
}

// TransferPartialSet represents a transfer partial set
type TransferPartialSet interface {
	Transfers() []Transfer
	Index() int
	Amount() int
	TotalAmount() int
}

// TransferService represents a transfer service
type TransferService interface {
	Save(trx Transfer) error
	Retrieve(index int, amount int) (TransferPartialSet, error)
	RetrieveByID(id *uuid.UUID) (Transfer, error)
	RetrieveByPublicKey(pubKey crypto.PublicKey) (Transfer, error)
}

/*
 * TransferClaim
 */

// TransferClaim represents a claim of transfer
type TransferClaim interface {
	ID() *uuid.UUID
	DepositTo() Wallet
	SignedContent() crypto.RingSignature
	Amount() int
}

// TransferClaimPartialSet represents a TransferClaimPartialSet instance
type TransferClaimPartialSet interface {
	Claims() []TransferClaim
	Index() int
	Amount() int
	TotalAmount() int
}

// TransferClaimService represents the TransferClaimService
type TransferClaimService interface {
	Save(claim TransferClaim) error
	RetrieveByID(id *uuid.UUID) (TransferClaim, error)
	RetrieveByToWalletID(toWalletID *uuid.UUID, index int, amount int) (TransferClaimPartialSet, error)
}

/*
 * Pledge
 */

// Pledge represents a pledge of tokens to a wallet
type Pledge interface {
	ID() *uuid.UUID
	From() Wallet
	To() Wallet
	Amount() int
}

// PledgePartialSet represents the pledge partial set
type PledgePartialSet interface {
	Pledges() []Pledge
	Index() int
	Amount() int
	TotalAmount() int
}

// PledgeService represents a pledge service
type PledgeService interface {
	Save(pledge Pledge) error
	Populate(stored *storedPledge) (Pledge, error)
	RetrieveByID(id *uuid.UUID) (Pledge, error)
	RetrieveByFromWalletID(fromWalletID *uuid.UUID, index int, amount int) (PledgePartialSet, error)
	RetrieveByToWalletID(toWalletID *uuid.UUID, index int, amount int) (PledgePartialSet, error)
}

/*
 * Balance
 */

// PledgeBalanceService represents the pledge balance service
type PledgeBalanceService interface {
	RetrieveByFromWalletID(fromWalletID *uuid.UUID, index int, amount int) (int, error)
	RetrieveByToWalletID(toWalletID *uuid.UUID, index int, amount int) (int, error)
}

/*
 * Validator
 */

// Validator represents a validator
type Validator interface {
	ID() *uuid.UUID
	Wallet() Wallet
	PubKey() tcrypto.PubKey
}

// ValidatorService represents the validator service
type ValidatorService interface {
	Retrieve(index int, amount int) ([]Validator, error)
	RetrieveByID(id *uuid.UUID) (Validator, error)
	Delete(val Validator) error
}
