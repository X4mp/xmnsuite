package xmn

import "github.com/xmnservices/xmnsuite/crypto"

type user struct {
	PKey  string `json:"pubkey"`
	Shres int    `json:"shares"`
	Wal   Wallet `json:"wallet"`
}

func createUser(pubKey crypto.PublicKey, shares int, wallet Wallet) User {
	out := user{
		PKey:  pubKey.String(),
		Shres: shares,
	}

	return &out
}

// PubKey returns the PublicKey
func (app *user) PubKey() crypto.PublicKey {
	return crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
		PubKeyAsString: app.PKey,
	})
}

// Shares returns the shares
func (app *user) Shares() int {
	return app.Shres
}

// Wallet returns the wallet
func (app *user) Wallet() Wallet {
	return app.Wal
}
