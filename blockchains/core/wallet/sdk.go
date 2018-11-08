package wallet

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

// Wallet represents a wallet
type Wallet interface {
	ID() *uuid.UUID
	ConcensusNeeded() int
}

// Normalized represents a normalized wallet
type Normalized interface {
}

// SDKFunc represents the Wallet SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
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
