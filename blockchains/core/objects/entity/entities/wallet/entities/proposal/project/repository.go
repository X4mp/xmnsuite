package project

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	approved_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
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

// RetrieveByProject retrieves a project by project
func (app *repository) RetrieveByProject(proj approved_project.Project) (Project, error) {
	keynames := []string{
		retrieveAllProjectKeyname(),
		retrieveProjectByProjectKeyname(proj),
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

// RetrieveByOwner retrieves a project by owner wallet
func (app *repository) RetrieveByOwner(owner wallet.Wallet) (Project, error) {
	keynames := []string{
		retrieveAllProjectKeyname(),
		retrieveProjectByOwnerWalletKeyname(owner),
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

// RetrieveByManager retrieves a project by manager wallet
func (app *repository) RetrieveByManager(mgr wallet.Wallet) (Project, error) {
	keynames := []string{
		retrieveAllProjectKeyname(),
		retrieveProjectByManagerWalletKeyname(mgr),
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

// RetrieveByLinker retrieves a project by linker wallet
func (app *repository) RetrieveByLinker(linker wallet.Wallet) (Project, error) {
	keynames := []string{
		retrieveAllProjectKeyname(),
		retrieveProjectByLinkerWalletKeyname(linker),
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
