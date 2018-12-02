package sell

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/buy"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/sell"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/external"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/link"
)

type application struct {
	sellAmountToRetrievePerBatch int
	buyAmountToRetrievePerBatch  int
	sleepAfterUpdateDuration     time.Duration
	localLink                    link.Link
	transferRepresentation       entity.Representation
	pledgeRepresentation         entity.Representation
	walletMetaData               entity.MetaData
	internalSellRepository       sell.Repository
	externalSellRepository       sell.Repository
	internalBuyRepository        buy.Repository
	externalBuyService           buy.Service
	externalTransferRepository   transfer.Repository
	entityRepository             entity.Repository
	entityService                entity.Service
	stop                         bool
}

func createApplication(
	sellAmountToRetrievePerBatch int,
	buyAmountToRetrievePerBatch int,
	sleepAfterUpdateDuration time.Duration,
	localLink link.Link,
	transferRepresentation entity.Representation,
	pledgeRepresentation entity.Representation,
	walletMetaData entity.MetaData,
	internalSellRepository sell.Repository,
	externalSellRepository sell.Repository,
	internalBuyRepository buy.Repository,
	externalBuyService buy.Service,
	externalTransferRepository transfer.Repository,
	entityRepository entity.Repository,
	entityService entity.Service,
	stop bool,
) Daemon {
	out := application{
		sellAmountToRetrievePerBatch: sellAmountToRetrievePerBatch,
		buyAmountToRetrievePerBatch:  buyAmountToRetrievePerBatch,
		sleepAfterUpdateDuration:     sleepAfterUpdateDuration,
		localLink:                    localLink,
		transferRepresentation:       transferRepresentation,
		pledgeRepresentation:         pledgeRepresentation,
		walletMetaData:               walletMetaData,
		internalSellRepository:       internalSellRepository,
		externalSellRepository:       externalSellRepository,
		internalBuyRepository:        internalBuyRepository,
		externalBuyService:           externalBuyService,
		externalTransferRepository:   externalTransferRepository,
		entityRepository:             entityRepository,
		entityService:                entityService,
	}

	return &out
}

// Start starts the daemon
func (app *application) Start() error {
	app.stop = false

	for {

		// sleep some time:
		log.Printf("Waiting %f seconds...", app.sleepAfterUpdateDuration.Seconds())
		time.Sleep(app.sleepAfterUpdateDuration)

		// if we stop:
		if app.stop {
			return nil
		}

		// execute the sell orders:
		app.executeSellOrders()

		// execute the buy orders:
		app.executeBuyOrders()

	}
}

// Stop stops the daemon
func (app *application) Stop() error {
	app.stop = true
	return nil
}

func (app *application) executeBuyOrders() error {
	index := 0
	for {

		// retrieve the buy instances:
		buyPS, buyPSErr := app.internalBuyRepository.RetrieveSet(index, app.buyAmountToRetrievePerBatch)
		if buyPSErr != nil {
			str := fmt.Sprintf("there was an errror while retrieving buy orders (index: %d, amount: %d): %s", index, app.buyAmountToRetrievePerBatch, buyPSErr.Error())
			return errors.New(str)
		}

		// for each buy orders:
		buysIns := buyPS.Instances()
		for _, oneBuyIns := range buysIns {
			if oneBuy, ok := oneBuyIns.(buy.Buy); ok {

				// retrieve the transfer on the external blockchain:
				extTransfer, extTransferErr := app.externalTransferRepository.RetrieveByID(oneBuy.Transfer().ResourceID())
				if extTransferErr != nil {
					// transfer does not exists!
				}

				// make sure the transfer matches the sell wishes:
				wish := oneBuy.Sell().Wish()
				if bytes.Compare(extTransfer.Deposit().Token().ID().Bytes(), wish.Token().ResourceID().Bytes()) != 0 {
					// invalid token
				}

				if extTransfer.Deposit().Amount() != wish.Amount() {
					// invalid amount
				}

				// there is a match, retrieve the to wallet:
				toWalletID := oneBuy.Sell().DepositToWallet().ResourceID()

				// save the transfer:
				_, trsfErr := app.transfer(oneBuy.Sell(), toWalletID)
				if trsfErr != nil {
					log.Printf(trsfErr.Error())
					continue
				}

				continue

			}

			log.Printf("the entity (ID: %s) was expected to be a Buy instance", oneBuyIns.ID().String())
			continue
		}

		if buyPS.IsLast() {
			return nil
		}

		// increment:
		index += app.buyAmountToRetrievePerBatch
	}

}

