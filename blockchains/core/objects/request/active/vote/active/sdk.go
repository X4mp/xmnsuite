package active

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	core_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Vote represents an active request vote
type Vote interface {
	ID() *uuid.UUID
	Vote() core_vote.Vote
	Power() int
}

// Normalized represents a normalized vote
type Normalized interface {
}

// Repository represents a vote repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Vote, error)
	RetrieveByVote(vot core_vote.Vote) (Vote, error)
	RetrieveByRequestVoter(voter user.User, req active_request.Request) (Vote, error)
	RetrieveSetByRequest(req active_request.Request, index int, amount int) (entity.PartialSet, error)
}

// Service represents the vote service
type Service interface {
	Save(ins Vote, rep entity.Representation) error
}

// Data represents human-redable data
type Data struct {
	ID    string
	Vote  *core_vote.Data
	Power int
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Votes       []*Data
}

// DataSetOfRequest represents human-redable data set of request
type DataSetOfRequest struct {
	Request        *active_request.Data
	Keyname        *keyname.Data
	MyUsers        *user.DataSet
	Votes          *DataSet
	ApprovedPow    int
	DisapprovedPow int
	NeutralPow     int
	TotalPow       int
}

// CreateParams represents the create params
type CreateParams struct {
	ID    *uuid.UUID
	Vote  core_vote.Vote
	Power int
}

// CreateServiceParams represents the CreateService params
type CreateServiceParams struct {
	DS               datastore.DataStore
	EntityRepository entity.Repository
	EntityService    entity.Service
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// RouteSetOfRequestParams represents the route set of request params
type RouteSetOfRequestParams struct {
	AmountOfElementsPerList int
	PK                      crypto.PrivateKey
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Vote
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	CreateService        func(params CreateServiceParams) Service
	ToData               func(vot Vote) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteSetOfRequest    func(params RouteSetOfRequestParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Vote {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		// create the request:
		out, outErr := createVote(params.ID, params.Vote, params.Power)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return createRepresentation()
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(params.EntityRepository, metaData)
		return out
	},
	CreateService: func(params CreateServiceParams) Service {
		if params.EntityService == nil && params.EntityRepository == nil {
			params.EntityRepository = entity.SDKFunc.CreateRepository(params.DS)
			params.EntityService = entity.SDKFunc.CreateService(params.DS)
		}

		voteRepresentation := createRepresentation()
		requestRepresentation := active_request.SDKFunc.CreateRepresentation()
		out := createVoteService(params.EntityRepository, params.EntityService, voteRepresentation, requestRepresentation)
		return out
	},
	ToData: func(vot Vote) *Data {
		return toData(vot)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	RouteSetOfRequest: func(params RouteSetOfRequestParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// retrieve the group name:
			vars := mux.Vars(r)
			if groupName, ok := vars["groupname"]; ok {
				// retrieve the keyname:
				if keynameAsString, ok := vars["keyname"]; ok {
					if requestIDAsString, ok := vars["id"]; ok {
						// create the representation:
						metaData := createMetaData()

						// create the repositories:
						voteRepository := createRepository(params.EntityRepository, metaData)
						requestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
							EntityRepository: params.EntityRepository,
						})

						groupRepository := group.SDKFunc.CreateRepository(group.CreateRepositoryParams{
							EntityRepository: params.EntityRepository,
						})

						keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
							EntityRepository: params.EntityRepository,
						})

						userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
							EntityRepository: params.EntityRepository,
						})

						// parse the id:
						id, idErr := uuid.FromString(requestIDAsString)
						if idErr != nil {
							w.WriteHeader(http.StatusNotFound)
							str := fmt.Sprintf("the id (%s) is invalid: %s", requestIDAsString, idErr.Error())
							w.Write([]byte(str))
							return
						}

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

						// retrieve the request:
						req, reqErr := requestRepository.RetrieveByID(&id)
						if reqErr != nil {
							w.WriteHeader(http.StatusNotFound)
							str := fmt.Sprintf("the given request (ID: %s) could not be found: %s", id.String(), reqErr.Error())
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

						// make sure the keyname and the keyname in the request fits:
						if bytes.Compare(kname.ID().Bytes(), req.Request().Keyname().ID().Bytes()) != 0 {
							w.WriteHeader(http.StatusNotFound)
							str := fmt.Sprintf("the given keyname (%s) does not fit the given request (ID: %s) keyname (%s)", kname.Name(), req.ID().String(), req.Request().Keyname().ID().String())
							w.Write([]byte(str))
							return
						}

						// retrieve
						myUsersPS, myUsersPSErr := userRepository.RetrieveSetByPubKey(params.PK.PublicKey(), 0, params.AmountOfElementsPerList)
						if myUsersPSErr != nil {
							w.WriteHeader(http.StatusInternalServerError)
							str := fmt.Sprintf("there was an error while retrieving the user set (PubKey: %s): %s", params.PK.PublicKey().String(), myUsersPSErr.Error())
							w.Write([]byte(str))
							return
						}

						// retrieve the votes associated with the request:
						votesPS, votesPSErr := voteRepository.RetrieveSetByRequest(req, 0, params.AmountOfElementsPerList)
						if votesPSErr != nil {
							w.WriteHeader(http.StatusInternalServerError)
							str := fmt.Sprintf("there was an  error while retrieving votes related to the given Request (ID: %s): %s", req.ID().String(), votesPSErr.Error())
							w.Write([]byte(str))
							return
						}

						approvedPow := 0
						disApprovedPow := 0
						neutralPow := 0
						votsIns := votesPS.Instances()
						for _, oneVoteIns := range votsIns {
							if vot, ok := oneVoteIns.(Vote); ok {
								pow := vot.Power()
								coreVote := vot.Vote()
								if coreVote.IsApproved() {
									approvedPow += pow
									continue
								}

								if coreVote.IsNeutral() {
									neutralPow += pow
									continue
								}

								disApprovedPow += vot.Power()
								continue
							}

							log.Printf("the given entity (ID: %s) is not a valid request Vote instance", oneVoteIns.ID().String())
							continue
						}

						// render:
						datSet, datSetErr := toDataSet(votesPS)
						if datSetErr != nil {
							w.WriteHeader(http.StatusInternalServerError)
							str := fmt.Sprintf("there was an error while converting the vote entity set to data: %s", datSetErr.Error())
							w.Write([]byte(str))
							return
						}

						w.WriteHeader(http.StatusOK)
						params.Tmpl.Execute(w, &DataSetOfRequest{
							Request:        active_request.SDKFunc.ToData(req),
							Keyname:        keyname.SDKFunc.ToData(kname),
							MyUsers:        user.SDKFunc.ToDataSet(myUsersPS),
							Votes:          datSet,
							ApprovedPow:    approvedPow,
							DisapprovedPow: disApprovedPow,
							NeutralPow:     neutralPow,
							TotalPow:       approvedPow + disApprovedPow + neutralPow,
						})
						return
					}

					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the id is mandatory")
					w.Write([]byte(str))

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
