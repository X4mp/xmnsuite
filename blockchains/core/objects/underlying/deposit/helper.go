package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

func retrieveAllDepositsKeyname() string {
	return "deposits"
}

func retrieveDepositsByToWalletIDKeyname(toWalletID *uuid.UUID) string {
	base := retrieveAllDepositsKeyname()
	return fmt.Sprintf("%s:by_to_wallet_id:%s", base, toWalletID.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Deposit",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableDeposit) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					str := fmt.Sprintf("the given storable Deposit ID (%s) is invalid: %s", storable.ID, idErr.Error())
					return nil, errors.New(str)
				}

				toWalletID, toWalletIDErr := uuid.FromString(storable.ToWalletID)
				if toWalletIDErr != nil {
					str := fmt.Sprintf("the given storable Deposit Wallet ID (%s) is invalid: %s", storable.ToWalletID, toWalletIDErr.Error())
					return nil, errors.New(str)
				}

				// retrieve the wallet:
				walletMetaData := wallet.SDKFunc.CreateMetaData()
				walletIns, walletInsErr := rep.RetrieveByID(walletMetaData, &toWalletID)
				if walletInsErr != nil {
					return nil, walletInsErr
				}

				if wal, ok := walletIns.(wallet.Wallet); ok {
					return createDeposit(&id, wal, storable.Amount)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", toWalletID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableDeposit); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedDeposit)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createDepositFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if deposit, ok := ins.(Deposit); ok {
				out, outErr := createNormalizedDeposit(deposit)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedDeposit); ok {
				return createDepositFromNormalized(normalized)
			}

			return nil, errors.New("the given normalized instance cannot be converted to a Deposit instance")
		},
		EmptyStorable:   new(storableDeposit),
		EmptyNormalized: new(normalizedDeposit),
	})
}
