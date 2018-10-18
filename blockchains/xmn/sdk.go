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
	To() Wallet
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
	ConcensusNeeded() float64
}

// WalletPartialSet represents a wallet partial set
type WalletPartialSet interface {
	Wallets() []Wallet
	Index() int
	Amount() int
	TotalAmount() int
}

// AddUserToWalletRequests represents a wallet with its current add-user-request + votes
type AddUserToWalletRequests interface {
	Wallet() Wallet
	Req() AddUserToWalletRequest
	Votes() []AddUserToWalletRequestVote
	IsApproved() (bool, bool)
}

// AddUserToWalletRequest represents an add-user-request
type AddUserToWalletRequest interface {
	ID() *uuid.UUID
	Wallet() Wallet
	User() User
}

// AddUserToWalletRequestVote represents an add-user-to-wallet-request-vote
type AddUserToWalletRequestVote interface {
	Request() AddUserToWalletRequest
	IsAccepted() bool
}

// DelUserFromWalletRequests represents a wallet with its current delete-user-request + votes
type DelUserFromWalletRequests interface {
	Wallet() Wallet
	Req() DelUserFromWalletRequest
	Votes() []DelUserFromWalletRequestVote
	IsApproved() (bool, bool)
}

// DelUserFromWalletRequest represents a delete-user-from-wallet-request
type DelUserFromWalletRequest interface {
	ID() *uuid.UUID
	User() User
}

// DelUserFromWalletRequestVote represents a delete-user-from-wallet-request-vote
type DelUserFromWalletRequestVote interface {
	Request() DelUserFromWalletRequest
	IsAccepted() bool
}

// DeleteWalletRequest represents a delete-wallet-request
type DeleteWalletRequest interface {
	Wallet() Wallet
}

// DeleteWalletRequestVote represents a delete-wallet-request-vote
type DeleteWalletRequestVote interface {
	Request() DeleteWalletRequest
	IsAccepted() bool
}

// DelWalletRequests represents a wallet with its current delete-request + votes
type DelWalletRequests interface {
	Wallet() Wallet
	Req() DeleteWalletRequest
	Votes() []DeleteWalletRequestVote
	IsApproved() (bool, bool)
}

// WalletService represents the wallet service
type WalletService interface {
	Save(wallet Wallet) error
	SaveAddUserToWalletRequest(obj AddUserToWalletRequest) error
	SaveAddUserToWalletRequestVote(obj AddUserToWalletRequestVote) error
	SaveDeleteUserFromWalletRequest(obj DelUserFromWalletRequest) error
	SaveDeleteUserFromWalletRequestVote(obj DelUserFromWalletRequestVote) error
	SaveDeleteWalletRequest(obj DeleteWalletRequest) error
	SaveDeleteWalletRequestVote(obj DeleteWalletRequestVote) error
	Retrieve(index int, amount int) (WalletPartialSet, error)
	RetrieveByID(id *uuid.UUID) (Wallet, error)
	RetrieveAddUserToWalletRequests(index int, amount int) ([]AddUserToWalletRequests, error)
	RetrieveAddUserToWalletRequestsByWalletID(id *uuid.UUID) (AddUserToWalletRequests, error)
	RetrieveDelUserFromWalletRequests(index int, amount int) ([]DelUserFromWalletRequests, error)
	RetrieveDelUserFromWalletRequestsByWalletID(id *uuid.UUID) (DelUserFromWalletRequests, error)
	RetrieveDelWalletRequests(index int, amount int) ([]DelWalletRequests, error)
	RetrieveDelWalletRequestsByWalletID(id *uuid.UUID) (DelWalletRequests, error)
}

/*
 * User
 */

// User represents a user
type User interface {
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
	RetrieveByWalletID(walletID *uuid.UUID, index int, amount int)
}

/*
 * Validator
 */

// Validator represents a validator
type Validator interface {
	ID() *uuid.UUID
	Wallet() Wallet
	PubKey() tcrypto.PubKey
	Pow() int
}

// ValidatorService represents the validator service
type ValidatorService interface {
	Retrieve(index int, amount int) ([]Validator, error)
	RetrieveByID(id *uuid.UUID) (Validator, error)
	Delete(val Validator) error
}
