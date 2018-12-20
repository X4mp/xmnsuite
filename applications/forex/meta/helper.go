package meta

import (
	"github.com/xmnservices/xmnsuite/applications/forex/objects/bank"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/deposit"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/fiatchain"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
)

func createMeta() (meta.Meta, error) {
	// create core metadata:
	walletMetaData := wallet.SDKFunc.CreateMetaData()
	tokenMetaData := token.SDKFunc.CreateMetaData()

	// create the representations:
	bankRepresentation := bank.SDKFunc.CreateRepresentation()
	categoryRepresentation := category.SDKFunc.CreateRepresentation()
	currencyRepresentation := currency.SDKFunc.CreateRepresentation()
	depositRepresentation := deposit.SDKFunc.CreateRepresentation()
	fiatChainRepresentation := fiatchain.SDKFunc.CreateRepresentation()

	// get the metadata:
	bankMetaData := bankRepresentation.MetaData()
	categoryMetaData := categoryRepresentation.MetaData()
	currencyMetaData := currencyRepresentation.MetaData()
	depositMetaData := depositRepresentation.MetaData()
	fiatChainMetaData := fiatChainRepresentation.MetaData()

	// create the meta:
	met := meta.SDKFunc.Create(meta.CreateParams{
		AdditionalRead: map[string]entity.MetaData{
			bankMetaData.Keyname():      bankMetaData,
			categoryMetaData.Keyname():  categoryMetaData,
			currencyMetaData.Keyname():  currencyMetaData,
			depositMetaData.Keyname():   depositMetaData,
			fiatChainMetaData.Keyname(): fiatChainMetaData,
		},
	})

	// the category must be voted by the token holders:
	addedCategoriesToTokErr := met.AddToWriteOnEntityRequest(tokenMetaData, categoryRepresentation)
	if addedCategoriesToTokErr != nil {
		return nil, addedCategoriesToTokErr
	}

	// the currency must be voted by the token holders:
	addedCurrenciesToTokErr := met.AddToWriteOnEntityRequest(tokenMetaData, currencyRepresentation)
	if addedCurrenciesToTokErr != nil {
		return nil, addedCurrenciesToTokErr
	}

	// the bank must be voted by the wallet owners:
	addedBankToWalletOwnersErr := met.AddToWriteOnEntityRequest(walletMetaData, bankRepresentation)
	if addedBankToWalletOwnersErr != nil {
		return nil, addedBankToWalletOwnersErr
	}

	// the deposit must be voted by the wallet owners:
	addedDepositToWalletOwnersErr := met.AddToWriteOnEntityRequest(walletMetaData, depositRepresentation)
	if addedDepositToWalletOwnersErr != nil {
		return nil, addedDepositToWalletOwnersErr
	}

	// returns:
	return met, nil
}
