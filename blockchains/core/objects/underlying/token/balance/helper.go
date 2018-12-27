package balance

import (
	"errors"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
)

func toData(bal Balance) *Data {
	out := Data{
		On:     wallet.SDKFunc.ToData(bal.On()),
		Of:     token.SDKFunc.ToData(bal.Of()),
		Amount: bal.Amount(),
	}

	return &out
}

func convertToDataSet(tok token.Token, walletPS entity.PartialSet, rep Repository) (*DataSet, error) {
	walletIns := walletPS.Instances()
	balances := []*Data{}
	for _, oneIns := range walletIns {
		if oneWallet, ok := oneIns.(wallet.Wallet); ok {
			// retrieve balance:
			balance, balanceErr := rep.RetrieveByWalletAndToken(oneWallet, tok)
			if balanceErr != nil {
				return nil, balanceErr
			}

			balances = append(balances, toData(balance))
			continue
		}

		return nil, errors.New("there is at least 1 element in the entity partial set that is not a valid Wallet instance")
	}

	out := DataSet{
		Index:       walletPS.Index(),
		Amount:      walletPS.Amount(),
		TotalAmount: walletPS.TotalAmount(),
		IsLast:      walletPS.IsLast(),
		Balances:    balances,
	}

	return &out, nil
}
