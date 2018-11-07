package user

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

type user struct {
	UUID  *uuid.UUID       `json:"id"`
	PKey  crypto.PublicKey `json:"pubkey"`
	Shres int              `json:"shares"`
	Wal   wallet.Wallet    `json:"wallet"`
}

func createUser(id *uuid.UUID, pubKey crypto.PublicKey, shares int, wallet wallet.Wallet) User {
	out := user{
		UUID:  id,
		PKey:  pubKey,
		Shres: shares,
		Wal:   wallet,
	}

	return &out
}

// ID returns the ID
func (app *user) ID() *uuid.UUID {
	return app.UUID
}

// PubKey returns the PublicKey
func (app *user) PubKey() crypto.PublicKey {
	return app.PKey
}

// Shares returns the shares
func (app *user) Shares() int {
	return app.Shres
}

// Wallet returns the wallet
func (app *user) Wallet() wallet.Wallet {
	return app.Wal
}