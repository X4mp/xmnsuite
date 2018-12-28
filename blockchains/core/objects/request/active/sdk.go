package active

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	core_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Request represents an active request
type Request interface {
	ID() *uuid.UUID
	Request() core_request.Request
	ConcensusNeeded() int
}

// Normalized represents a normalized request
type Normalized interface {
}

// Repository represents a Request repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Request, error)
	RetrieveByRequest(req core_request.Request) (Request, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByFromUser(usr user.User, index int, amount int) (entity.PartialSet, error)
	RetrieveSetByKeyname(kname keyname.Keyname, index int, amount int) (entity.PartialSet, error)
}

// Data represents human-redable data
type Data struct {
	ID              string
	Request         *core_request.Data
	ConcensusNeeded int
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Requests    []*Data
}

// DataSetWithKeyname represents human-redable data set with keyname
type DataSetWithKeyname struct {
	Keyname  *keyname.Data
	Requests *DataSet
}

// CreateParams represents the create params
type CreateParams struct {
	ID              *uuid.UUID
	Request         core_request.Request
	ConcensusNeeded int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// RouteSetOfKeynameParams represents the route set params
type RouteSetOfKeynameParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// SDKFunc represents the request SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Request
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(req Request) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteSetOfKeyname    func(params RouteSetOfKeynameParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Request {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		// create the request:
		out, outErr := createRequest(params.ID, params.Request, params.ConcensusNeeded)
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
				if req, ok := ins.(Request); ok {
					out := createStorable(req)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid active Request instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if req, ok := ins.(Request); ok {
					return []string{
						retrieveAllRequestsKeyname(),
						retrieveAllRequestsByRequestKeyname(req.Request()),
						retrieveAllRequestsFromUserKeyname(req.Request().From()),
						retrieveAllRequestsByKeynameKeyname(req.Request().Keyname()),
					}, nil
				}

				return nil, errors.New("the given entity is not a valid active Request instance")
			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				if req, ok := ins.(Request); ok {
					// metadata:
					metaData := createMetaData()
					coreRequestRepresentation := core_request.SDKFunc.CreateRepresentation()

					// create the repository and service:
					entityRepository := entity.SDKFunc.CreateRepository(ds)
					service := entity.SDKFunc.CreateService(ds)

					// make sure the request does not exists:
					_, retReqErr := entityRepository.RetrieveByID(metaData, req.ID())
					if retReqErr == nil {
						str := fmt.Sprintf("the Request (ID: %s) already exists", req.ID().String())
						return errors.New(str)
					}

					// make sure the request does not exits, then save it:
					_, retPrevReqErr := entityRepository.RetrieveByID(coreRequestRepresentation.MetaData(), req.Request().ID())
					if retPrevReqErr == nil {
						str := fmt.Sprintf("the given Request (ID: %s) already exists: %s", req.Request().ID().String(), retPrevReqErr.Error())
						return errors.New(str)
					}

					// save the request:
					saveReqErr := service.Save(req.Request(), coreRequestRepresentation)
					if saveReqErr != nil {
						return saveReqErr
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid active Request instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		return createRepository(params.EntityRepository, metaData)
	},
	ToData: func(req Request) *Data {
		return toData(req)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	RouteSetOfKeyname: func(params RouteSetOfKeynameParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// retrieve the group name:
			vars := mux.Vars(r)
			if groupName, ok := vars["groupname"]; ok {
				// retrieve the keyname:
				if keynameAsString, ok := vars["keyname"]; ok {
					// create the metadata:
					metaData := createMetaData()

					// create the repositories:
					requestRepository := createRepository(params.EntityRepository, metaData)
					groupRepository := group.SDKFunc.CreateRepository(group.CreateRepositoryParams{
						EntityRepository: params.EntityRepository,
					})

					keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
						EntityRepository: params.EntityRepository,
					})

					// retrieve the group:
					grp, grpErr := groupRepository.RetrieveByName(groupName)
					if grpErr != nil {
						w.WriteHeader(http.StatusNotFound)
						str := fmt.Sprintf("the given group (name: %s) could not be found: %s", groupName, grpErr.Error())
						w.Write([]byte(str))
						return
					}

					// retrieve the keyname:
					kname, knameErr := keynameRepository.RetrieveByName(keynameAsString)
					if knameErr != nil {
						w.WriteHeader(http.StatusNotFound)
						str := fmt.Sprintf("the given keyname (name: %s) could not be found: %s", keynameAsString, knameErr.Error())
						w.Write([]byte(str))
						return
					}

					// make sure the group in the keyname fits the group:
					if bytes.Compare(grp.ID().Bytes(), kname.Group().ID().Bytes()) != 0 {
						w.WriteHeader(http.StatusNotFound)
						str := fmt.Sprintf("the given group (%s) does not fit the given keyname (%s) group (%s)", grp.Name(), kname.Name(), kname.Group().Name())
						w.Write([]byte(str))
						return
					}

					// retrieve the request related to our keyname:
					reqPS, reqPSErr := requestRepository.RetrieveSetByKeyname(kname, 0, params.AmountOfElementsPerList)
					if reqPSErr != nil {
						w.WriteHeader(http.StatusNotFound)
						str := fmt.Sprintf("there was an error while retrieving the requests related to the given keyname (ID: %s): %s", kname.ID().String(), reqPSErr.Error())
						w.Write([]byte(str))
						return
					}

					// render:
					datSet, datSetErr := toDataSet(reqPS)
					if datSetErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while converting the request entity set to data: %s", datSetErr.Error())
						w.Write([]byte(str))
						return
					}

					w.WriteHeader(http.StatusOK)
					params.Tmpl.Execute(w, &DataSetWithKeyname{
						Keyname:  keyname.SDKFunc.ToData(kname),
						Requests: datSet,
					})
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("the keyname is mandatory")
				w.Write([]byte(str))
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the group name is mandatory")
			w.Write([]byte(str))
		}
	},
}
