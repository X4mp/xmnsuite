package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

type repository struct {
	userMetaData     entity.MetaData
	entityRepository entity.Repository
}

func createRepository(userMetaData entity.MetaData, entityRepository entity.Repository) Repository {
	out := repository{
		userMetaData:     userMetaData,
		entityRepository: entityRepository,
	}

	return &out
}

// RetrieveByID retrieves a User instance by its ID
func (app *repository) RetrieveByID(id *uuid.UUID) (User, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.userMetaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if usr, ok := ins.(User); ok {
		return usr, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveByPubKeyAndWallet retrieves a User instance by its publicKey and wallet
func (app *repository) RetrieveByPubKeyAndWallet(pubKey crypto.PublicKey, wal wallet.Wallet) (User, error) {
	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.userMetaData, []string{
		retrieveAllUserKeyname(),
		retrieveUserByPubKeyKeyname(pubKey),
		retrieveUserByWalletIDKeyname(wal.ID()),
	})

	if insErr != nil {
		return nil, insErr
	}

	if usr, ok := ins.(User); ok {
		return usr, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) retrieved (using pubKey: %s, walletID: %s) is not a valid User instance", ins.ID().String(), pubKey.String(), wal.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByPubKey retrieves a list of users connected to this pubkey
func (app *repository) RetrieveSetByPubKey(pubKey crypto.PublicKey, index int, amount int) (entity.PartialSet, error) {
	insPS, insPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.userMetaData, []string{
		retrieveAllUserKeyname(),
		retrieveUserByPubKeyKeyname(pubKey),
	}, index, amount)

	if insPSErr != nil {
		return nil, insPSErr
	}

	return insPS, nil
}

// RetrieveSetByWallet retrieves a list of users connected to wallet
func (app *repository) RetrieveSetByWallet(wal wallet.Wallet, index int, amount int) (entity.PartialSet, error) {
	insPS, insPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.userMetaData, []string{
		retrieveAllUserKeyname(),
		retrieveUserByWalletIDKeyname(wal.ID()),
	}, index, amount)

	if insPSErr != nil {
		return nil, insPSErr
	}

	return insPS, nil
}
