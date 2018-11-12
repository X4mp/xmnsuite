package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
)

// Pledge represents a pledge
type Pledge interface {
	ID() *uuid.UUID
	From() withdrawal.Withdrawal
	To() wallet.Wallet
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
				if pledge, ok := ins.(Pledge); ok {
					base := retrieveAllPledgesKeyname()
					return []string{
						base,
						fmt.Sprintf("%s:by_from_withdrawal_id:%s", base, pledge.From().ID().String()),
						fmt.Sprintf("%s:by_to_wallet_id:%s", base, pledge.To().ID().String()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Sync: func(rep entity.Repository, service entity.Service, ins entity.Entity) error {
				withdrawalRepresentation := withdrawal.SDKFunc.CreateRepresentation()
				walletRepresentation := wallet.SDKFunc.CreateRepresentation()

				if pledge, ok := ins.(Pledge); ok {
					// try to retrieve the withdrawal, send an error if it exists:
					from := pledge.From()
					_, retFromErr := rep.RetrieveByID(withdrawalRepresentation.MetaData(), from.ID())
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
					_, retToErr := rep.RetrieveByID(walletRepresentation.MetaData(), to.ID())
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
