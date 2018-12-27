package token

import (
	"errors"
	"fmt"

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

// Data represents human-redable data
type Data struct {
	ID          string
	Symbol      string
	Name        string
	Description string
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Tokens      []*Data
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
	ToData               func(tok Token) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
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
	ToData: func(tok Token) *Data {
		return toData(tok)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
