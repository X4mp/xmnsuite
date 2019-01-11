package milestone

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
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

// RetrieveByID retrieves a milestone by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Milestone, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if mils, ok := ins.(Milestone); ok {
		return mils, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveByWallet retrieves a milestone by wallet
func (app *repository) RetrieveByWallet(wal wallet.Wallet) (Milestone, error) {
	keynames := []string{
		retrieveAllMilestoneKeyname(),
		retrieveMilestoneByWalletKeyname(wal),
	}

	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if insErr != nil {
		return nil, insErr
	}

	if mils, ok := ins.(Milestone); ok {
		return mils, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Milestone instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByCategory retrieves a milestone partial set by project
func (app *repository) RetrieveSetByProject(proj project.Project, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllMilestoneKeyname(),
		retrieveMilestoneByProjectKeyname(proj),
	}

	return app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
}
