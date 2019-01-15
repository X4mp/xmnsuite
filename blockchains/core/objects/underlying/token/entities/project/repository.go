package project

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal"
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

// RetrieveByID retrieves a project by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Project, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if proj, ok := ins.(Project); ok {
		return proj, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveByProposal retrieves a project by proposal
func (app *repository) RetrieveByProposal(prop proposal.Proposal) (Project, error) {
	keynames := []string{
		retrieveAllProjectKeyname(),
		retrieveProjectByProposalKeyname(prop),
	}

	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if insErr != nil {
		return nil, insErr
	}

	if proj, ok := ins.(Project); ok {
		return proj, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByCategory retrieves a project set by category
func (app *repository) RetrieveSetByCategory(cat category.Category, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllProjectKeyname(),
		retrieveProjectByCategoryWalletKeyname(cat),
	}

	projPS, projPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if projPSErr != nil {
		return nil, projPSErr
	}

	return projPS, nil
}
