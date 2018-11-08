package genesis

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type service struct {
	repository        Repository
	entityRepository  entity.Repository
	entityService     entity.Service
	genRepresentation entity.Representation
}

func createService(serv entity.Service, entityRepository entity.Repository, rep Repository, genRepresentation entity.Representation) Service {
	out := service{
		repository:        rep,
		entityService:     serv,
		entityRepository:  entityRepository,
		genRepresentation: genRepresentation,
	}

	return &out
}

// Save saves an InitialDeposit instance
func (app *service) Save(ins Genesis) error {
	// if there is already a gensis instance, return an error:
	_, retGenErr := app.repository.Retrieve()
	if retGenErr == nil {
		return retGenErr
	}

	// save the genesis instance:
	saveErr := app.entityService.Save(ins, app.genRepresentation)
	if saveErr != nil {
		return saveErr
	}

	return nil
}
