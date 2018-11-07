package genesis

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type service struct {
	repository        entity.Repository
	service           entity.Service
	genRepresentation entity.Representation
	genMetaData       entity.MetaData
}

func createService(serv entity.Service, repository entity.Repository, genRepresentation entity.Representation, genMetaData entity.MetaData) Service {
	out := service{
		repository:        repository,
		service:           serv,
		genRepresentation: genRepresentation,
		genMetaData:       genMetaData,
	}

	return &out
}

// Save saves an InitialDeposit instance
func (app *service) Save(ins Genesis) error {
	// if there is already a gensis instance, return an error:
	retGen, retGenErr := app.repository.RetrieveByIntersectKeynames(app.genMetaData, []string{keyname()})
	if retGenErr == nil {
		str := fmt.Sprintf("an genesis instance has already been created (ID: %s)", retGen.ID().String())
		return errors.New(str)
	}

	// save the genesis instance:
	saveErr := app.service.Save(ins, app.genRepresentation)
	if saveErr != nil {
		return saveErr
	}

	return nil
}
