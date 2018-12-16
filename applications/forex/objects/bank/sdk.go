package bank

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
)

// Bank represents a bank
type Bank interface {
	ID() *uuid.UUID
	Pledge() pledge.Pledge
	Currency() currency.Currency
	Amount() int
	Price() int
}

// Interval represents an interval
type Interval interface {
	Min() int
	Max() int
}

// RetrieverCriteria represents a bank retriever criteria
type RetrieverCriteria interface {
	Amount() Interval
	Price() Interval
}

// Repository represents the bank repository
type Repository interface {
	RetrieveByPledgeID(id *uuid.UUID) (Bank, error)
	RetrieveSetByCurrencyID(id *uuid.UUID, criteria RetrieverCriteria) (Bank, error)
	RetrieveSetByCurrencyIDs(ids *uuid.UUID, criteria RetrieverCriteria) (Bank, error)
}

// Normalized represents a normalized bank
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID       *uuid.UUID
	Pledge   pledge.Pledge
	Currency currency.Currency
	Amount   int
	Price    int
}

// SDKFunc represents the Bank SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Bank
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Bank {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createBank(params.ID, params.Pledge, params.Currency, params.Amount, params.Price)
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
				if bnk, ok := ins.(Bank); ok {
					out := createStorableBank(bnk)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Bank instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if bnk, ok := ins.(Bank); ok {
					return []string{
						retrieveAllBanksKeyname(),
						retrieveBankByCurrencyIDKeyname(bnk.Currency().ID()),
						retrieveBankByPledgeIDKeyname(bnk.Pledge().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Bank instance", ins.ID().String())
				return nil, errors.New(str)
			},
		})
	},
}
