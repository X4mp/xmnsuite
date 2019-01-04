package information

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type service struct {
	repository        Repository
	entityRepository  entity.Repository
	entityService     entity.Service
	infoRepresentation entity.Representation
}

func createService(serv entity.Service, entityRepository entity.Repository, rep Repository, infoRepresentation entity.Representation) Service {
	out := service{
		repository:        rep,
		entityService:     serv,
		entityRepository:  entityRepository,
		infoRepresentation: infoRepresentation,
	}

	return &out
}

// Save saves an InitialDeposit instance
func (app *service) Save(ins Information) error {
	// if there is already a Information instance, return an error:
	_, retInfoErr := app.repository.Retrieve()
	if retInfoErr == nil {
		return retInfoErr
	}

	// save the information instance:
	saveErr := app.entityService.Save(ins, app.infoRepresentation)
	if saveErr != nil {
		return saveErr
	}

	return nil
}
