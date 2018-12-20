package sell

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/external"
)

// Wish is an external wanted proposition
type Wish interface {
	Token() external.External
	Amount() int
}

// Sell represents a sell order
type Sell interface {
	ID() *uuid.UUID
	From() pledge.Pledge
	Wish() Wish
	DepositToWallet() external.External
}

// Repository represents the sell repository
type Repository interface {
	RetrieveMatch(with Wish) (Sell, error)
	RetrieveMatches(wish Wish) (entity.PartialSet, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the create params
type CreateParams struct {
	ID              *uuid.UUID
	From            pledge.Pledge
	Wish            Wish
	DepositToWallet external.External
}

// SDKFunc represents the Sell SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Sell
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Sell {
		out := createSell(params.ID, params.From, params.Wish, params.DepositToWallet)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if sell, ok := ins.(Sell); ok {
					out := createStorableSell(sell)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Sell instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllSellsKeyname(),
				}, nil
			},
		})
	},
}
