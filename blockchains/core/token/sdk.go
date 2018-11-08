package token

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

// Token represents the token
type Token interface {
	ID() *uuid.UUID
	Symbol() string
	Name() string
	Description() string
}

// Normalized represents the normalized Token
type Normalized interface {
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
			Met:        createMetaData(),
			ToStorable: toData,
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllTokensKeyname(),
				}, nil
			},
		})
	},
}
