package feature

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
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

// RetrieveByID retrieves a Feature by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Feature, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if feat, ok := ins.(Feature); ok {
		return feat, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Feature instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByProject retrieves a Feature set by Project
func (app *repository) RetrieveSetByProject(proj project.Project, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllFeatureKeyname(),
		retrieveFeatureByProjectKeyname(proj),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
}

// RetrieveSetByCreatedByUser retrieves a Feature by createdBy user
func (app *repository) RetrieveSetByCreatedByUser(createdBy user.User, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllFeatureKeyname(),
		retrieveFeatureByCreatedByUserKeyname(createdBy),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
}
