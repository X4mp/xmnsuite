package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
)

func retrieveAllDepositsKeyname() string {
	return "deposits"
}

func retrieveDepositsByTokenIDKeyname(tokenID *uuid.UUID) string {
	base := retrieveAllDepositsKeyname()
	return fmt.Sprintf("%s:by_token_id:%s", base, tokenID.String())
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

				tokenID, tokenIDErr := uuid.FromString(storable.TokenID)
				if tokenIDErr != nil {
					str := fmt.Sprintf("the given storable Deposit Token ID (%s) is invalid: %s", storable.TokenID, tokenIDErr.Error())
					return nil, errors.New(str)
				}

				// retrieve the wallet:
				walletMetaData := wallet.SDKFunc.CreateMetaData()
				walletIns, walletInsErr := rep.RetrieveByID(walletMetaData, &toWalletID)
				if walletInsErr != nil {
					return nil, walletInsErr
				}

				// retrieve the token:
				tokenMetaData := token.SDKFunc.CreateMetaData()
				tokenIns, tokenInsErr := rep.RetrieveByID(tokenMetaData, &tokenID)
				if tokenInsErr != nil {
					return nil, tokenInsErr
				}

				if wal, ok := walletIns.(wallet.Wallet); ok {
					if tok, ok := tokenIns.(token.Token); ok {
						return createDeposit(&id, wal, tok, storable.Amount)
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Token instance", tokenID.String())
					return nil, errors.New(str)
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

func toData(dep Deposit) *Data {
	out := Data{
		ID:     dep.ID().String(),
		To:     wallet.SDKFunc.ToData(dep.To()),
		Token:  token.SDKFunc.ToData(dep.Token()),
		Amount: dep.Amount(),
	}

	return &out
}

func toDataSet(ins entity.PartialSet) (*DataSet, error) {
	data := []*Data{}
	instances := ins.Instances()
	for _, oneIns := range instances {
		if dep, ok := oneIns.(Deposit); ok {
			data = append(data, toData(dep))
			continue
		}

		str := fmt.Sprintf("at least one of the elements (ID: %s) in the entity partial set is not a valid Deposit instance", oneIns.ID().String())
		return nil, errors.New(str)
	}

	out := DataSet{
		Index:       ins.Index(),
		Amount:      ins.Amount(),
		TotalAmount: ins.TotalAmount(),
		IsLast:      ins.IsLast(),
		Deposits:    data,
	}

	return &out, nil
}
