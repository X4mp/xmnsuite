package group

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Group represents a group a request can made on
type Group interface {
	ID() *uuid.UUID
	Name() string
}

// Normalized represents a normalized group
type Normalized interface {
}

// Repository represents a group repository
type Repository interface {
	RetrieveByName(name string) (Group, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the create params
type CreateParams struct {
	ID   *uuid.UUID
	Name string
}

// CreateRepositoryParams represents the create repository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// Data represents human-redable data
type Data struct {
	ID   string
	Name string
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Groups      []*Data
}

// RouteSetParams represents the route set params
type RouteSetParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Group
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(grp Group) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteSet             func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Group {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createGroup(params.ID, params.Name)
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
				if grp, ok := ins.(Group); ok {
					out := createStorableGroup(grp)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Group instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if grp, ok := ins.(Group); ok {
					return []string{
						retrieveAllGroupsKeyname(),
						retrieveGroupByNameKeyname(grp.Name()),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid Group instance")
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				if grp, ok := ins.(Group); ok {
					// metadata:
					metaData := createMetaData()

					// create the repository and service:
					repository := entity.SDKFunc.CreateRepository(ds)
					kanameRepository := createRepository(repository, metaData)

					// the group must not exists:
					_, retGrpErr := repository.RetrieveByID(metaData, grp.ID())
					if retGrpErr == nil {
						str := fmt.Sprintf("the Group (ID: %s) already exists", grp.ID().String())
						return errors.New(str)
					}

					// the name must be unique:
					_, retGrpByNameErr := kanameRepository.RetrieveByName(grp.Name())
					if retGrpByNameErr == nil {
						str := fmt.Sprintf("there is already a Group instance under that name: %s", grp.Name())
						return errors.New(str)
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Group instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(params.EntityRepository, metaData)
		return out
	},
	ToData: func(grp Group) *Data {
		return toData(grp)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	RouteSet: func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create metadata:
			metaData := createMetaData()

			// create the repositories:
			groupRepository := createRepository(params.EntityRepository, metaData)

			// retrieve the groups:
			retGroupsPS, retGroupsPSErr := groupRepository.RetrieveSet(0, params.AmountOfElementsPerList)
			if retGroupsPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retireving request groups: %s", retGroupsPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// render:
			datSet, datSetErr := toDataSet(retGroupsPS)
			if datSetErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while converting the group entity set to data: %s", datSetErr.Error())
				w.Write([]byte(str))
				return
			}

			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, datSet)
		}
	},
}
