package deposit

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllDepositsKeyname() string {
	return "cryptodeposits"
}

func retrieveDepositByOfferKeyname(off offer.Offer) string {
	base := retrieveAllDepositsKeyname()
	return fmt.Sprintf("%s:by_offer_id:%s", base, off.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "CryptoDeposit",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableDeposit); ok {
				return createDepositFromStorable(storable, rep)
			}

			ptr := new(normalizedDeposit)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createDepositFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if dep, ok := ins.(Deposit); ok {
				out, outErr := createNormalizedDeposit(dep)
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

			return nil, errors.New("the given instance is not a valid normalized Deposit instance")
		},
		EmptyNormalized: new(normalizedDeposit),
		EmptyStorable:   new(storableDeposit),
	})
}

func createRepresentation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if dep, ok := ins.(Deposit); ok {
				out := createStorableDeposit(dep)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if dep, ok := ins.(Deposit); ok {
				return []string{
					retrieveAllDepositsKeyname(),
					retrieveDepositByOfferKeyname(dep.Offer()),
				}, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
			return nil, errors.New(str)
		},
	})
}

func toData(dep Deposit) *Data {
	out := Data{
		ID:     dep.ID().String(),
		Offer:  offer.SDKFunc.ToData(dep.Offer()),
		From:   address.SDKFunc.ToData(dep.From()),
		Amount: dep.Amount(),
	}

	return &out
}

func toDataSet(ps entity.PartialSet) (*DataSet, error) {
	ins := ps.Instances()
	deposits := []*Data{}
	for _, oneIns := range ins {
		if dep, ok := oneIns.(Deposit); ok {
			deposits = append(deposits, toData(dep))
			continue
		}

		return nil, errors.New("there is at least one entity that is not a valid Deposit instance")
	}

	out := DataSet{
		Index:       ps.Index(),
		Amount:      ps.Amount(),
		TotalAmount: ps.TotalAmount(),
		IsLast:      ps.IsLast(),
		Deposits:    deposits,
	}

	return &out, nil
}
