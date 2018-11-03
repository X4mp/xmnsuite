package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
	"github.com/xmnservices/xmnsuite/blockchains/framework/wallet"
)

// Deposit represents the initial deposit
type Deposit interface {
	ID() *uuid.UUID
	To() wallet.Wallet
	Amount() int
}

// CreateRepresentationParams represents the CreateRepresentation params
type CreateRepresentationParams struct {
	WalletMetaData       entity.MetaData
	WalletRepresentation entity.Representation
}

// SDKFunc represents the Deposit SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func(params CreateRepresentationParams) entity.Representation
}{
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func(params CreateRepresentationParams) entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if deposit, ok := ins.(Deposit); ok {
					out := createStorableDeposit(deposit)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if deposit, ok := ins.(Deposit); ok {
					base := retrieveAllDepositsKeyname()
					return []string{
						base,
						fmt.Sprintf("%s:by_to_wallet_id:%s", base, deposit.To().ID().String()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)

			},
			Sync: func(rep entity.Repository, service entity.Service, ins entity.Entity) error {
				if deposit, ok := ins.(Deposit); ok {
					// try to retrieve the wallet:
					toUser := deposit.To()
					_, retToUserErr := rep.RetrieveByID(params.WalletMetaData, toUser.ID())
					if retToUserErr != nil {
						// save the wallet:
						saveErr := service.Save(toUser, params.WalletRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
