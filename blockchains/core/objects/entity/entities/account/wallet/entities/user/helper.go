package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

func retrieveAllUserKeyname() string {
	return "users"
}

func retrieveUserByPubKeyKeyname(pubKey crypto.PublicKey) string {
	base := retrieveAllUserKeyname()
	return fmt.Sprintf("%s:by_public_key:%s", base, pubKey.String())
}

func retrieveUserByWalletIDKeyname(walletID *uuid.UUID) string {
	base := retrieveAllUserKeyname()
	return fmt.Sprintf("%s:by_wallet_id:%s", base, walletID.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "User",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableUser) (entity.Entity, error) {
				// create the metadata:
				walletMetaData := wallet.SDKFunc.CreateMetaData()

				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
					PubKeyAsString: storable.PubKey,
				})

				walletID, walletIDErr := uuid.FromString(storable.WalletID)
				if walletIDErr != nil {
					return nil, walletIDErr
				}

				ins, insErr := rep.RetrieveByID(walletMetaData, &walletID)
				if insErr != nil {
					return nil, insErr
				}

				if wal, ok := ins.(wallet.Wallet); ok {
					out := createUser(&id, pubKey, storable.Shares, wal)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance and thererfore the given data cannot be transformed to a User instance", walletID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableUser); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedUser)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createUserFromNormalizedUser(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if usr, ok := ins.(User); ok {
				return createNormalizedUser(usr)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedUser); ok {
				return createUserFromNormalizedUser(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized User instance")
		},
		EmptyStorable:   new(storableUser),
		EmptyNormalized: new(normalizedUser),
	})
}

func toData(usr User) *Data {
	out := Data{
		ID:     usr.ID().String(),
		PubKey: usr.PubKey().String(),
		Shares: usr.Shares(),
		Wallet: wallet.SDKFunc.ToData(usr.Wallet()),
	}

	return &out
}

func toDataSet(ins entity.PartialSet) (*DataSet, error) {
	data := []*Data{}
	instances := ins.Instances()
	for _, oneIns := range instances {
		if usr, ok := oneIns.(User); ok {
			data = append(data, toData(usr))
			continue
		}

		str := fmt.Sprintf("at least one of the elements (ID: %s) in the entity partial set is not a valid User instance", oneIns.ID().String())
		return nil, errors.New(str)
	}

	out := DataSet{
		Index:       ins.Index(),
		Amount:      ins.Amount(),
		TotalAmount: ins.TotalAmount(),
		IsLast:      ins.IsLast(),
		Users:       data,
	}

	return &out, nil
}
