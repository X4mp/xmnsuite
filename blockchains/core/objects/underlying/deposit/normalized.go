package deposit

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type normalizedDeposit struct {
	ID     string            `json:"id"`
	To     wallet.Normalized `json:"to"`
	Amount int               `json:"amount"`
}

func createNormalizedDeposit(ins Deposit) (*normalizedDeposit, error) {
	normalizedToWallet, normalizedToWalletErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.To())
	if normalizedToWalletErr != nil {
		return nil, normalizedToWalletErr
	}

	out := normalizedDeposit{
		ID:     ins.ID().String(),
		To:     normalizedToWallet,
		Amount: ins.Amount(),
	}

	return &out, nil
}
