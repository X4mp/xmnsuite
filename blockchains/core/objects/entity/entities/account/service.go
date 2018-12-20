package account

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
)

type service struct {
	userRepository     user.Repository
	walletRepository   wallet.Repository
	entityService      entity.Service
	userRepresentation entity.Representation
}

func createService(userRepository user.Repository, walletRepository wallet.Repository, entityService entity.Service, userRepresentation entity.Representation) Service {
	out := service{
		userRepository:     userRepository,
		walletRepository:   walletRepository,
		entityService:      entityService,
		userRepresentation: userRepresentation,
	}

	return &out
}

// Save saves a new account
func (app *service) Save(ins Account, amountOfWorkToVerify int) error {
	// make sure the user does not exists:
	_, retUserErr := app.userRepository.RetrieveByID(ins.User().ID())
	if retUserErr == nil {
		str := fmt.Sprintf("the given User (ID: %s) already exists and therefore can't be added in a new account", ins.User().ID().String())
		return errors.New(str)
	}

	// make sure the wallet of the user does not exists:
	_, retWalletErr := app.walletRepository.RetrieveByID(ins.User().Wallet().ID())
	if retWalletErr == nil {
		str := fmt.Sprintf("the given Wallet (ID: %s) already exists and therefore can't be added in a new account", ins.User().Wallet().ID().String())
		return errors.New(str)
	}

	// verify the work:
	verErr := ins.Work().PartialVerify(amountOfWorkToVerify)
	if verErr != nil {
		str := fmt.Sprintf("there was an error while veryfing the work: %s", verErr.Error())
		return errors.New(str)
	}

	// save the new user:
	saveErr := app.entityService.Save(ins.User(), app.userRepresentation)
	if saveErr != nil {
		str := fmt.Sprintf("there was an error while saving a User instance: %s", saveErr.Error())
		return errors.New(str)
	}

	return nil
}
