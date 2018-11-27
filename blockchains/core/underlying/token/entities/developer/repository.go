package developer

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
)

type repository struct {
	entityRepository entity.Repository
	devMetaData      entity.MetaData
}

func createRepository(entityRepository entity.Repository, devMetaData entity.MetaData) Repository {
	out := repository{
		entityRepository: entityRepository,
		devMetaData:      devMetaData,
	}

	return &out
}

// RetrieveByUser retrieves a Developer by its user
func (app *repository) RetrieveByUser(usr user.User) (Developer, error) {
	devIns, devInsErr := app.entityRepository.RetrieveByIntersectKeynames(app.devMetaData, []string{
		retrieveDevelopersByUserIDKeyname(usr.ID()),
	})

	if devInsErr != nil {
		return nil, devInsErr
	}

	if dev, ok := devIns.(Developer); ok {
		return dev, nil
	}

	str := fmt.Sprintf("the returned entity (ID: %s) was expected to be a Developer instance", devIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveByPledge retrieves a Developer by its pledge
func (app *repository) RetrieveByPledge(pldge pledge.Pledge) (Developer, error) {
	devIns, devInsErr := app.entityRepository.RetrieveByIntersectKeynames(app.devMetaData, []string{
		retrieveDevelopersByPledgeIDKeyname(pldge.ID()),
	})

	if devInsErr != nil {
		return nil, devInsErr
	}

	if dev, ok := devIns.(Developer); ok {
		return dev, nil
	}

	str := fmt.Sprintf("the returned entity (ID: %s) was expected to be a Developer instance", devIns.ID().String())
	return nil, errors.New(str)
}
