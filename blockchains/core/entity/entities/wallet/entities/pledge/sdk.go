package pledge

import (
	"bytes"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Pledge represents a pledge
type Pledge interface {
	ID() *uuid.UUID
	From() withdrawal.Withdrawal
	To() wallet.Wallet
}

// Normalized represents a normalized pledge
type Normalized interface {
}

// CreateParams represents the create params
type CreateParams struct {
	ID   *uuid.UUID
	From withdrawal.Withdrawal
	To   wallet.Wallet
}

// SDKFunc represents the Pledge SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Pledge
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Pledge {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createPledge(params.ID, params.From, params.To)
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
				if pledge, ok := ins.(Pledge); ok {
					out := createStorablePledge(pledge)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Pledge instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if pledge, ok := ins.(Pledge); ok {
					base := retrieveAllPledgesKeyname()
					return []string{
						base,
						fmt.Sprintf("%s:by_from_withdrawal_id:%s", base, pledge.From().ID().String()),
						fmt.Sprintf("%s:by_to_wallet_id:%s", base, pledge.To().ID().String()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Pledge instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)
				service := entity.SDKFunc.CreateService(ds)

				// create the representations:
				withdrawalRepresentation := withdrawal.SDKFunc.CreateRepresentation()
				walletRepresentation := wallet.SDKFunc.CreateRepresentation()

				if pledge, ok := ins.(Pledge); ok {
					// make sure the from wallet is not the same as the to wallet:
					if bytes.Compare(pledge.From().From().ID().Bytes(), pledge.To().ID().Bytes()) == 0 {
						str := fmt.Sprintf("the wallet of the from withdrawal (ID: %s) cannot be the same as the to wallet (ID: %s)", pledge.From().From().ID().String(), pledge.To().ID().String())
						return errors.New(str)
					}

					// try to retrieve the withdrawal, send an error if it exists:
					from := pledge.From()
					_, retFromErr := repository.RetrieveByID(withdrawalRepresentation.MetaData(), from.ID())
					if retFromErr == nil {
						str := fmt.Sprintf("the Pledge instance (ID: %s) contains a Withdrawal instance that already exists", from.ID().String())
						return errors.New(str)
					}

					// save the withdrawal:
					saveErr := service.Save(from, withdrawalRepresentation)
					if saveErr != nil {
						return saveErr
					}

					// try to retrieve the wallet:
					to := pledge.To()
					_, retToErr := repository.RetrieveByID(walletRepresentation.MetaData(), to.ID())
					if retToErr != nil {
						// save the wallet:
						saveErr := service.Save(to, walletRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Pledge instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
