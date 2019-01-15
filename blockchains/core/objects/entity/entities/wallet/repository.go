package wallet

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/crypto"
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

// RetrieveByID retrieves the wallet by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Wallet, error) {
	walIns, walInsErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if walInsErr != nil {
		return nil, walInsErr
	}

	if wal, ok := walIns.(Wallet); ok {
		return wal, nil
	}

	str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", walIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveByName retrieves the wallet by name
func (app *repository) RetrieveByName(name string) (Wallet, error) {
	keynames := []string{
		retrieveByNameKeyname(name),
	}

	walIns, walInsErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if walInsErr != nil {
		return nil, walInsErr
	}

	if wal, ok := walIns.(Wallet); ok {
		return wal, nil
	}

	str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", walIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByCreatorPublicKey retrieves the wallet set by public key
func (app *repository) RetrieveSetByCreatorPublicKey(pubKey crypto.PublicKey, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllWalletKeyname(),
		retrieveByPublicKeyWalletKeyname(pubKey),
	}

	walPS, walPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if walPSErr != nil {
		return nil, walPSErr
	}

	return walPS, nil
}

// RetrieveSet retrieves the wallet by set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllWalletKeyname()
	walPS, walPSErr := app.entityRepository.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if walPSErr != nil {
		return nil, walPSErr
	}

	return walPS, nil
}
