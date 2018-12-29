package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type repository struct {
	entityRepository entity.Repository
	metaData         entity.MetaData
}

func createRepository(metaData entity.MetaData, entityRepository entity.Repository) Repository {
	out := repository{
		entityRepository: entityRepository,
		metaData:         metaData,
	}

	return &out
}

// RetrieveByID retrieves a deposit by id
func (app *repository) RetrieveByID(id *uuid.UUID) (Deposit, error) {
	depIns, depInsErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if depInsErr != nil {
		return nil, depInsErr
	}

	if dep, ok := depIns.(Deposit); ok {
		return dep, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Deposit instance", depIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet retrieves a deposit partial set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllDepositsKeyname()
	depPS, depPSErr := app.entityRepository.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if depPSErr != nil {
		return nil, depPSErr
	}

	return depPS, nil
}

// RetrieveSetByOffer retrieves a deposit set by offer
func (app *repository) RetrieveSetByOffer(off offer.Offer, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllDepositsKeyname(),
		retrieveDepositByOfferKeyname(off),
	}

	depPS, depPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if depPSErr != nil {
		return nil, depPSErr
	}

	return depPS, nil
}

// RetrieveSetByFromAddress retrieves a deposit set by from address
func (app *repository) RetrieveSetByFromAddress(frmAddress address.Address, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllDepositsKeyname(),
		retrieveDepositByFromAddressKeyname(frmAddress),
	}

	depPS, depPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if depPSErr != nil {
		return nil, depPSErr
	}

	return depPS, nil
}
