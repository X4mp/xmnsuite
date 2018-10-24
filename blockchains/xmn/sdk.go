package xmn

import (
	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/crypto"
)

/*
 * Genesis
 */

// Genesis represents the genesis instance
type Genesis interface {
	GazPricePerKb() int
	MaxAmountOfValidators() int
	Deposit() InitialDeposit
	Token() Token
}

// GenesisService represents the init service
type GenesisService interface {
	Save(obj Genesis) error
	Retrieve() (Genesis, error)
}

/*
 * InitialDeposit
 */

// InitialDeposit represents the initial deposit
type InitialDeposit interface {
	To() User
	Amount() int
}

// InitialDepositService represents the initial deposit service
type InitialDepositService interface {
	Retrieve() (InitialDeposit, error)
	Save(initialDep InitialDeposit) error
}

/*
 * Token
 */

// Token represents the token
type Token interface {
	Symbol() string
	Name() string
	Description() string
}

// TokenService represents the token service
type TokenService interface {
	Retrieve() (Token, error)
	Save(tok Token) error
}

/*
 * Wallet
 */

// Wallet represents a wallet
type Wallet interface {
	ID() *uuid.UUID
	ConcensusNeeded() int
}

// WalletPartialSet represents a wallet partial set
type WalletPartialSet interface {
	Wallets() []Wallet
	Index() int
	Amount() int
	TotalAmount() int
}

// WalletService represents the wallet service
type WalletService interface {
	Save(wallet Wallet) error
	Retrieve(index int, amount int) (WalletPartialSet, error)
	RetrieveByID(id *uuid.UUID) (Wallet, error)
}

/*
 * UserRequest
 */

// UserRequest represents a user request
type UserRequest interface {
	User() User
}

// UserRequestPartialSet represents the user request partial set
type UserRequestPartialSet interface {
	Requests() []UserRequest
	Index() int
	Amount() int
	TotalAmount() int
}

// UserRequestService represents a user request service
type UserRequestService interface {
	Save(req UserRequest) error
	Delete(req UserRequest) error
	RetrieveByID(id *uuid.UUID) (UserRequest, error)
	FromStoredToUserRequest(req *storedUserRequest) (UserRequest, error)
	RetrieveByPubkeyAndWalletID(pubKey crypto.PublicKey, walletID *uuid.UUID) (UserRequest, error)
	RetrieveByWalletID(walletID *uuid.UUID, index int, amount int) (UserRequestPartialSet, error)
	RetrieveByPubKey(pubKey crypto.PublicKey, index int, amount int) (UserRequestPartialSet, error)
}

/*
 * UserRequestVote
 */

// UserRequestVote represents a user request vote
type UserRequestVote interface {
	ID() *uuid.UUID
	Request() UserRequest
	Voter() User
	IsApproved() bool
}

// UserRequestVotePartialSet represents the user request vote partial set
type UserRequestVotePartialSet interface {
	UserRequestVotes() []UserRequestVote
	Index() int
	Amount() int
	TotalAmount() int
}

// UserRequestVoteService represents a user request vote service
type UserRequestVoteService interface {
	Save(vote UserRequestVote) error
	Delete(vote UserRequestVote) error
	RetrieveByID(id *uuid.UUID) (UserRequestVote, error)
	FromStoredToUserRequestVote(vote *storedUserRequestVote) (UserRequestVote, error)
	RetrieveByVoterIDAndUserRequestID(voterID *uuid.UUID, requestID *uuid.UUID) (UserRequestVote, error)
	RetrieveByUserRequestID(requestID *uuid.UUID, index int, amount int) (UserRequestVotePartialSet, error)
	RetrieveByRequesterWalletID(walletID *uuid.UUID, index int, amount int) (UserRequestVotePartialSet, error)
	RetrieveByVoterID(voterID *uuid.UUID, index int, amount int) (UserRequestVotePartialSet, error)
	RetrieveByVoterWalletID(walletID *uuid.UUID, index int, amount int) (UserRequestVotePartialSet, error)
}

/*
 * User
 */

// User represents a user
type User interface {
	ID() *uuid.UUID
	PubKey() crypto.PublicKey
	Shares() int
	Wallet() Wallet
}

// UserPartialSet represents the user partial set
type UserPartialSet interface {
	Users() []User
	Index() int
	Amount() int
	TotalAmount() int
}

// UserService represents a user service
type UserService interface {
	Save(usr User) error
	RetrieveByID(id *uuid.UUID) (User, error)
	RetrieveByWalletID(walletID *uuid.UUID, index int, amount int) (UserPartialSet, error)
	RetrieveByPubKey(pubKey crypto.PublicKey, index int, amount int) (UserPartialSet, error)
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
	RetrieveByID(id *uuid.UUID) (Pledge, error)
	FromStoredToPledge(stored *storedPledge) (Pledge, error)
	RetrieveByFromWalletID(fromWalletID *uuid.UUID, index int, amount int) (PledgePartialSet, error)
	RetrieveByToWalletID(toWalletID *uuid.UUID, index int, amount int) (PledgePartialSet, error)
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
