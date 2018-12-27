package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/bank"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

// Deposit represents a currency deposit to a bank
type Deposit interface {
	ID() *uuid.UUID
	Amount() int
	Bank() bank.Bank
}

// Normalized represents a normalized deposit
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	Bank   bank.Bank
	Amount int
}

// SDKFunc represents the Deposit SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Deposit
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Deposit {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createDeposit(params.ID, params.Amount, params.Bank)
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
				if dep, ok := ins.(Deposit); ok {
					out := createStorableDeposit(dep)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if dep, ok := ins.(Deposit); ok {
					return []string{
						retrieveAllDepositsKeyname(),
						retrieveDepositsByBankIDKeyname(dep.Bank().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)
			},
		})
	},
}
