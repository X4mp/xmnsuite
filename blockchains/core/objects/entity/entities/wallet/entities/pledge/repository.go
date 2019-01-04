package pledge

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

func createRepository(entityRepository entity.Repository, metaData entity.MetaData) Repository {
	out := repository{
		entityRepository: entityRepository,
		metaData:         metaData,
	}

	return &out
}

// RetrieveByID retrieves a pledge by id
func (app *repository) RetrieveByID(id *uuid.UUID) (Pledge, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if pldge, ok := ins.(Pledge); ok {
		return pldge, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByFromWallet retrieves a pledge partial set by from wallet
func (app *repository) RetrieveSetByFromWallet(frm wallet.Wallet, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllPledgesKeyname(),
		retrievePledgesByFromWalletKeyname(frm),
	}

	ps, psErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if psErr != nil {
		return nil, psErr
	}

	return ps, nil
}

// RetrieveSetByToWallet retrieves a pledge partial set by to wallet
func (app *repository) RetrieveSetByToWallet(to wallet.Wallet, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllPledgesKeyname(),
		retrievePledgesByToWalletKeyname(to),
	}

	ps, psErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if psErr != nil {
		return nil, psErr
	}

	return ps, nil
}
