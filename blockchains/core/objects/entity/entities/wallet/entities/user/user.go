package user

import (
	"errors"
	"fmt"
	"regexp"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

type user struct {
	UUID  *uuid.UUID       `json:"id"`
	Nme   string           `json:"name"`
	PKey  crypto.PublicKey `json:"pubkey"`
	Shres int              `json:"shares"`
	Wal   wallet.Wallet    `json:"wallet"`
	Ref   wallet.Wallet    `json:"referral"`
}

func createUser(id *uuid.UUID, nme string, pubKey crypto.PublicKey, shares int, wallet wallet.Wallet) (User, error) {
	return createUserWithReferral(id, nme, pubKey, shares, wallet, nil)
}

func createUserWithReferral(id *uuid.UUID, nme string, pubKey crypto.PublicKey, shares int, wallet wallet.Wallet, ref wallet.Wallet) (User, error) {

	pattern, patternErr := regexp.Compile("[a-zA-Z0-9-]{3,}")
	if patternErr != nil {
		return nil, patternErr
	}

	found := pattern.FindString(nme)
	if found != nme {
		str := fmt.Sprintf("the name (%s) must only contain letters, numbers and hyphens (-) with at least 3 characters", nme)
		return nil, errors.New(str)
	}

	out := user{
		UUID:  id,
		Nme:   nme,
		PKey:  pubKey,
		Shres: shares,
		Wal:   wallet,
		Ref:   ref,
	}

	return &out, nil
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
		return createUser(&id, ins.Name, pubKey, ins.Shares, wal)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", walIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (app *user) ID() *uuid.UUID {
	return app.UUID
}

// Name returns the name of the user
func (app *user) Name() string {
	return app.Nme
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

// HasBeenReferred returns true if the user has been referred, false otherwise
func (app *user) HasBeenReferred() bool {
	return app.Ref != nil
}

// Referral returns the referral, if any
func (app *user) Referral() wallet.Wallet {
	return app.Ref
}
