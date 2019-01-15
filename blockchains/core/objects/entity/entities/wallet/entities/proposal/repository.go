package proposal

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/category"
)

type repository struct {
	metaData         entity.MetaData
	entityRepository entity.Repository
}

func createRepository(metaData entity.MetaData, entityRepository entity.Repository) Repository {
	out := repository{
		metaData:         metaData,
		entityRepository: entityRepository,
	}

	return &out
}

// RetrieveByID retrieves a proposal by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Proposal, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if prop, ok := ins.(Proposal); ok {
		return prop, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Proposal instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByCategory retrieves a proposal partial set by category
func (app *repository) RetrieveSetByCategory(cat category.Category, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllProposalKeyname(),
		retrieveProposalByCategoryKeyname(cat),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
}