func (app *application) executeSellOrders() error {
	index := 0
	for {

		// retrieve all the sell orders:
		sellPS, sellPSErr := app.internalSellRepository.RetrieveSet(index, app.sellAmountToRetrievePerBatch)
		if sellPSErr != nil {
			str := fmt.Sprintf("there was an errror while retrieving sell orders (index: %d, amount: %d): %s", index, app.sellAmountToRetrievePerBatch, sellPSErr.Error())
			return errors.New(str)
		}

		// for each sell orders:
		sellsIns := sellPS.Instances()
		for _, oneSellIns := range sellsIns {
			if oneSell, ok := oneSellIns.(sell.Sell); ok {
				// get the wish:
				wish := oneSell.Wish()

				// ask the external blockchain if there is any sell order that matches ours:
				matchedSell, matchedSellErr := app.externalSellRepository.RetrieveMatch(wish)
				if matchedSellErr != nil {
					log.Printf("it appears the Sell instance (ID: %s) has no match: %s", oneSell.ID().String(), matchedSellErr.Error())
					continue
				}

				// there is a match, get the walletID:
				toWalletID := matchedSell.DepositToWallet().ResourceID()

				// transfer:
				trsf, trsfErr := app.transfer(oneSell, toWalletID)
				if trsfErr != nil {
					log.Printf(trsfErr.Error())
					continue
				}

				// save the buy instance on the external blockchain:
				saveBuyErr := app.externalBuyService.Save(buy.SDKFunc.Create(buy.CreateParams{
					Sell: matchedSell,
					Transfer: external.SDKFunc.Create(external.CreateParams{
						Link:         app.localLink,
						ResourceName: app.transferRepresentation.MetaData().Keyname(),
						ResourceID:   trsf.ID(),
					}),
				}))

				if saveBuyErr != nil {
					log.Printf("there was an error while saving a Buy instance on an external blockchain: %s", saveBuyErr.Error())
					continue
				}

				continue
			}

			log.Printf("the entity (ID: %s) was expected to be a Sell instance", oneSellIns.ID().String())
			continue
		}

		if sellPS.IsLast() {
			return nil
		}

		// increment:
		index += app.sellAmountToRetrievePerBatch
	}
}

func (app *application) transfer(sld sell.Sell, toWalletID *uuid.UUID) (transfer.Transfer, error) {
	// retrieve the to wallet:
	toWalletIns, toWalletInsErr := app.entityRepository.RetrieveByID(app.walletMetaData, toWalletID)
	if toWalletInsErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the to-Wallet instance (ID: %s) to execute a sell order: %s", toWalletID, toWalletInsErr.Error())
		return nil, errors.New(str)
	}

	if toWallet, ok := toWalletIns.(wallet.Wallet); ok {
		// delete the pledge:
		delPledgeErr := app.entityService.Delete(sld.From(), app.pledgeRepresentation)
		if delPledgeErr != nil {
			str := fmt.Sprintf("there was an error while deleting a Pledge instance: %s", delPledgeErr.Error())
			return nil, errors.New(str)
		}

		// create the transfer instance:
		with := sld.From().From()
		trsf := transfer.SDKFunc.Create(transfer.CreateParams{
			Withdrawal: with,
			Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
				To:     toWallet,
				Token:  with.Token(),
				Amount: with.Amount(),
			}),
		})

		// transfer the tokens:
		saveTransferErr := app.entityService.Save(trsf, app.transferRepresentation)
		if saveTransferErr != nil {
			str := fmt.Sprintf("there was an error while saving a Transfer instance: %s", saveTransferErr.Error())
			return nil, errors.New(str)
		}

		// return:
		return trsf, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) was expected to be a Wallet instance", toWalletIns.ID().String())
	return nil, errors.New(str)
}
