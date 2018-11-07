package wallet

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

// Wallet represents a wallet
type Wallet interface {
	ID() *uuid.UUID
	ConcensusNeeded() int
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
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if wallet, ok := ins.(Wallet); ok {
					out := createStoredWallet(wallet)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllWalletKeyname(),
				}, nil
			},
		})
	},
}
