package sell

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/external"
)

func retrieveAllSellsKeyname() string {
	return "sells"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Sell",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableSell) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				fromID, fromIDErr := uuid.FromString(storable.FromPledgeID)
				if fromIDErr != nil {
					return nil, fromIDErr
				}

				depositToWalletID, depositToWalletIDErr := uuid.FromString(storable.DepositToWalletID)
				if depositToWalletIDErr != nil {
					return nil, depositToWalletIDErr
				}

				tokID, tokIDErr := uuid.FromString(storable.Wish.ExternalTokenID)
				if tokIDErr != nil {
					return nil, tokIDErr
				}

				// retrieve the pledge:
				pledgeMetaData := pledge.SDKFunc.CreateMetaData()
				pledgeIns, pledgeInsErr := rep.RetrieveByID(pledgeMetaData, &fromID)
				if pledgeInsErr != nil {
					return nil, pledgeInsErr
				}

				// retrieve the wallet external resource:
				externalMetaData := external.SDKFunc.CreateMetaData()
				externalWalletIns, externalWalletInsErr := rep.RetrieveByID(externalMetaData, &depositToWalletID)
				if externalWalletInsErr != nil {
					return nil, externalWalletInsErr
				}

				// retrieve the token external resource:
				externalTokIns, externalTokInsErr := rep.RetrieveByID(externalMetaData, &tokID)
				if externalTokInsErr != nil {
					return nil, externalTokInsErr
				}

				if pledge, ok := pledgeIns.(pledge.Pledge); ok {
					if extWallet, ok := externalWalletIns.(external.External); ok {
						if extTok, ok := externalTokIns.(external.External); ok {
							wish := createWish(extTok, storable.Wish.Amount)
							out := createSell(&id, pledge, wish, extWallet)
							return out, nil
						}

						str := fmt.Sprintf("the entity (ID: %s) is not a valid External (token) instance", tokID.String())
						return nil, errors.New(str)
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid External (wallet) instance", depositToWalletID.String())
					return nil, errors.New(str)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", fromID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableSell); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableSell)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableSell),
	})
}
