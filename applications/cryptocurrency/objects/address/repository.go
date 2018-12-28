package address

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
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

// RetrieveByID retrieves an address by ID:
func (app *repository) RetrieveByID(id *uuid.UUID) (Address, error) {
	addrIns, addrInsErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if addrInsErr != nil {
		return nil, addrInsErr
	}

	if addr, ok := addrIns.(Address); ok {
		return addr, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Address instance", addrIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveByAddress retrieves an address by address:
func (app *repository) RetrieveByAddress(addr string) (Address, error) {
	keynames := []string{
		retrieveAllAddressKeyname(),
		retrieveAddressByAddressKeyname(addr),
	}

	addrIns, addrInsErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if addrInsErr != nil {
		return nil, addrInsErr
	}

	if addr, ok := addrIns.(Address); ok {
		return addr, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Address instance", addrIns.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet retrieves an address partial set:
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllAddressKeyname()
	addrPS, addrPSErr := app.entityRepository.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if addrPSErr != nil {
		return nil, addrPSErr
	}

	return addrPS, nil
}

// RetrieveSetByWallet retrieves an address partial set by wallet:
func (app *repository) RetrieveSetByWallet(wal wallet.Wallet, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllAddressKeyname(),
		retrieveAddressByWalletKeyname(wal),
	}

	addrPS, addrPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if addrPSErr != nil {
		return nil, addrPSErr
	}

	return addrPS, nil
}
