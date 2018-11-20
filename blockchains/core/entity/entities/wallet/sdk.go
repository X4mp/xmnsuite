package wallet

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Wallet represents a wallet
type Wallet interface {
	ID() *uuid.UUID
	Creator() crypto.PublicKey
	ConcensusNeeded() int
}

// Normalized represents a normalized wallet
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID              *uuid.UUID
	Creator         crypto.PublicKey
	ConcensusNeeded int
}

// SDKFunc represents the Wallet SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Wallet
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Wallet {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createWallet(params.ID, params.Creator, params.ConcensusNeeded)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met:        createMetaData(),
			ToStorable: toData,
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllWalletKeyname(),
				}, nil
			},
		})
	},
}
