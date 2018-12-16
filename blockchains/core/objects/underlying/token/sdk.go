package token

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
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

// CreateParams represents the Create params
type CreateParams struct {
	ID          *uuid.UUID
	Symbol      string
	Name        string
	Description string
}

// SDKFunc represents the Token SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Token
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Token {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createToken(params.ID, params.Symbol, params.Name, params.Description)
		return out
	},
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
