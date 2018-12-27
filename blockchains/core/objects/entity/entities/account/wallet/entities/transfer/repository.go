package transfer

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
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

// RetrieveByID retrieves a transfer instance by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Transfer, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if trsf, ok := ins.(Transfer); ok {
		return trsf, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Transfer instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet retrieves a transfer set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllTransfersKeyname()
	trsfPS, trsfPSErr := app.entityRepository.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if trsfPSErr != nil {
		return nil, trsfPSErr
	}

	return trsfPS, nil
}

// RetrieveSetByDeposit retrieves a transfer set by deposit
func (app *repository) RetrieveSetByDeposit(dep deposit.Deposit, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllTransfersKeyname(),
		retrieveTransfersByDepositKeyname(dep),
	}

	trsfPS, trsfPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if trsfPSErr != nil {
		return nil, trsfPSErr
	}

	return trsfPS, nil
}

// RetrieveSetByWithdrawal retrieves a transfer set by withdrawal
func (app *repository) RetrieveSetByWithdrawal(with withdrawal.Withdrawal, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllTransfersKeyname(),
		retrieveTransfersByWithdrawalKeyname(with),
	}

	trsfPS, trsfPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if trsfPSErr != nil {
		return nil, trsfPSErr
	}

	return trsfPS, nil
}
