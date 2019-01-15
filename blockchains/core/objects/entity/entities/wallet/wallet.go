package wallet

import (
	"errors"
	"fmt"
	"regexp"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/crypto"
)

type wallet struct {
	UUID          *uuid.UUID       `json:"id"`
	Nme           string           `json:"name"`
	CreatorPubKey crypto.PublicKey `json:"creator"`
	CNeeded       int              `json:"concensus_needed"`
}

func createWallet(id *uuid.UUID, name string, creatorPubKey crypto.PublicKey, concensusNeeded int) (Wallet, error) {

	pattern, patternErr := regexp.Compile("[a-zA-Z0-9-]{3,}")
	if patternErr != nil {
		return nil, patternErr
	}

	found := pattern.FindString(name)
	if found != name {
		str := fmt.Sprintf("the name (%s) must only contain letters, numbers and hyphens (-) with at least 3 characters", name)
		return nil, errors.New(str)
	}

	out := wallet{
		UUID:          id,
		Nme:           name,
		CreatorPubKey: creatorPubKey,
		CNeeded:       concensusNeeded,
	}

	return &out, nil
}

func createWalletFromStorable(storable *storableWallet) (Wallet, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
		PubKeyAsString: storable.Creator,
	})

	return createWallet(&id, storable.Name, pubKey, storable.CNeeded)
}

// ID returns the ID
func (app *wallet) ID() *uuid.UUID {
	return app.UUID
}

// Name returns the name
func (app *wallet) Name() string {
	return app.Nme
}

// Creator returns the creator public key
func (app *wallet) Creator() crypto.PublicKey {
	return app.CreatorPubKey
}

// ConcensusNeeded returns the concensus needed
func (app *wallet) ConcensusNeeded() int {
	return app.CNeeded
}
