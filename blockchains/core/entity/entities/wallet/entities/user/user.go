package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
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

func createUserFromNormalizedUser(ins *normalizedUser) (User, error) {
	id, idErr := uuid.FromString(ins.ID)
	if idErr != nil {
		return nil, idErr
	}

	pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
		PubKeyAsString: ins.PubKey,
	})

	walIns, walInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(ins.Wallet)
	if walInsErr != nil {
		return nil, walInsErr
	}

	if wal, ok := walIns.(wallet.Wallet); ok {
		out := createUser(&id, pubKey, ins.Shares, wal)
		return out, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", walIns.ID().String())
	return nil, errors.New(str)
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
