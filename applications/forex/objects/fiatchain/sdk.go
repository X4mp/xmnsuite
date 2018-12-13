package fiatchain

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
)

// FiatChain represents a blockchain that holds fiat deposit
type FiatChain interface {
	ID() *uuid.UUID
	Seeds() []string
	Deposits() []deposit.Deposit
}

// Normalized represents a normalized fiatchain
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID       *uuid.UUID
	Seeds    []string
	Deposits []deposit.Deposit
}

// SDKFunc represents the FiatChain SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) FiatChain
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) FiatChain {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createFiatChain(params.ID, params.Seeds, params.Deposits)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if fc, ok := ins.(FiatChain); ok {
					out := createStorableFiatChain(fc)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid FiatChain instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if fc, ok := ins.(FiatChain); ok {

					deps := fc.Deposits()
					keynames := []string{
						retrieveAllFiatChainsKeyname(),
					}

					for _, oneDep := range deps {
						keynames = append(keynames, retrieveFiatChainByDepositIDKeyname(oneDep.ID()))
					}

					return keynames, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid FiatChain instance", ins.ID().String())
				return nil, errors.New(str)
			},
		})
	},
}
