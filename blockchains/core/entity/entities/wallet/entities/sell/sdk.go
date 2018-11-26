package sell

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/external"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
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

// Daemon represents the daemon exchange
type Daemon interface {
	Start() error
	Stop() error
}

// SDKFunc represents the Exchange SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
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
