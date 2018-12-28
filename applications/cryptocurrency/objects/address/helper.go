package address

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
)

func retrieveAllAddressKeyname() string {
	return "addresses"
}

func retrieveAddressByWalletKeyname(wal wallet.Wallet) string {
	base := retrieveAllAddressKeyname()
	return fmt.Sprintf("%s:by_wallet_id:%s", base, wal.ID().String())
}

func retrieveAddressByAddressKeyname(addr []byte) string {
	base := retrieveAllAddressKeyname()
	return fmt.Sprintf("%s:by_address:%s", base, hex.EncodeToString(addr))
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Address",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableAddress); ok {
				return createAddressFromStorable(storable, rep)
			}

			ptr := new(normalizedAddress)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createAddressFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if addr, ok := ins.(Address); ok {
				out, outErr := createNormalizedAddress(addr)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Address instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedAddress); ok {
				return createAddressFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Address instance")
		},
		EmptyNormalized: new(normalizedAddress),
		EmptyStorable:   new(storableAddress),
	})
}

func createRepresentation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if addr, ok := ins.(Address); ok {
				out := createStorableAddress(addr)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Address instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if addr, ok := ins.(Address); ok {
				return []string{
					retrieveAllAddressKeyname(),
					retrieveAddressByWalletKeyname(addr.Wallet()),
					retrieveAddressByAddressKeyname(addr.Address()),
				}, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Address instance", ins.ID().String())
			return nil, errors.New(str)
		},
	})
}

func toData(addr Address) *Data {
	out := Data{
		ID:      addr.ID().String(),
		Wallet:  wallet.SDKFunc.ToData(addr.Wallet()),
		Address: hex.EncodeToString(addr.Address()),
	}

	return &out
}

func toDataSet(ps entity.PartialSet) (*DataSet, error) {
	ins := ps.Instances()
	addresses := []*Data{}
	for _, oneIns := range ins {
		if addr, ok := oneIns.(Address); ok {
			addresses = append(addresses, toData(addr))
			continue
		}

		return nil, errors.New("there is at least one entity that is not a valid Address instance")
	}

	out := DataSet{
		Index:       ps.Index(),
		Amount:      ps.Amount(),
		TotalAmount: ps.TotalAmount(),
		IsLast:      ps.IsLast(),
		Addresses:   addresses,
	}

	return &out, nil
}
