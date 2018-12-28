package vote

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/crypto"
)

// CalculateFn represents the vote calculation func.
// First bool = concensus is reached
// Second bool = the vote passed
type CalculateFn func(votes entity.PartialSet) (bool, bool, error)

// CreateRouteFn creates a route
type CreateRouteFn func(ins Vote, rep entity.Representation) (string, error)

// Vote represents a request vote
type Vote interface {
	ID() *uuid.UUID
	Request() request.Request
	Voter() user.User
	Reason() string
	IsNeutral() bool
	IsApproved() bool
}

// Normalized represents a normalized Vote
type Normalized interface {
}

// Service represents the vote service
type Service interface {
	Save(ins Vote, rep entity.Representation) error
}

// Data represents human-redable data
type Data struct {
	ID         string
	Request    *request.Data
	Voter      *user.Data
	Reason     string
	IsNeutral  bool
	IsApproved bool
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Votes       []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID         *uuid.UUID
	Request    request.Request
	Voter      user.User
	Reason     string
	IsApproved bool
	IsNeutral  bool
}

// CreateSDKServiceParams represents the CreateSDKService params
type CreateSDKServiceParams struct {
	PK              crypto.PrivateKey
	Client          applications.Client
	CreateRouteFunc CreateRouteFn
}

// RouteNewParams represents the route new params
type RouteNewParams struct {
	PK                  crypto.PrivateKey
	Client              applications.Client
	FetchRepresentation func(groupName string, keyname string) (entity.Representation, error)
	EntityRepository    entity.Repository
}

// SDKFunc represents the vote SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Vote
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateSDKService     func(params CreateSDKServiceParams) Service
	ToData               func(vot Vote) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteNew             func(params RouteNewParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Vote {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createVote(params.ID, params.Request, params.Voter, params.Reason, params.IsNeutral, params.IsApproved)
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
	CreateSDKService: func(params CreateSDKServiceParams) Service {
		out := createSDKService(params.PK, params.Client, params.CreateRouteFunc)
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
	RouteNew: func(params RouteNewParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			if parseFormErr := r.ParseForm(); parseFormErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while parsing form elements: %s", parseFormErr.Error())
				w.Write([]byte(str))
				return
			}

			// retrieve the group name:
			vars := mux.Vars(r)
			if groupName, ok := vars["groupname"]; ok {
				// retrieve the keyname:
				if keynameAsString, ok := vars["keyname"]; ok {
					if requestIDAsString, ok := vars["id"]; ok {
						// create the repositories:
						voteService := createSDKService(params.PK, params.Client, func(ins Vote, rep entity.Representation) (string, error) {
							return fmt.Sprintf("/%s/requests/%s", rep.MetaData().Keyname(), ins.Request().ID().String()), nil
						})

						requestRepository := request.SDKFunc.CreateRepository(request.CreateRepositoryParams{
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

						myUserIDAsString := r.FormValue("myuser")
						decision := r.FormValue("decision")
						reason := r.FormValue("reason")

						isApproved := false
						isNeutral := false
						if decision == "is_approved" {
							isApproved = true
						}

						if decision == "is_neutral" {
							isNeutral = true
						}

						// parse the walletID:
						myUserID, myUserIDErr := uuid.FromString(myUserIDAsString)
						if myUserIDErr != nil {
							w.WriteHeader(http.StatusNotFound)
							str := fmt.Sprintf("the posted userID (%s) is invalid: %s", myUserID, myUserIDErr.Error())
							w.Write([]byte(str))
							return
						}

						// retrieve the user:
						myUsr, myUsrErr := userRepository.RetrieveByID(&myUserID)
						if myUsrErr != nil {
							w.WriteHeader(http.StatusInternalServerError)
							str := fmt.Sprintf("the posted user (ID: %s) could not be found: %s", myUserID.String(), myUsrErr.Error())
							w.Write([]byte(str))
							return
						}

						// create the vote:
						voteID := uuid.NewV4()
						vote, voteErr := createVote(&voteID, req, myUsr, reason, isNeutral, isApproved)
						if voteErr != nil {
							w.WriteHeader(http.StatusInternalServerError)
							str := fmt.Sprintf("there was an error while creating the Vote instance: %s", voteErr.Error())
							w.Write([]byte(str))
							return
						}

						// retrieve the vote representation:
						voteRepresentation, voteRepresentationErr := params.FetchRepresentation(groupName, keynameAsString)
						if voteRepresentationErr != nil {
							w.WriteHeader(http.StatusInternalServerError)
							str := fmt.Sprintf("there is no representation that match Group (%s) and Keyname (%s): %s", groupName, keynameAsString, voteRepresentationErr.Error())
							w.Write([]byte(str))
							return
						}

						// save the vote:
						saveVoteErr := voteService.Save(vote, voteRepresentation)
						if saveVoteErr != nil {
							w.WriteHeader(http.StatusInternalServerError)
							str := fmt.Sprintf("there was an error while saving the vote: %s", saveVoteErr.Error())
							w.Write([]byte(str))
							return
						}

						// redirect:
						url := fmt.Sprintf("/requests/%s/%s/%s", req.Request().Keyname().Group().Name(), req.Request().Keyname().Name(), req.ID().String())
						http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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
