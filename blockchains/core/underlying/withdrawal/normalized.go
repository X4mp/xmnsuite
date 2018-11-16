package withdrawal

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
)

type normalizedWithdrawal struct {
	ID     string            `json:"id"`
	From   wallet.Normalized `json:"from"`
	Token  token.Normalized  `json:"token"`
	Amount int               `json:"amount"`
}

func createNormalizedWithdrawal(ins Withdrawal) (*normalizedWithdrawal, error) {
	from, fromErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.From())
	if fromErr != nil {
		return nil, fromErr
	}

	tok, tokErr := token.SDKFunc.CreateMetaData().Normalize()(ins.Token())
	if tokErr != nil {
		return nil, tokErr
	}

	out := normalizedWithdrawal{
		ID:     ins.ID().String(),
		From:   from,
		Token:  tok,
		Amount: ins.Amount(),
	}

	return &out, nil
}
