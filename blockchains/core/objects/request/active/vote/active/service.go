package active

import (
	"errors"
	"fmt"
	"log"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/completed"
)

type voteService struct {
	repository                     entity.Repository
	service                        entity.Service
	voteRepresentation             entity.Representation
	requestRepresentation          entity.Representation
	completedRequestRepresentation entity.Representation
}

func createVoteService(
	repository entity.Repository,
	service entity.Service,
	voteRepresentation entity.Representation,
	requestRepresentation entity.Representation,
	completedRequestRepresentation entity.Representation,
) Service {
	out := voteService{
		repository:                     repository,
		service:                        service,
		voteRepresentation:             voteRepresentation,
		requestRepresentation:          requestRepresentation,
		completedRequestRepresentation: completedRequestRepresentation,
	}

	return &out
}

// Save saves a Vote instance
func (app *voteService) Save(vote Vote, rep entity.Representation) error {
	// saves the vote:
	saveErr := app.service.Save(vote, app.voteRepresentation)
	if saveErr != nil {
		return saveErr
	}

	//retrieve all the votes by requestID:
	req := vote.Vote().Request()
	reqID := req.ID()
	keyname := retrieveVotesByRequestIDKeyname(reqID)
	votes, votesErr := app.repository.RetrieveSetByKeyname(app.voteRepresentation.MetaData(), keyname, 0, -1)
	if votesErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the vote partial set related to the Request (ID: %s): %s", reqID.String(), votesErr.Error())
		return errors.New(str)
	}

	// voting still going on...
	if votes.Amount() == 0 {
		return nil
	}

	// retrieve the concensus needed:
	votesIns := votes.Instances()
	if firstVote, ok := votesIns[0].(Vote); ok {
		// check the amount of concensus needed:
		neededConcensus := firstVote.Vote().Request().ConcensusNeeded()

		// calculate the balance:
		approved := 0
		disapproved := 0
		neutral := 0
		for _, oneVote := range votesIns {
			if vot, ok := oneVote.(Vote); ok {
				coreVote := vot.Vote()
				if coreVote.IsApproved() {
					approved += vot.Power()
					continue
				}

				if coreVote.IsNeutral() {
					neutral += vot.Power()
					continue
				}

				disapproved += vot.Power()
				continue
			}
		}

		isApproved := neededConcensus <= approved
		concensusReached := (approved + disapproved + neutral) >= neededConcensus

		// if concensus is reached:
		if concensusReached {
			// if vote is approved, insert the new entity:
			if isApproved {

				// get the request:
				reqReq := req.Request()

				// if we must save the entity:
				if reqReq.HasSave() {
					// save the entity:
					toSaveEntity := reqReq.Save()
					saveErr := app.service.Save(toSaveEntity, rep)
					if saveErr != nil {
						str := fmt.Sprintf("there was an error while saving the Entity instance (ID: %s): %s", toSaveEntity.ID().String(), saveErr.Error())
						return errors.New(str)
					}
				}

				// if we must delete the entity:
				if reqReq.HasDelete() {
					//delete the entity:
					toDeleteEntity := reqReq.Delete()
					delErr := app.service.Delete(toDeleteEntity, rep)
					if delErr != nil {
						str := fmt.Sprintf("there was an error while deleting the Entity instance (ID: %s): %s", toDeleteEntity.ID().String(), delErr.Error())
						return errors.New(str)
					}
				}
			}

			// delete the votes:
			for _, oneVote := range votesIns {
				delVoteErr := app.service.Delete(oneVote, app.voteRepresentation)
				if delVoteErr != nil {
					log.Printf("there was an error while deleting a Vote (ID: %s) after concensus was reached: %s", oneVote.ID().String(), delVoteErr.Error())
				}
			}

			// create the completed request:
			completedReq := completed.SDKFunc.Create(completed.CreateParams{
				Request:         req.Request(),
				ConcensusNeeded: neededConcensus,
				Approved:        approved,
				Disapproved:     disapproved,
				Neutral:         neutral,
			})

			// save the completed request:
			saveCompletedReqErr := app.service.Save(completedReq, app.completedRequestRepresentation)
			if saveCompletedReqErr != nil {
				str := fmt.Sprintf("there was an error while saving a completed request: %s", saveCompletedReqErr.Error())
				return errors.New(str)
			}

			// delete the active request:
			delActiveReqErr := app.service.Delete(req, app.requestRepresentation)
			if delActiveReqErr != nil {
				str := fmt.Sprintf("there was an error while deleting an ActiveRequest (ID: %s) after concensus was reached: %s", req.ID().String(), delActiveReqErr.Error())
				return errors.New(str)
			}

			// everything worked:
			return nil
		}

		// the voting is still going on:
		return nil

	}

	str := fmt.Sprintf("the entity (ID: %s) was expected to be a Vote instance", votesIns[0].ID().String())
	return errors.New(str)
}
