package vote

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllVotesKeyname() string {
	return "votes"
}

func retrieveVotesByRequestIDKeyname(reqID *uuid.UUID) string {
	base := retrieveAllVotesKeyname()
	return fmt.Sprintf("%s:by_request_id:%s", base, reqID.String())
}

func retrieveVotesByVoterIDKeyname(voterID *uuid.UUID) string {
	base := retrieveAllVotesKeyname()
	return fmt.Sprintf("%s:by_voter_id:%s", base, voterID.String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Vote",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableVote) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				reqID, reqIDErr := uuid.FromString(storable.ReqID)
				if reqIDErr != nil {
					return nil, reqIDErr
				}

				voterID, voterIDErr := uuid.FromString(storable.VoterID)
				if voterIDErr != nil {
					return nil, voterIDErr
				}

				// retrieve the request:
				reqMet := request.SDKFunc.CreateMetaData()
				reqIns, reqInsErr := rep.RetrieveByID(reqMet, &reqID)
				if reqInsErr != nil {
					return nil, reqInsErr
				}

				// retrieve the user:
				usrMet := user.SDKFunc.CreateMetaData()
				usrIns, usrInsErr := rep.RetrieveByID(usrMet, &voterID)
				if usrInsErr != nil {
					return nil, usrInsErr
				}

				if req, ok := reqIns.(request.Request); ok {
					if usr, ok := usrIns.(user.User); ok {
						out, outErr := createVote(&id, req, usr, storable.Reason, storable.IsNeutral, storable.IsAppr)
						if outErr != nil {
							return nil, outErr
						}

						return out, nil
					}

					str := fmt.Sprintf("the entity (ID: %s) is not a valid user instance", voterID.String())
					return nil, errors.New(str)

				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid request instance", reqID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableVote); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedVote)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createVoteFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if vote, ok := ins.(Vote); ok {
				return createNormalizedVote(vote)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Vote instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedVote); ok {
				return createVoteFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Vote instance")
		},
		EmptyStorable:   new(storableVote),
		EmptyNormalized: new(normalizedVote),
	})
}

func createRepresentation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if vote, ok := ins.(Vote); ok {
				out := createStorableVote(vote)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Vote instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if vote, ok := ins.(Vote); ok {
				return []string{
					retrieveAllVotesKeyname(),
					retrieveVotesByRequestIDKeyname(vote.Request().ID()),
					retrieveVotesByVoterIDKeyname(vote.Voter().ID()),
				}, nil
			}

			return nil, errors.New("the given entity is not a valid Vote instance")
		},
		Sync: func(ds datastore.DataStore, ins entity.Entity) error {
			if vot, ok := ins.(Vote); ok {
				// metadata:
				metaData := createMetaData()
				requestMetaData := request.SDKFunc.CreateMetaData()

				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)

				// make sure the vote does not already exists:
				_, retVoteErr := repository.RetrieveByID(metaData, vot.ID())
				if retVoteErr == nil {
					str := fmt.Sprintf("the given Vote (ID: %s) already exists: %s", vot.ID().String(), retVoteErr.Error())
					return errors.New(str)
				}

				// make sure the request already exists:
				req, retRequestErr := repository.RetrieveByID(requestMetaData, vot.Request().ID())
				if retRequestErr != nil {
					str := fmt.Sprintf("the Request (ID: %s) does not exists: %s", req.ID().String(), retRequestErr.Error())
					return errors.New(str)
				}

				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Vote instance", ins.ID().String())
			return errors.New(str)
		},
	})
}

func toData(vot Vote) *Data {
	out := Data{
		ID:         vot.ID().String(),
		Request:    request.SDKFunc.ToData(vot.Request()),
		Voter:      user.SDKFunc.ToData(vot.Voter()),
		Reason:     vot.Reason(),
		IsNeutral:  vot.IsNeutral(),
		IsApproved: vot.IsApproved(),
	}

	return &out
}

func toDataSet(ins entity.PartialSet) (*DataSet, error) {
	data := []*Data{}
	instances := ins.Instances()
	for _, oneIns := range instances {
		if vot, ok := oneIns.(Vote); ok {
			data = append(data, toData(vot))
			continue
		}

		str := fmt.Sprintf("at least one of the elements (ID: %s) in the entity partial set is not a valid Vote instance", oneIns.ID().String())
		return nil, errors.New(str)
	}

	out := DataSet{
		Index:       ins.Index(),
		Amount:      ins.Amount(),
		TotalAmount: ins.TotalAmount(),
		IsLast:      ins.IsLast(),
		Votes:       data,
	}

	return &out, nil
}
