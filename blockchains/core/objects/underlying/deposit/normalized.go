package deposit

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type normalizedDeposit struct {
	ID     string            `json:"id"`
	To     wallet.Normalized `json:"to"`
	Token  token.Normalized  `json:"token"`
	Amount int               `json:"amount"`
}

func createNormalizedDeposit(ins Deposit) (*normalizedDeposit, error) {
	normalizedToWallet, normalizedToWalletErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.To())
	if normalizedToWalletErr != nil {
		return nil, normalizedToWalletErr
	}

	normalizedToken, normalizedTokenErr := token.SDKFunc.CreateMetaData().Normalize()(ins.Token())
	if normalizedTokenErr != nil {
		return nil, normalizedTokenErr
	}

	out := normalizedDeposit{
		ID:     ins.ID().String(),
		To:     normalizedToWallet,
		Token:  normalizedToken,
		Amount: ins.Amount(),
	}

	return &out, nil
}
