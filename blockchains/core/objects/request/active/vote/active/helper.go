package active

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	core_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllVotesKeyname() string {
	return "activevotes"
}

func retrieveVotesByVoteIDKeyname(voteID *uuid.UUID) string {
	base := retrieveAllVotesKeyname()
	return fmt.Sprintf("%s:by_vote_id:%s", base, voteID.String())
}

func retrieveVotesByRequestIDKeyname(reqID *uuid.UUID) string {
	base := retrieveAllVotesKeyname()
	return fmt.Sprintf("%s:by_request_id:%s", base, reqID.String())
}

func retrieveVotesByVoterIDKeyname(voterID *uuid.UUID) string {
	base := retrieveAllVotesKeyname()
	return fmt.Sprintf("%s:by_voter_id:%s", base, voterID.String())
}

func retrieveVotesIsApprovedKeyname(isApproved bool) string {
	base := retrieveAllVotesKeyname()
	return fmt.Sprintf("%s:is_approved:%t", base, isApproved)
}

func retrieveVotesIsNeutralKeyname(isNeutral bool) string {
	base := retrieveAllVotesKeyname()
	return fmt.Sprintf("%s:is_neutral:%t", base, isNeutral)
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "ActiveVote",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableVote); ok {
				return createVoteFromStorable(rep, storable)
			}

			ptr := new(normalizedVote)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createVoteFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if vot, ok := ins.(Vote); ok {
				return createNormalizedVote(vot)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid active Vote instance", ins.ID().String())
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
			if vot, ok := ins.(Vote); ok {
				out := createStorableVote(vot)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid active Vote instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if vot, ok := ins.(Vote); ok {
				return []string{
					retrieveAllVotesKeyname(),
					retrieveVotesByVoteIDKeyname(vot.Vote().ID()),
					retrieveVotesByRequestIDKeyname(vot.Vote().Request().ID()),
					retrieveVotesByVoterIDKeyname(vot.Vote().Voter().ID()),
					retrieveVotesIsApprovedKeyname(vot.Vote().IsApproved()),
					retrieveVotesIsNeutralKeyname(vot.Vote().IsNeutral()),
				}, nil
			}

			return nil, errors.New("the given entity is not a valid active Vote instance")
		},
		Sync: func(ds datastore.DataStore, ins entity.Entity) error {
			if vot, ok := ins.(Vote); ok {
				// metadata:
				metaData := createMetaData()
				voteRepresentation := core_vote.SDKFunc.CreateRepresentation()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)
				entityService := entity.SDKFunc.CreateService(ds)
				repository := createRepository(entityRepository, metaData)

				// make sure the vote does not exists:
				_, retVoteErr := entityRepository.RetrieveByID(metaData, vot.ID())
				if retVoteErr == nil {
					str := fmt.Sprintf("the Vote (ID: %s) already exists", vot.ID().String())
					return errors.New(str)
				}

				// make sure the voter did not already vote:
				_, retVoteByVoterErr := repository.RetrieveByRequestVoter(vot.Vote().Voter(), vot.Vote().Request())
				if retVoteByVoterErr == nil {
					str := fmt.Sprintf("the Request (ID: %s) has already been voted on by the given Voter (ID: %s)", vot.Vote().Request().ID().String(), vot.Vote().Voter().ID().String())
					return errors.New(str)
				}

				// make sure the core vote does not exits, then save it:
				_, retPrevVoteErr := entityRepository.RetrieveByID(voteRepresentation.MetaData(), vot.Vote().ID())
				if retPrevVoteErr == nil {
					str := fmt.Sprintf("the given Vote (ID: %s) already exists: %s", vot.Vote().ID().String(), retPrevVoteErr.Error())
					return errors.New(str)
				}

				// save the core vote:
				saveReqErr := entityService.Save(vot.Vote(), voteRepresentation)
				if saveReqErr != nil {
					return saveReqErr
				}

				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid active Vote instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
