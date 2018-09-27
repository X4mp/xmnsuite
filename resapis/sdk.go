package restapis

import (
	"time"

	crypto "github.com/xmnservices/xmnsuite/crypto"
)

// Server represents a server api
type Server interface {
	Start() error
	Stop() error
}

// Token represents an authenticated token
type Token interface {
	Method() string
	RequestURI() string
	Data() map[string][]string
	Hash() string
}

// Client represents the client sdk
type Client interface {
	CreateAccount(name string, seedwords []string) error
	RetrieveAccounts() ([]Account, error)
	RetrieveAccountByName(name string) (Account, error)
}

// Account represents an aplication account
type Account interface {
	Name() string
	EncryptedPK() string
	CreatedOn() time.Time
	DecryptPK(seedWords []string) crypto.PrivateKey
}
