package wallet

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
)

type wallet struct {
	UUID          *uuid.UUID       `json:"id"`
	CreatorPubKey crypto.PublicKey `json:"creator"`
	CNeeded       int              `json:"concensus_needed"`
}

func createWallet(id *uuid.UUID, creatorPubKey crypto.PublicKey, concensusNeeded int) Wallet {
	out := wallet{
		UUID:          id,
		CreatorPubKey: creatorPubKey,
		CNeeded:       concensusNeeded,
	}

	return &out
}

func createWalletFromStorable(storable *storableWallet) (Wallet, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
		PubKeyAsString: storable.Creator,
	})

	out := createWallet(&id, pubKey, storable.CNeeded)
	return out, nil
}

// ID returns the ID
func (app *wallet) ID() *uuid.UUID {
	return app.UUID
}

// Creator returns the creator public key
func (app *wallet) Creator() crypto.PublicKey {
	return app.CreatorPubKey
}

// ConcensusNeeded returns the concensus needed
func (app *wallet) ConcensusNeeded() int {
	return app.CNeeded
}
