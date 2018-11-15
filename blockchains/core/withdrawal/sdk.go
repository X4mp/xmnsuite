package withdrawal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

// Withdrawal represents a withdrawal
type Withdrawal interface {
	ID() *uuid.UUID
	From() wallet.Wallet
	Token() token.Token
	Amount() int
}

// Normalized represents the normalized withdrawal
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	From   wallet.Wallet
	Token  token.Token
	Amount int
}

// SDKFunc represents the Withdrawal SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Withdrawal
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Withdrawal {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createWithdrawal(params.ID, params.From, params.Token, params.Amount)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
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
						fmt.Sprintf("%s:by_token_id:%s", base, withdrawal.Token().ID().String()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
				return nil, errors.New(str)

			},
		})
	},
}
