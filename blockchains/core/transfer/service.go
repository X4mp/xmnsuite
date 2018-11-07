package transfer

import (
	"errors"
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/balance"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type service struct {
	repository               entity.Repository
	service                  entity.Service
	metaData                 entity.MetaData
	representation           entity.Representation
	withdrawalRepresentation entity.Representation
	balanceRepository        balance.Repository
}

func createService(
	repository entity.Repository,
	serv entity.Service,
	metaData entity.MetaData,
	representation entity.Representation,
	withdrawalRepresentation entity.Representation,
	balanceRepository balance.Repository,
) Service {
	out := service{
		repository:               repository,
		service:                  serv,
		metaData:                 metaData,
		representation:           representation,
		withdrawalRepresentation: withdrawalRepresentation,
		balanceRepository:        balanceRepository,
	}

	return &out
}

// Save saves a Transfer instance
func (app *service) Save(ins Transfer) error {
	// make sure the transfer does not already exists:
	_, transErr := app.repository.RetrieveByID(app.metaData, ins.ID())
	if transErr == nil {
		str := fmt.Sprintf("the Transfer instance (ID: %s) already exists", ins.ID().String())
		return errors.New(str)
	}

	// retrieve the balance:
	withdr := ins.Withdrawal()
	balance, balanceErr := app.balanceRepository.RetrieveByWalletAndToken(withdr.From(), withdr.Token())
	if balanceErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the balance of the Wallet (ID: %s), for the Token (ID: %s): %s", withdr.From().ID().String(), withdr.Token().ID().String(), balanceErr.Error())
		return errors.New(str)
	}

	// make sure the balance is bigger than the transfer:
	if balance.Amount() < withdr.Amount() {
		str := fmt.Sprintf("the balance of the wallet (ID: %s) for the token (ID: %s) is %d, but the transfer needed %d", balance.On().ID().String(), balance.Of().ID().String(), balance.Amount(), withdr.Amount())
		return errors.New(str)
	}

	// execute the withdrawal:
	wID := uuid.NewV4()
	withIns := withdrawal.SDKFunc.Create(withdrawal.CreateParams{
		ID:     &wID,
		From:   withdr.From(),
		Amount: withdr.Amount(),
	})

	// save the withdrawal instance:
	saveWithdraErr := app.service.Save(withIns, app.withdrawalRepresentation)
	if saveWithdraErr != nil {
		log.Printf("there was an error while saving a Withdrawal instance: %s", saveWithdraErr.Error())
	}

	// save thre tranfer instance:
	saveErr := app.service.Save(ins, app.representation)
	if saveErr != nil {
		str := fmt.Sprintf("there was an error while saving a Transfer instance: %s", saveErr.Error())
		return errors.New(str)
	}

	return nil
}
