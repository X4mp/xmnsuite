package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
)

// Pledge represents a pledge
type Pledge interface {
	ID() *uuid.UUID
	From() withdrawal.Withdrawal
	To() wallet.Wallet
	Amount() int
}

// SDKFunc represents the Pledge SDK func
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
				if pledge, ok := ins.(Pledge); ok {
					out := createStorablePledge(pledge)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Pledge instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllPledgesKeyname(),
				}, nil
			},
		})
	},
}
