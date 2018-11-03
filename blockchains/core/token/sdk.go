package token

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
)

// Token represents the token
type Token interface {
	ID() *uuid.UUID
	Symbol() string
	Name() string
	Description() string
}

// SDKFunc represents the Token SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if tok, ok := ins.(Token); ok {
					out := createStorableToken(tok)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Token instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllTokensKeyname(),
				}, nil
			},
		})
	},
}
