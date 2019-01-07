package active

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	core_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
)

type repository struct {
	entityRepository entity.Repository
	metaData         entity.MetaData
}

func createRepository(entityRepository entity.Repository, metaData entity.MetaData) Repository {
	out := repository{
		entityRepository: entityRepository,
		metaData:         metaData,
	}

	return &out
}

// RetrieveByID retrieves a vote by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Vote, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if vot, ok := ins.(Vote); ok {
		return vot, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid active Vote instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveByVote retrieves a vote by core vote
func (app *repository) RetrieveByVote(vot core_vote.Vote) (Vote, error) {
	keynames := []string{
		retrieveAllVotesKeyname(),
		retrieveVotesByVoteIDKeyname(vot.ID()),
	}

	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if insErr != nil {
		return nil, insErr
	}

	if vot, ok := ins.(Vote); ok {
		return vot, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid active Vote instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveByVoterOnRequest retrieves an active vote by voter and requst
func (app *repository) RetrieveByRequestVoter(voter user.User, req active_request.Request) (Vote, error) {
	keynames := []string{
		retrieveAllVotesKeyname(),
		retrieveVotesByRequestIDKeyname(req.ID()),
		retrieveVotesByVoterIDKeyname(voter.ID()),
	}

	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if insErr != nil {
		return nil, insErr
	}

	if vot, ok := ins.(Vote); ok {
		return vot, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid active Vote instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByRequest retrieves a vote set by request
func (app *repository) RetrieveSetByRequest(req active_request.Request, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllVotesKeyname(),
		retrieveVotesByRequestIDKeyname(req.ID()),
	}

	votePS, votePSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if votePSErr != nil {
		return nil, votePSErr
	}

	return votePS, nil
}

// RetrieveSetByRequestWithDirection retrieves a vote set by request and direction
func (app *repository) RetrieveSetByRequestWithDirection(req active_request.Request, index int, amount int, isApproved bool, isNeutral bool) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllVotesKeyname(),
		retrieveVotesByRequestIDKeyname(req.ID()),
		retrieveVotesIsApprovedKeyname(isApproved),
		retrieveVotesIsNeutralKeyname(isNeutral),
	}

	votePS, votePSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if votePSErr != nil {
		return nil, votePSErr
	}

	return votePS, nil
}
