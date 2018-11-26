package pledge

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/withdrawal"
)

func retrieveAllPledgesKeyname() string {
	return "pledges"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Pledge",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storablePledge); ok {
				return fromStorableToPledge(rep, storable)
			}

			ptr := new(normalizedPledge)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createPledgeFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if pledge, ok := ins.(Pledge); ok {
				return createNormalizedPledge(pledge)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Pledge instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedPledge); ok {
				return createPledgeFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Pledge instance")
		},
		EmptyStorable:   new(storablePledge),
		EmptyNormalized: new(normalizedPledge),
	})
}

func fromStorableToPledge(repository entity.Repository, storable *storablePledge) (Pledge, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		str := fmt.Sprintf("the given storable Pledge ID (%s) is invalid: %s", storable.ID, idErr.Error())
		return nil, errors.New(str)
	}

	withdrawalID, withdrawalIDErr := uuid.FromString(storable.FromWithdrawalID)
	if withdrawalIDErr != nil {
		str := fmt.Sprintf("the given storable Pledge Withdrawal ID (%s) is invalid: %s", storable.FromWithdrawalID, withdrawalIDErr.Error())
		return nil, errors.New(str)
	}

	walletID, walletIDErr := uuid.FromString(storable.ToWalletID)
	if walletIDErr != nil {
		str := fmt.Sprintf("the given storable Pledge Wallet ID (%s) is invalid: %s", storable.ToWalletID, walletIDErr.Error())
		return nil, errors.New(str)
	}

	// retrieve the withdrawal:
	withdrawalMetaData := withdrawal.SDKFunc.CreateMetaData()
	fromIns, fromInsErr := repository.RetrieveByID(withdrawalMetaData, &withdrawalID)
	if fromInsErr != nil {
		return nil, fromInsErr
	}

	// retrieve the wallet:
	walletMetaData := wallet.SDKFunc.CreateMetaData()
	toIns, toInsErr := repository.RetrieveByID(walletMetaData, &walletID)
	if toInsErr != nil {
		return nil, toInsErr
	}

	if from, ok := fromIns.(withdrawal.Withdrawal); ok {
		if to, ok := toIns.(wallet.Wallet); ok {
			out := createPledge(&id, from, to)
			return out, nil
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", walletID.String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Withdrawal instance", withdrawalID.String())
	return nil, errors.New(str)
}
