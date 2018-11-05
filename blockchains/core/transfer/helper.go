package transfer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
	"github.com/xmnservices/xmnsuite/crypto"
)

func retrieveAllTransfersKeyname() string {
	return "transfers"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Transfer",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableTransfer) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				withdrawalID, withdrawalIDErr := uuid.FromString(storable.WithdrawalID)
				if withdrawalIDErr != nil {
					return nil, withdrawalIDErr
				}

				pubkey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
					PubKeyAsString: storable.PubKey,
				})

				// retrieve the transfer:
				withdrawalMetaData := withdrawal.SDKFunc.CreateMetaData()
				withdrawalIns, withdrawalInsErr := rep.RetrieveByID(withdrawalMetaData, &withdrawalID)
				if withdrawalInsErr != nil {
					return nil, withdrawalInsErr
				}

				if withdrawl, ok := withdrawalIns.(withdrawal.Withdrawal); ok {
					out := createTransfer(&id, withdrawl, storable.Content, pubkey)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", withdrawalID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableTransfer); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableTransfer)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableTransfer),
	})
}
