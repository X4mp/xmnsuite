package keyname

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
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

// Data represents human-redable data
type Data struct {
	ID    string
	Group *group.Data
	Name  string
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Keynames    []*Data
}

// DataSetWithGroup represents human-redable data set with group
type DataSetWithGroup struct {
	Group    *group.Data
	Keynames *DataSet
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

// RouteSetOfGroupParams represents the route set of group params
type RouteSetOfGroupParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Keyname
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(kname Keyname) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteSetOfGroup      func(params RouteSetOfGroupParams) func(w http.ResponseWriter, r *http.Request)
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
	ToData: func(kname Keyname) *Data {
		return toData(kname)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	RouteSetOfGroup: func(params RouteSetOfGroupParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// retrieve the group name:
			vars := mux.Vars(r)
			if groupName, ok := vars["name"]; ok {
				// create metadata:
				metaData := createMetaData()

				// create repositories:
				groupRepository := group.SDKFunc.CreateRepository(group.CreateRepositoryParams{
					EntityRepository: params.EntityRepository,
				})

				keynameRepository := createRepository(params.EntityRepository, metaData)

				// retrieve the group:
				grp, grpErr := groupRepository.RetrieveByName(groupName)
				if grpErr != nil {
					w.WriteHeader(http.StatusNotFound)
					str := fmt.Sprintf("the given group (name: %s) could not be found: %s", groupName, grpErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the keynames by group:
				knamePS, knamePSErr := keynameRepository.RetrieveSetByGroup(grp, 0, params.AmountOfElementsPerList)
				if knamePSErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retireving request keynames by grpoup (ID: %s): %s", grp.ID().String(), knamePSErr.Error())
					w.Write([]byte(str))
					return
				}

				// render:
				datSet, datSetErr := toDataSet(knamePS)
				if datSetErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while converting the keyname entity set to data: %s", datSetErr.Error())
					w.Write([]byte(str))
					return
				}

				w.WriteHeader(http.StatusOK)
				params.Tmpl.Execute(w, &DataSetWithGroup{
					Group:    group.SDKFunc.ToData(grp),
					Keynames: datSet,
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the group name is mandatory")
			w.Write([]byte(str))
		}
	},
}
