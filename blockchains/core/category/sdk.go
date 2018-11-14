package category

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

// Category represents a category
type Category interface {
	ID() *uuid.UUID
	Name() string
	Description() string
}

// Normalized represents a normalized category
type Normalized interface {
}

// SDKFunc represents the Category SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if cat, ok := ins.(Category); ok {
					out := createStorableCategory(cat)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Category instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllCategoriesKeyname(),
				}, nil
			},
		})
	},
}
