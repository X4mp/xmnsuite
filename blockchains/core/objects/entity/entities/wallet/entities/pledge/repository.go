package pledge

import (
	"bytes"
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

// RetrieveByFromAndToWallet retrieves a pledge by from and to wallet
func (app *repository) RetrieveByFromAndToWallet(frm wallet.Wallet, to wallet.Wallet) (Pledge, error) {
	pldgePS, pldgePSErr := app.RetrieveSetByFromWallet(frm, 0, -1)
	if pldgePSErr != nil {
		return nil, pldgePSErr
	}

	pldges := pldgePS.Instances()
	for _, onePledgeIns := range pldges {
		if onePledge, ok := onePledgeIns.(Pledge); ok {
			if bytes.Compare(onePledge.To().ID().Bytes(), to.ID().Bytes()) == 0 {
				return onePledge, nil
			}

			continue
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", onePledgeIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("there is no pledge that match a from wallet (ID: %s) and to wallet (ID: %s)", frm.ID().String(), to.ID().String())
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
