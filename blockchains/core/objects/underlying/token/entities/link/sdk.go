package link

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Link represents a blockchain link
type Link interface {
	ID() *uuid.UUID
	Title() string
	Description() string
}

// Normalized represents the normalized link
type Normalized interface {
}

// Repository represents a repository
type Repository interface {
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents a Create params
type CreateParams struct {
	ID          *uuid.UUID
	Title       string
	Description string
}

// SDKFunc represents the Link SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Link
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Link {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createLink(params.ID, params.Title, params.Description)

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
				if lnk, ok := ins.(Link); ok {
					out := createStorableLink(lnk)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Link instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllLinksKeyname(),
				}, nil
			},
			OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)

				// create the metadata and representations:
				metaData := createMetaData()

				if lnk, ok := ins.(Link); ok {
					// if the link already exists:
					_, retLinkInsErr := repository.RetrieveByID(metaData, lnk.ID())
					if retLinkInsErr == nil {
						str := fmt.Sprintf("the Link (ID: %s) already exists", lnk.ID().String())
						return errors.New(str)
					}

					// the link doesnt exists, so everything is fine:
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Link instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
