package wallet

import (
	uuid "github.com/satori/go.uuid"
)

type wallet struct {
	UUID    *uuid.UUID `json:"id"`
	CNeeded int        `json:"concensus_needed"`
}

func createWallet(id *uuid.UUID, concensusNeeded int) Wallet {
	out := wallet{
		UUID:    id,
		CNeeded: concensusNeeded,
	}

	return &out
}

func createWalletFromStorable(storable *storableWallet) (Wallet, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	out := createWallet(&id, storable.CNeeded)
	return out, nil
}

// ID returns the ID
func (app *wallet) ID() *uuid.UUID {
	return app.UUID
}

// ConcensusNeeded returns the concensus needed
func (app *wallet) ConcensusNeeded() int {
	return app.CNeeded
}
