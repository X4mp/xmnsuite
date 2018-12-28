package chain

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllChainsKeyname() string {
	return "chains"
}

func retrieveChainByOfferKeyname(off offer.Offer) string {
	base := retrieveAllChainsKeyname()
	return fmt.Sprintf("%s:by_offer_id:%s", base, off.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Chain",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableChain); ok {
				return createChainFromStorable(storable, rep)
			}

			ptr := new(normalizedChain)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createChainFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if chn, ok := ins.(Chain); ok {
				out, outErr := createNormalizedChain(chn)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Chain instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedChain); ok {
				return createChainFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Chain instance")
		},
		EmptyNormalized: new(normalizedChain),
		EmptyStorable:   new(storableChain),
	})
}

func createRepresentation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if chn, ok := ins.(Chain); ok {
				out := createStorableChain(chn)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Chain instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if chn, ok := ins.(Chain); ok {
				return []string{
					retrieveAllChainsKeyname(),
					retrieveChainByOfferKeyname(chn.Offer()),
				}, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Chain instance", ins.ID().String())
			return nil, errors.New(str)
		},
	})
}

func toData(chn Chain) *Data {
	out := Data{
		ID:          chn.ID().String(),
		Offer:       offer.SDKFunc.ToData(chn.Offer()),
		TotalAmount: chn.TotalAmount(),
	}

	return &out
}

func toDataSet(ps entity.PartialSet) (*DataSet, error) {
	ins := ps.Instances()
	chains := []*Data{}
	for _, oneIns := range ins {
		if chn, ok := oneIns.(Chain); ok {
			chains = append(chains, toData(chn))
			continue
		}

		return nil, errors.New("there is at least one entity that is not a valid Chain instance")
	}

	out := DataSet{
		Index:       ps.Index(),
		Amount:      ps.Amount(),
		TotalAmount: ps.TotalAmount(),
		IsLast:      ps.IsLast(),
		Chains:      chains,
	}

	return &out, nil
}
