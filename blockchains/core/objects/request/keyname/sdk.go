package keyname

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Keyname represents a keyname a request can made on
type Keyname interface {
	ID() *uuid.UUID
	Group() group.Group
	Name() string
}

// Normalized represents a normalized keyname
type Normalized interface {
}

// Repository represents a keyname repository
type Repository interface {
	RetrieveByName(name string) (Keyname, error)
	RetrieveSetByGroup(grp group.Group, index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the create params
type CreateParams struct {
	ID    *uuid.UUID
	Group group.Group
	Name  string
}

// CreateRepositoryParams represents the create repository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Keyname
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Keyname {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createKeyname(params.ID, params.Group, params.Name)
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
				if kname, ok := ins.(Keyname); ok {
					out := createStorableKeyname(kname)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Keyname instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if kname, ok := ins.(Keyname); ok {
					return []string{
						retrieveAllKeynamesKeyname(),
						retrieveKeynameByNameKeyname(kname.Name()),
						retrieveKeynameByGroupKeyname(kname.Group()),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid Keyname instance")
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				if kname, ok := ins.(Keyname); ok {
					// metadata:
					metaData := createMetaData()
					groupRepresentation := group.SDKFunc.CreateRepresentation()

					// create the repository and service:
					repository := entity.SDKFunc.CreateRepository(ds)
					service := entity.SDKFunc.CreateService(ds)
					kanameRepository := createRepository(repository, metaData)
					groupRepository := group.SDKFunc.CreateRepository(group.CreateRepositoryParams{
						EntityRepository: repository,
					})

					// the keyname must not exists:
					_, retKnameErr := repository.RetrieveByID(metaData, kname.ID())
					if retKnameErr == nil {
						str := fmt.Sprintf("the Keyname (ID: %s) already exists", kname.ID().String())
						return errors.New(str)
					}

					// the name must be unique:
					_, retKnameByNameErr := kanameRepository.RetrieveByName(kname.Name())
					if retKnameByNameErr == nil {
						str := fmt.Sprintf("there is already a Keyname instance under that name: %s", kname.Name())
						return errors.New(str)
					}

					// if the group does not exists, create it:
					_, retGrpErr := groupRepository.RetrieveByName(kname.Group().Name())
					if retGrpErr != nil {
						saveErr := service.Save(kname.Group(), groupRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Keyname instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(params.EntityRepository, metaData)
		return out
	},
}
