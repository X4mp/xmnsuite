package vote

import (
	"errors"
	"fmt"
	"log"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type voteService struct {
	repository            entity.Repository
	service               entity.Service
	voteRepresentation    entity.Representation
	requestRepresentation entity.Representation
}

func createVoteService(
	repository entity.Repository,
	service entity.Service,
	voteRepresentation entity.Representation,
	requestRepresentation entity.Representation,
) Service {
	out := voteService{
		repository:            repository,
		service:               service,
		voteRepresentation:    voteRepresentation,
		requestRepresentation: requestRepresentation,
	}

	return &out
}

// Save saves a Vote instance
func (app *voteService) Save(vote Vote, rep entity.Representation) error {
	// make sure the vote doesnt exists:
	_, retVoteErr := app.repository.RetrieveByID(app.voteRepresentation.MetaData(), vote.ID())
	if retVoteErr == nil {
		str := fmt.Sprintf("the Vote (ID: %s) already exists", vote.ID().String())
		return errors.New(str)
	}

	// make sure the voter did not already vote:
	_, retVoteByVoterErr := app.repository.RetrieveByIntersectKeynames(app.voteRepresentation.MetaData(), []string{
		retrieveVotesByRequestIDKeyname(vote.Request().ID()),
		retrieveVotesByVoterIDKeyname(vote.Voter().ID()),
	})

	if retVoteByVoterErr == nil {
		str := fmt.Sprintf("the Request (ID: %s) has already been voted on by the given Voter (ID: %s)", vote.Request().ID().String(), vote.Voter().ID().String())
		return errors.New(str)
	}

	// saves the vote:
	saveErr := app.service.Save(vote, app.voteRepresentation)
	if saveErr != nil {
		return saveErr
	}

	//retrieve all the votes by requestID:
	req := vote.Request()
	reqID := req.ID()
	keyname := retrieveVotesByRequestIDKeyname(reqID)
	votes, votesErr := app.repository.RetrieveSetByKeyname(app.voteRepresentation.MetaData(), keyname, 0, -1)
	if votesErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the vote partial set related to the Request (ID: %s): %s", reqID.String(), votesErr.Error())
		return errors.New(str)
	}

	// retrieve the needed concensus from the requester wallet:
	neededConcensus := vote.Request().From().Wallet().ConcensusNeeded()

	// compile the vote's concensus:
	approved := 0
	disapproved := 0
	votesIns := votes.Instances()
	for _, oneVoteIns := range votesIns {
		if oneVote, ok := oneVoteIns.(Vote); ok {
			if oneVote.IsApproved() {
				approved += oneVote.Voter().Shares()
				continue
			}

			disapproved += oneVote.Voter().Shares()
			continue
		}

		log.Printf("the entity (ID: %s) is not a valid Vote instance", oneVoteIns.ID().String())
	}

	// vote is approved, insert the new entity:
	isApproved := approved >= neededConcensus
	if isApproved {
		// insert the new entity:
		newEntity := req.New()
		saveNewErr := app.service.Save(newEntity, rep)
		if saveNewErr != nil {
			str := fmt.Sprintf("there was an error while saving the new Entity instance (ID: %s): %s", newEntity.ID().String(), saveNewErr.Error())
			return errors.New(str)
		}
	}

	// if vote is approved or disapproved, delete the votes and request:
	isDisapproved := disapproved >= neededConcensus
	if isApproved || isDisapproved {
		// delete the votes:
		for _, oneVote := range votesIns {
			delVoteErr := app.service.Delete(oneVote, app.voteRepresentation)
			if delVoteErr != nil {
				log.Printf("there was an error while deleting a Vote (ID: %s) after concensus was reached: %s", oneVote.ID().String(), delVoteErr.Error())
			}
		}

		// delete the request:
		delReqErr := app.service.Delete(req, app.requestRepresentation)
		if delReqErr != nil {
			log.Printf("there was an error while deleting a Request (ID: %s) after concensus was reached: %s", reqID.String(), delReqErr.Error())
		}

		return nil
	}

	// the voting is still going on:
	return nil
}
