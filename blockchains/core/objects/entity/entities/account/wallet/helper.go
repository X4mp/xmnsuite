package wallet

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/crypto"
)

func retrieveAllWalletKeyname() string {
	return "wallets"
}

func retrieveByPublicKeyWalletKeyname(pubKey crypto.PublicKey) string {
	base := retrieveAllWalletKeyname()
	return fmt.Sprintf("%s:by_pubkey:%s", base, pubKey.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Wallet",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableWallet); ok {
				return createWalletFromStorable(storable)
			}

			ptr := new(storableWallet)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createWalletFromStorable(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if wallet, ok := ins.(Wallet); ok {
				out := createStoredWallet(wallet)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if storable, ok := ins.(*storableWallet); ok {
				return createWalletFromStorable(storable)
			}

			return nil, errors.New("the given instance is not a valid normalized Wallet instance")
		},
		EmptyStorable:   new(storableWallet),
		EmptyNormalized: new(storableWallet),
	})
}

func toData(wal Wallet) *Data {
	out := Data{
		ID:              wal.ID().String(),
		Creator:         wal.Creator().String(),
		ConcensusNeeded: wal.ConcensusNeeded(),
	}

	return &out
}

func toDataSet(ins entity.PartialSet) (*DataSet, error) {
	data := []*Data{}
	instances := ins.Instances()
	for _, oneIns := range instances {
		if wal, ok := oneIns.(Wallet); ok {
			data = append(data, toData(wal))
			continue
		}

		str := fmt.Sprintf("at least one of the elements (ID: %s) in the entity partial set is not a valid Wallet instance", oneIns.ID().String())
		return nil, errors.New(str)
	}

	out := DataSet{
		Index:       ins.Index(),
		Amount:      ins.Amount(),
		TotalAmount: ins.TotalAmount(),
		IsLast:      ins.IsLast(),
		Wallets:     data,
	}

	return &out, nil
}
