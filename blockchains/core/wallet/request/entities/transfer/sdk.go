package transfer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
)

// Transfer represents a transfer of token that can be claimed
type Transfer interface {
	ID() *uuid.UUID
	Withdrawal() withdrawal.Withdrawal
	Deposit() deposit.Deposit
}

// SDKFunc represents the Transfer SDK func
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
				if trans, ok := ins.(Transfer); ok {
					out := createStorableTransfer(trans)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Transfer instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllTransfersKeyname(),
				}, nil
			},
		})
	},
}
