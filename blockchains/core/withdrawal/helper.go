package withdrawal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
)

func retrieveAllWithdrawalsKeyname() string {
	return "withdrawals"
}

func retrieveWithdrawalsByTokenIDKeyname(tokenID *uuid.UUID) string {
	base := retrieveAllWithdrawalsKeyname()
	return fmt.Sprintf("%s:by_token_id:%s", base, tokenID.String())
}

func retrieveWithdrawalsByFromWalletIDKeyname(toWalletID *uuid.UUID) string {
	base := retrieveAllWithdrawalsKeyname()
	return fmt.Sprintf("%s:by_from_wallet_id:%s", base, toWalletID.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Withdrawal",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableWithdrawal) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				fromWalletID, fromWalletIDErr := uuid.FromString(storable.FromWalletID)
				if fromWalletIDErr != nil {
					return nil, fromWalletIDErr
				}

				tokenID, tokenIDErr := uuid.FromString(storable.TokenID)
				if tokenIDErr != nil {
					return nil, tokenIDErr
				}

				// retrieve the wallet:
				walletMetaData := wallet.SDKFunc.CreateMetaData()
				walletIns, walletInsErr := rep.RetrieveByID(walletMetaData, &fromWalletID)
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
						out := createWithdrawal(&id, wal, tok, storable.Amount)
						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid Token instance", tokenID.String())
					return nil, errors.New(str)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", fromWalletID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableWithdrawal); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedWithdrawal)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createWithdrawalFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if with, ok := ins.(Withdrawal); ok {
				return createNormalizedWithdrawal(with)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Withdrawal instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedWithdrawal); ok {
				return createWithdrawalFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Withdrawal instance")
		},
		EmptyStorable:   new(storableWithdrawal),
		EmptyNormalized: new(storableWithdrawal),
	})
}
