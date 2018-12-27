package token

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllTokensKeyname() string {
	return "tokens"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Token",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableToken); ok {
				return createTokenFromStorable(storable)
			}

			ptr := new(storableToken)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createTokenFromStorable(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if tok, ok := ins.(Token); ok {
				out := createStorableToken(tok)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Token instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if storable, ok := ins.(*storableToken); ok {
				return createTokenFromStorable(storable)
			}

			return nil, errors.New("the given instance is not a valid normalized Token instance")
		},
		EmptyStorable:   new(storableToken),
		EmptyNormalized: new(storableToken),
	})
}

func toData(tok Token) *Data {
	out := Data{
		ID:          tok.ID().String(),
		Symbol:      tok.Symbol(),
		Name:        tok.Name(),
		Description: tok.Description(),
	}

	return &out
}

func toDataSet(ins entity.PartialSet) (*DataSet, error) {
	data := []*Data{}
	instances := ins.Instances()
	for _, oneIns := range instances {
		if tok, ok := oneIns.(Token); ok {
			data = append(data, toData(tok))
			continue
		}

		str := fmt.Sprintf("at least one of the elements (ID: %s) in the entity partial set is not a valid Token instance", oneIns.ID().String())
		return nil, errors.New(str)
	}

	out := DataSet{
		Index:       ins.Index(),
		Amount:      ins.Amount(),
		TotalAmount: ins.TotalAmount(),
		IsLast:      ins.IsLast(),
		Tokens:      data,
	}

	return &out, nil
}
