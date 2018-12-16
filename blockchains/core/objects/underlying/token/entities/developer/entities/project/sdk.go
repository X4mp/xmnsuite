package project

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Project represents a project
type Project interface {
	ID() *uuid.UUID
	Title() string
	Description() string
}

// Normalized represents a normalized project
type Normalized interface {
}

// CreateParams represents the create params
type CreateParams struct {
	ID          *uuid.UUID
	Title       string
	Description string
}

// SDKFunc represents the Project SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Project
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Project {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createProject(params.ID, params.Title, params.Description)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if proj, ok := ins.(Project); ok {
					out := createStorableProject(proj)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Project instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllProjectsKeyname(),
				}, nil
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)

				// create the metadata:
				metaData := createMetaData()

				if proj, ok := ins.(Project); ok {
					// if the project already exists, return an error:
					_, retProjErr := repository.RetrieveByID(metaData, proj.ID())
					if retProjErr == nil {
						str := fmt.Sprintf("the Project (ID: %s) already exists", proj.ID().String())
						return errors.New(str)
					}

					// everything is alright:
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Project instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
