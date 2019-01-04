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
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/configs"
)

func processWalletRequest(c *cliapp.Context, representation entity.Representation, storable interface{}) (request.Normalized, error) {
	// retrieve the configurations:
	fileAsString := c.String("file")
	confRepository := configs.SDKFunc.CreateRepository()
	conf, confErr := confRepository.Retrieve(fileAsString, c.String("pass"))
	if confErr != nil {
		str := fmt.Sprintf("the given file (%s) either does not exist or the given password is invalid", fileAsString)
		return nil, errors.New(str)
	}

	// convert the walletid:
	walletIDAsString := c.String("walletid")
	walID, walIDErr := uuid.FromString(walletIDAsString)
	if walIDErr != nil {
		str := fmt.Sprintf("the given walletID (ID: %s) is invalid, but mandatory", walletIDAsString)
		return nil, errors.New(str)
	}

	// metadata:
	metaData := representation.MetaData()

	// create the blockchain client:
	client := tendermint.SDKFunc.CreateClient(tendermint.CreateClientParams{
		IPAsString: c.String("host"),
	})

	// create the repositories:
	entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
		PK:          conf.WalletPK(),
		Client:      client,
		RoutePrefix: "",
	})

	// create the user repository:
	walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	// create the services:
	requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
		PK:          conf.WalletPK(),
		Client:      client,
		RoutePrefix: "",
	})

	// retrieve the wallet:
	wal, walErr := walletRepository.RetrieveByID(&walID)
	if walErr != nil {
		str := fmt.Sprintf("there was an error while retrieving a wallet (ID: %s): %s", walID.String(), walErr.Error())
		return nil, errors.New(str)
	}

	// convert the storable to an entity:
	ent, entErr := metaData.ToEntity()(entityRepository, storable)
	if entErr != nil {
		str := fmt.Sprintf("there was an error while converting the storable to an entity instance: %s", entErr.Error())
		return nil, errors.New(str)
	}

	// retrieve the from user:
	pubKey := conf.WalletPK().PublicKey()
	fromUser, fromUserErr := userRepository.RetrieveByPubKeyAndWallet(pubKey, wal)
	if fromUserErr != nil {
		str := fmt.Sprintf("there was an error while retrieving a user (pubKey: %s, walletID: %s): %s", pubKey.String(), wal.ID().String(), fromUserErr.Error())
		return nil, errors.New(str)
	}

	// retrieve the keyname:
	kname, knameErr := keynameRepository.RetrieveByName(metaData.Keyname())
	if knameErr != nil {
		str := fmt.Sprintf("there was an  error while retrieving the keyname (name; %s): %s", metaData.Keyname(), knameErr.Error())
		return nil, errors.New(str)
	}

	// create the request:
	req := request.SDKFunc.Create(request.CreateParams{
		FromUser:  fromUser,
		NewEntity: ent,
		Reason:    c.String("reason"),
		Keyname:   kname,
	})

	// save the request:
	saveErr := requestService.Save(req, representation)
	if saveErr != nil {
		str := fmt.Sprintf("there was an error while saving a request instance: %s", saveErr.Error())
		return nil, errors.New(str)
	}

	// normalize:
	normalized, normalizedErr := request.SDKFunc.CreateMetaData().Normalize()(req)
	if normalizedErr != nil {
		return nil, normalizedErr
	}

	return normalized, nil
}
