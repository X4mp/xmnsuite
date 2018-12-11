package user

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
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

// RetrieveByPubKeyAndWallet retrieves a User instance by its publicKey and wallet
func (app *repository) RetrieveByPubKeyAndWallet(pubKey crypto.PublicKey, wal wallet.Wallet) (User, error) {
	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.userMetaData, []string{
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
		retrieveUserByPubKeyKeyname(pubKey),
	}, index, amount)

	if insPSErr != nil {
		return nil, insPSErr
	}

	return insPS, nil
}
