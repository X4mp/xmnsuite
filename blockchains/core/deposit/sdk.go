package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

// Deposit represents the initial deposit
type Deposit interface {
	ID() *uuid.UUID
	To() wallet.Wallet
	Token() token.Token
	Amount() int
}

// Normalized represents the normalized deposit
type Normalized interface {
}

// SDKFunc represents the Deposit SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
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

				walletRepresentation := wallet.SDKFunc.CreateRepresentation()
				tokRepresentation := token.SDKFunc.CreateRepresentation()

				if deposit, ok := ins.(Deposit); ok {
					// try to retrieve the wallet:
					toWallet := deposit.To()
					_, retToWalletErr := rep.RetrieveByID(walletRepresentation.MetaData(), toWallet.ID())
					if retToWalletErr != nil {
						// save the wallet:
						saveErr := service.Save(toWallet, walletRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}

					// try to retrieve the token:
					tok := deposit.Token()
					_, retTokErr := rep.RetrieveByID(tokRepresentation.MetaData(), tok.ID())
					if retTokErr != nil {
						// save the token:
						saveErr := service.Save(tok, tokRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
