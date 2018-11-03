package transfer

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/balance"
	"github.com/xmnservices/xmnsuite/blockchains/core/withdrawal"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
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
	balance, balanceErr := app.balanceRepository.RetrieveByWalletAndToken(ins.From(), ins.Token())
	if balanceErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the balance of the Wallet (ID: %s), for the Token (ID: %s): %s", ins.From().ID().String(), ins.Token().ID().String(), balanceErr.Error())
		return errors.New(str)
	}

	// make sure the balance is bigger than the transfer:
	if balance.Amount() < ins.Amount() {
		str := fmt.Sprintf("the balance of the wallet (ID: %s) for the token (ID: %s) is %d, but the transfer needed %d", balance.On().ID().String(), balance.Of().ID().String(), balance.Amount(), ins.Amount())
		return errors.New(str)
	}

	// create multiple withdrawals to make them harder to link with the transfer:
	compoundedAmount := 0
	withdrawals := []withdrawal.Withdrawal{}
	totalAmount := ins.Amount()
	amountOfWithdrawals := rand.Int() % 20
	amountPerWithdrawal := int(math.Ceil(float64(totalAmount / amountOfWithdrawals)))
	for i := 0; i < amountOfWithdrawals; i++ {

		// find the right amount to withdraw:
		amount := amountPerWithdrawal
		if (compoundedAmount + amount) > totalAmount {
			amount = totalAmount - compoundedAmount
		}

		// create the instance:
		wID := uuid.NewV4()
		withdrawals = append(withdrawals, withdrawal.SDKFunc.Create(withdrawal.CreateParams{
			ID:     &wID,
			From:   ins.From(),
			Amount: amount,
		}))

		// compound the amount:
		compoundedAmount += amount
	}

	// save the withdrawals instances:
	for _, oneWithdrawal := range withdrawals {
		saveErr := app.service.Save(oneWithdrawal, app.withdrawalRepresentation)
		if saveErr != nil {
			log.Printf("there was an error while saving a Withdrawal instance: %s", saveErr.Error())
		}
	}

	// save thre tranfer instance:
	saveErr := app.service.Save(ins, app.representation)
	if saveErr != nil {
		str := fmt.Sprintf("there was an error while saving a Transfer instance: %s", saveErr.Error())
		return errors.New(str)
	}

	return nil
}
