package helpers

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

func saveRequest(c *cliapp.Context, entityRepresentation entity.Representation, ins entity.Entity) (request.Request, error) {
	// retrieve conf with client:
	conf, client, confErr := retrieveConfWithClient(c)
	if confErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the conf instance: %s", confErr.Error())
		return nil, errors.New(str)
	}

	// create the request service:
	reqService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:          conf.WalletPK(),
		Client:      client,
		RoutePrefix: "",
	})

	// create the repositories:
	entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
		PK:     conf.WalletPK(),
		Client: client,
	})

	walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	// retrieve the walletID:
	walletIDAsString := c.String("walletid")
	walletID, walletIDErr := uuid.FromString(walletIDAsString)
	if walletIDErr != nil {
		str := fmt.Sprintf("the given walletid (ID: %s) is not a valid id", walletIDAsString)
		return nil, errors.New(str)
	}

	// retrieve the wallet:
	fromWallet, fromWalletErr := walletRepository.RetrieveByID(&walletID)
	if fromWalletErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", walletID.String(), fromWalletErr)
		return nil, errors.New(str)
	}

	// retrieve my user:
	fromUser, fromUSerErr := userRepository.RetrieveByPubKeyAndWallet(conf.WalletPK().PublicKey(), fromWallet)
	if fromUSerErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the user (PubKey: %s, WalletID): %s", conf.WalletPK().PublicKey(), fromWallet.ID().String())
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
		FromUser:  fromUser,
		NewEntity: ins,
		Reason:    c.String("reason"),
		Keyname:   kname,
	})

	saveErr := reqService.Save(newReq, entityRepresentation)
	if saveErr != nil {
		str := fmt.Sprintf("there was an error while saving the request: %s", saveErr.Error())
		return nil, errors.New(str)
	}

	return newReq, nil
}
