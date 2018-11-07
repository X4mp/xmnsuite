package user

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
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

// RetrieveByPubKey retrieves a User instance by its publicKey
func (app *repository) RetrieveByPubKey(pubKey crypto.PublicKey) (User, error) {
	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.userMetaData, []string{
		retrieveUserByPubKeyKeyname(pubKey),
	})

	if insErr != nil {
		return nil, insErr
	}

	if usr, ok := ins.(User); ok {
		return usr, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) retrieved (using pubKey: %s) is not a valid User instance", ins.ID().String(), pubKey.String())
	return nil, errors.New(str)
}
