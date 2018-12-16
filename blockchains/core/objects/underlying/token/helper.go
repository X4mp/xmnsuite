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
		Normalize: toData,
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

func toData(ins entity.Entity) (interface{}, error) {
	if tok, ok := ins.(Token); ok {
		out := createStorableToken(tok)
		return out, nil
	}

	str := fmt.Sprintf("the given entity (ID: %s) is not a valid Token instance", ins.ID().String())
	return nil, errors.New(str)
}
