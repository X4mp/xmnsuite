package offer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

type repository struct {
	entityRepository entity.Repository
	metaData         entity.MetaData
}

func createRepository(metaData entity.MetaData, entityRepository entity.Repository) Repository {
	out := repository{
		metaData:         metaData,
		entityRepository: entityRepository,
	}

	return &out
}

// RetrieveByID retrieves an Offer by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Offer, error) {
	offIns, offInsErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if offInsErr != nil {
		return nil, offInsErr
	}

	if off, ok := offIns.(Offer); ok {
		return off, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Offer instance", offIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveByPledge retrieves an Offer by Pledge
func (app *repository) RetrieveByPledge(pldge pledge.Pledge) (Offer, error) {
	keynames := []string{
		retrieveOfferByPledge(pldge),
	}

	offIns, offInsErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if offInsErr != nil {
		return nil, offInsErr
	}

	if off, ok := offIns.(Offer); ok {
		return off, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Offer instance", offIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet retrieves an Offer set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllOffersKeyname()
	offPS, offPSErr := app.entityRepository.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if offPSErr != nil {
		return nil, offPSErr
	}

	return offPS, nil
}

// RetrieveSetByToAddress retrieves an Offer set by toAddress
func (app *repository) RetrieveSetByToAddress(toAddr address.Address, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllOffersKeyname(),
		retrieveOfferByToAddress(toAddr),
	}

	offPS, offPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if offPSErr != nil {
		return nil, offPSErr
	}

	return offPS, nil
}
