package vote

import (
	"errors"
	"fmt"
	"log"

	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
)

type voteService struct {
	repository              entity.Repository
	service                 entity.Service
	voteRepresentation      entity.Representation
	voteMetaData            entity.MetaData
	requestRepresentation   entity.Representation
	newEntityRepresentation entity.Representation
}

func createVoteService(repository entity.Repository, service entity.Service, voteRepresentation entity.Representation, voteMetaData entity.MetaData, requestRepresentation entity.Representation, newEntityRepresentation entity.Representation) Service {
	out := voteService{
		repository:              repository,
		service:                 service,
		voteRepresentation:      voteRepresentation,
		voteMetaData:            voteMetaData,
		requestRepresentation:   requestRepresentation,
		newEntityRepresentation: newEntityRepresentation,
	}

	return &out
}

// Save saves a Vote instance
func (app *voteService) Save(vote Vote) error {
	// saves the entity:
	saveErr := app.service.Save(vote, app.voteRepresentation)
	if saveErr != nil {
		return saveErr
	}

	//retrieve all the votes by requestID:
	req := vote.Request()
	reqID := req.ID()
	keyname := retrieveVotesByRequestIDKeyname(reqID)
	votes, votesErr := app.repository.RetrieveSetByKeyname(app.voteMetaData, keyname, 0, -1)
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
	if approved >= neededConcensus {
		// insert the new entity:
		newEntity := req.New()
		saveNewErr := app.service.Save(newEntity, app.newEntityRepresentation)
		if saveNewErr != nil {
			str := fmt.Sprintf("there was an error while saving the new Entity instance (ID: %s): %s", newEntity.ID().String(), saveNewErr.Error())
			return errors.New(str)
		}
	}

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
