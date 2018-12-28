package meta

import (
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/chain"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/deposit"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
)

func createMeta() (meta.Meta, error) {
	// create core metadata:
	walletMetaData := wallet.SDKFunc.CreateMetaData()

	// create the representations:
	addressRepresentation := address.SDKFunc.CreateRepresentation()
	offerRepresentation := offer.SDKFunc.CreateRepresentation()
	depositRepresentation := deposit.SDKFunc.CreateRepresentation()
	chainRepresentation := chain.SDKFunc.CreateRepresentation()

	// get the metadata:
	addressMetaData := addressRepresentation.MetaData()
	offerMetaData := offerRepresentation.MetaData()
	depositMetaData := depositRepresentation.MetaData()
	chainMetaData := chainRepresentation.MetaData()

	// create the meta:
	met := meta.SDKFunc.Create(meta.CreateParams{
		AdditionalRead: map[string]entity.MetaData{
			addressMetaData.Keyname(): addressMetaData,
			offerMetaData.Keyname():   offerMetaData,
			depositMetaData.Keyname(): depositMetaData,
			chainMetaData.Keyname():   chainMetaData,
		},
	})

	// the address must be voted by the wallet owners:
	addedAddressToWalletOwnersErr := met.AddToWriteOnEntityRequest(walletMetaData, addressRepresentation)
	if addedAddressToWalletOwnersErr != nil {
		return nil, addedAddressToWalletOwnersErr
	}

	// the offer must be voted by the wallet owners:
	addedOfferToWalletOwnersErr := met.AddToWriteOnEntityRequest(walletMetaData, offerRepresentation)
	if addedOfferToWalletOwnersErr != nil {
		return nil, addedOfferToWalletOwnersErr
	}

	// the deposit must be voted by the wallet owners:
	addedDepositToWalletOwnersErr := met.AddToWriteOnEntityRequest(walletMetaData, depositRepresentation)
	if addedDepositToWalletOwnersErr != nil {
		return nil, addedDepositToWalletOwnersErr
	}

	// returns:
	return met, nil
}
