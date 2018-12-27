package currency

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

const (
	amountOfCharactersForSymbol         = 6
	maxAountOfCharactersForName         = 50
	maxAmountOfCharactersForDescription = 500
)

// Currency represents a currency
type Currency interface {
	ID() *uuid.UUID
	Category() category.Category
	Symbol() string
	Name() string
	Description() string
}

// Repository represents the currency repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Currency, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// Normalized represents a normalized currency
type Normalized interface {
}

// Data represents the currency data
type Data struct {
	ID          string
	Category    *category.Data
	Symbol      string
	Name        string
	Description string
}

// DataSet represents the human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Currencies  []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID          *uuid.UUID
	Category    category.Category
	Symbol      string
	Name        string
	Description string
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Currency SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Currency
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(curr Currency) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
}{
	Create: func(params CreateParams) Currency {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createCurrency(params.ID, params.Category, params.Symbol, params.Name, params.Description)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if curr, ok := ins.(Currency); ok {
					out := createStorableCurrency(curr)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Currency instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if _, ok := ins.(Currency); ok {
					return []string{
						retrieveAllCurrenciesKeyname(),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Currency instance", ins.ID().String())
				return nil, errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(metaData, params.EntityRepository)
		return out
	},
	ToData: func(curr Currency) *Data {
		return toData(curr)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
