package vote

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
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

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Vote instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveByRequestVoter retrieves a vote by request and voter
func (app *repository) RetrieveByRequestVoter(req request.Request, voter user.User) (Vote, error) {
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

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Vote instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByRequest retrieves a vote set by request
func (app *repository) RetrieveSetByRequest(req request.Request, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveVotesByRequestIDKeyname(req.ID()),
	}

	votePS, votePSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if votePSErr != nil {
		return nil, votePSErr
	}

	return votePS, nil
}
