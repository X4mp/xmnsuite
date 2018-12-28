package address

import (
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
	return nil, nil
}

// RetrieveByAddress retrieves an address by address:
func (app *repository) RetrieveByAddress(addr []byte) (Address, error) {
	return nil, nil
}

// RetrieveSet retrieves an address partial set:
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	return nil, nil
}

// RetrieveSetByWallet retrieves an address partial set by wallet:
func (app *repository) RetrieveSetByWallet(wal wallet.Wallet, index int, amount int) (entity.PartialSet, error) {
	return nil, nil
}
