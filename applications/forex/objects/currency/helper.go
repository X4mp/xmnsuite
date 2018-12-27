package currency

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllCurrenciesKeyname() string {
	return "currencies"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Currency",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableCurrency); ok {
				return fromStorableToCurrency(storable, rep)
			}

			ptr := new(normalizedCurrency)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromNormalizedToCurrency(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if curr, ok := ins.(Currency); ok {
				return createNormalizedCurrency(curr)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Currency instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedCurrency); ok {
				return fromNormalizedToCurrency(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Currency instance")
		},
		EmptyNormalized: new(normalizedCurrency),
		EmptyStorable:   new(storableCurrency),
	})
}

func toData(curr Currency) *Data {
	out := Data{
		ID:          curr.ID().String(),
		Category:    category.SDKFunc.ToData(curr.Category()),
		Symbol:      curr.Symbol(),
		Name:        curr.Name(),
		Description: curr.Description(),
	}

	return &out
}

func toDataSet(ps entity.PartialSet) (*DataSet, error) {
	ins := ps.Instances()
	currencies := []*Data{}
	for _, oneIns := range ins {
		if curr, ok := oneIns.(Currency); ok {
			currencies = append(currencies, toData(curr))
			continue
		}

		return nil, errors.New("there is at least one entity that is not a valid Currency instance")
	}

	out := DataSet{
		Index:       ps.Index(),
		Amount:      ps.Amount(),
		TotalAmount: ps.TotalAmount(),
		IsLast:      ps.IsLast(),
		Currencies:  currencies,
	}

	return &out, nil
}
