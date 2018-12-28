package offer

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

func retrieveAllOffersKeyname() string {
	return "offers"
}

func retrieveOfferByPledge(pldge pledge.Pledge) string {
	base := retrieveAllOffersKeyname()
	return fmt.Sprintf("%s:by_pledge_id:%s", base, pldge.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Offer",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableOffer); ok {
				return createOfferFromStorable(storable, rep)
			}

			ptr := new(normalizedOffer)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createOfferFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if off, ok := ins.(Offer); ok {
				out, outErr := createNormalizedOffer(off)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Offer instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedOffer); ok {
				return createOfferFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Offer instance")
		},
		EmptyNormalized: new(normalizedOffer),
		EmptyStorable:   new(storableOffer),
	})
}

func createRepresentation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if off, ok := ins.(Offer); ok {
				out := createStorableOffer(off)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Offer instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if off, ok := ins.(Offer); ok {
				return []string{
					retrieveAllOffersKeyname(),
					retrieveOfferByPledge(off.Pledge()),
				}, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Offer instance", ins.ID().String())
			return nil, errors.New(str)
		},
	})
}

func toData(off Offer) *Data {
	out := Data{
		ID:        off.ID().String(),
		Pledge:    pledge.SDKFunc.ToData(off.Pledge()),
		ToAddress: address.SDKFunc.ToData(off.To()),
		Amount:    off.Amount(),
		Price:     off.Price(),
	}

	return &out
}

func toDataSet(ps entity.PartialSet) (*DataSet, error) {
	ins := ps.Instances()
	offers := []*Data{}
	for _, oneIns := range ins {
		if off, ok := oneIns.(Offer); ok {
			offers = append(offers, toData(off))
			continue
		}

		return nil, errors.New("there is at least one entity that is not a valid Offer instance")
	}

	out := DataSet{
		Index:       ps.Index(),
		Amount:      ps.Amount(),
		TotalAmount: ps.TotalAmount(),
		IsLast:      ps.IsLast(),
		Offers:      offers,
	}

	return &out, nil
}
