package helpers

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

func saveRequest(c *cliapp.Context, entityRepresentation entity.Representation, saveIns entity.Entity, delIns entity.Entity) (request.Request, error) {
	// retrieve conf with client:
	conf, client, confErr := retrieveConfWithClient(c)
	if confErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the conf instance: %s", confErr.Error())
		return nil, errors.New(str)
	}

	// create the request service:
	reqService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:     conf.WalletPK(),
		Client: client,
	})

	// create the repositories:
	entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
		PK:     conf.WalletPK(),
		Client: client,
	})

	userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	// retrieve my user:
	fromUser, fromUserErr := userRepository.RetrieveByPubKey(conf.WalletPK().PublicKey())
	if fromUserErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the user (PubKey: %s): %s", conf.WalletPK().PublicKey(), fromUserErr.Error())
		return nil, errors.New(str)
	}

	// retrieve the keyname:
	kname, knameErr := keynameRepository.RetrieveByName(entityRepresentation.MetaData().Keyname())
	if knameErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the keyname instance (Name: %s): %s", entityRepresentation.MetaData().Keyname(), knameErr.Error())
		return nil, errors.New(str)
	}

	// create the new request:
	newReq := request.SDKFunc.Create(request.CreateParams{
		FromUser:     fromUser,
		SaveEntity:   saveIns,
		DeleteEntity: delIns,
		Reason:       c.String("reason"),
		Keyname:      kname,
	})

	saveErr := reqService.Save(newReq, entityRepresentation)
	if saveErr != nil {
		str := fmt.Sprintf("there was an error while saving the request: %s", saveErr.Error())
		return nil, errors.New(str)
	}

	return newReq, nil
}
