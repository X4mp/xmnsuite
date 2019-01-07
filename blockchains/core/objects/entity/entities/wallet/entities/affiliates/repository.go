package affiliates

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
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

// RetrieveByID retrieves a Affiliate instance by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Affiliate, error) {
	affIns, affInsErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if affInsErr != nil {
		return nil, affInsErr
	}

	if aff, ok := affIns.(Affiliate); ok {
		return aff, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Affiliate instance", affIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveByWallet retrieves a Affiliate by wallet
func (app *repository) RetrieveByWallet(wal wallet.Wallet) (Affiliate, error) {
	keynames := []string{
		retrieveAllAffiliatesKeyname(),
		retrieveAffiliatesByWalletKeyname(wal),
	}

	affIns, affInsErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if affInsErr != nil {
		return nil, affInsErr
	}

	if aff, ok := affIns.(Affiliate); ok {
		return aff, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Affiliate instance", affIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet retrieves a Affiliate set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllAffiliatesKeyname()
	affPS, affPSErr := app.entityRepository.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if affPSErr != nil {
		return nil, affPSErr
	}

	return affPS, nil
}
