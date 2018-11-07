package withdrawal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

// Withdrawal represents a withdrawal
type Withdrawal interface {
	ID() *uuid.UUID
	From() wallet.Wallet
	Token() token.Token
	Amount() int
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	From   wallet.Wallet
	Tok    token.Token
	Amount int
}

// CreateRepresentationParams represents the CreateRepresentation params
type CreateRepresentationParams struct {
	WalletMetaData       entity.MetaData
	WalletRepresentation entity.Representation
}

// SDKFunc represents the Withdrawal SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Withdrawal
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func(params CreateRepresentationParams) entity.Representation
}{
	Create: func(params CreateParams) Withdrawal {
		out := createWithdrawal(params.ID, params.From, params.Tok, params.Amount)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func(params CreateRepresentationParams) entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if withdrawal, ok := ins.(Withdrawal); ok {
					out := createStorableWithdrawal(withdrawal)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if withdrawal, ok := ins.(Withdrawal); ok {
					base := retrieveAllWithdrawalsKeyname()
					return []string{
						base,
						fmt.Sprintf("%s:by_from_wallet_id:%s", base, withdrawal.From().ID().String()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
				return nil, errors.New(str)

			},
			Sync: func(rep entity.Repository, service entity.Service, ins entity.Entity) error {
				if withdrawal, ok := ins.(Withdrawal); ok {
					// try to retrieve the wallet:
					toWallet := withdrawal.From()
					_, retToWalletErr := rep.RetrieveByID(params.WalletMetaData, toWallet.ID())
					if retToWalletErr != nil {
						// save the wallet:
						saveErr := service.Save(toWallet, params.WalletRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
