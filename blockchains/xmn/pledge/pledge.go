package xmn

import (
	"errors"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
	"github.com/xmnservices/xmnsuite/datastore"
)

type storedPledge struct {
	ID   string `json:"id"`
	From string `json:"from_wallet_id"`
	To   string `json:"to_wallet_id"`
	Am   int    `json:"amount"`
}

func createStoredPledge(pledge Pledge) *storedPledge {
	out := storedPledge{
		ID:   pledge.ID().String(),
		From: pledge.From().ID().String(),
		To:   pledge.To().ID().String(),
		Am:   pledge.Amount(),
	}

	return &out
}

type pledge struct {
	UUID       *uuid.UUID `json:"id"`
	FromWallet Wallet     `json:"from"`
	ToWallet   Wallet     `json:"to"`
	Am         int        `json:"amount"`
}

func createPledge(id *uuid.UUID, frm Wallet, to Wallet, amount int) Pledge {
	out := pledge{
		UUID:       id,
		FromWallet: frm,
		ToWallet:   to,
		Am:         amount,
	}

	return &out
}

// ID returns the ID
func (obj *pledge) ID() *uuid.UUID {
	return obj.UUID
}

// From returns the from wallet
func (obj *pledge) From() Wallet {
	return obj.FromWallet
}

// To returns the to wallet
func (obj *pledge) To() Wallet {
	return obj.ToWallet
}

// Amount returns the amount
func (obj *pledge) Amount() int {
	return obj.Am
}

type pledgePartialSet struct {
	Plds  []Pledge `json:"pledges"`
	Idx   int      `json:"index"`
	TotAm int      `json:"total_amount"`
}

func createPledgePartialSet(pledges []Pledge, idx int, totAm int) PledgePartialSet {
	out := pledgePartialSet{
		Plds:  pledges,
		Idx:   idx,
		TotAm: totAm,
	}

	return &out
}

// Pledges returns the pledges
func (obj *pledgePartialSet) Pledges() []Pledge {
	return obj.Plds
}

// Index returns the index
func (obj *pledgePartialSet) Index() int {
	return obj.Idx
}

// Amount returns the amount
func (obj *pledgePartialSet) Amount() int {
	return len(obj.Plds)
}

// TotalAmount returns the totalAmount
func (obj *pledgePartialSet) TotalAmount() int {
	return obj.TotAm
}

type pledgeService struct {
	name          string
	keyname       string
	store         datastore.DataStore
	walletService WalletService
}

func createPledgeService(store datastore.DataStore, walletService WalletService) PledgeService {
	name := "Pledge"
	out := pledgeService{
		name:          name,
		keyname:       strings.ToLower(name),
		store:         store,
		walletService: walletService,
	}

	return &out
}

// Save saves a Pledge instance
func (app *pledgeService) Save(pledge Pledge) error {
	return app.entityService([]string{
		app.keynameByFromWalletID(pledge.From().ID()),
		app.keynameByToWalletID(pledge.To().ID()),
	}).Save(pledge)
}

// RetrieveByID retrieves a Pledge by ID
func (app *pledgeService) RetrieveByID(id *uuid.UUID) (Pledge, error) {
	ins, insErr := app.entityServiceWithEmptyKeys().RetrieveByID(id)
	if insErr != nil {
		return nil, insErr
	}

	return app.fromEntityToPledge(ins)
}

// RetrieveByFromWalletID retrieves a PledgePartialSet instance by the from WalletID
func (app *pledgeService) RetrieveByFromWalletID(fromWalletID *uuid.UUID, index int, amount int) (PledgePartialSet, error) {
	keyname := app.keynameByFromWalletID(fromWalletID)
	return app.retrievePartialSetByKeyname(keyname, index, amount)
}

// RetrieveByToWalletID retrieves a PledgePartialSet instance by the to WalletID
func (app *pledgeService) RetrieveByToWalletID(toWalletID *uuid.UUID, index int, amount int) (PledgePartialSet, error) {
	keyname := app.keynameByToWalletID(toWalletID)
	return app.retrievePartialSetByKeyname(keyname, index, amount)
}

// Populate populates a storable instance to a Pledge instance
func (app *pledgeService) Populate(stored *storedPledge) (Pledge, error) {
	// cast the ID:
	id, idErr := uuid.FromString(stored.ID)
	if idErr != nil {
		return nil, idErr
	}

	// cast the fromWalletID:
	fromWalletID, fromWalletIDErr := uuid.FromString(stored.From)
	if fromWalletIDErr != nil {
		return nil, fromWalletIDErr
	}

	// cast the toWalletID:
	toWalletID, toWalletIDErr := uuid.FromString(stored.To)
	if toWalletIDErr != nil {
		return nil, toWalletIDErr
	}

	// retrieve the from wallet:
	fromWallet, fromWalletErr := app.walletService.RetrieveByID(&fromWalletID)
	if fromWalletErr != nil {
		return nil, fromWalletErr
	}

	// retrieve the to wallet:
	toWallet, toWalletErr := app.walletService.RetrieveByID(&toWalletID)
	if toWalletErr != nil {
		return nil, toWalletErr
	}

	out := createPledge(&id, fromWallet, toWallet, stored.Am)
	return out, nil
}

func (app *pledgeService) fromEntityToPledge(ins entity.Entity) (Pledge, error) {
	if pledge, ok := ins.(Pledge); ok {
		return pledge, nil
	}

	return nil, errors.New("invalid entity type")
}

func (app *pledgeService) fromEntityToPledgePartialSet(set entity.EntityPartialSet) (PledgePartialSet, error) {
	pledges := []Pledge{}
	for _, oneIns := range set.Instances() {
		onePledge, onePledgeErr := app.fromEntityToPledge(oneIns)
		if onePledgeErr != nil {
			return nil, onePledgeErr
		}

		pledges = append(pledges, onePledge)
	}

	return createPledgePartialSet(pledges, set.Index(), set.TotalAmount()), nil
}

func (app *pledgeService) retrievePartialSetByKeyname(keyname string, index int, amount int) (PledgePartialSet, error) {
	entityPartialSet, entityPartialSetErr := app.entityServiceWithEmptyKeys().RetrieveSetByKeyname(keyname, index, amount)
	if entityPartialSetErr != nil {
		return nil, entityPartialSetErr
	}

	return app.fromEntityToPledgePartialSet(entityPartialSet)
}

func (app *pledgeService) keynameByID(id *uuid.UUID) string {
	return fmt.Sprintf("%s:by_id:%s", app.keyname, id.String())
}

func (app *pledgeService) keynameByFromWalletID(fromWalletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_from_wallet_id:%s", app.keyname, fromWalletID.String())
}

func (app *pledgeService) keynameByToWalletID(toWalletID *uuid.UUID) string {
	return fmt.Sprintf("%s:by_to_wallet_id:%s", app.keyname, toWalletID.String())
}

func (app *pledgeService) entityServiceWithEmptyKeys() entity.EntityService {
	return app.entityService([]string{})
}

func (app *pledgeService) entityService(keys []string) entity.EntityService {
	return entity.SDKFunc.CreateEntityService(entity.CreateEntityServiceParams{
		Met: entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
			Name: app.name,
			Keys: keys,
			ToData: func(ins entity.Entity) (interface{}, error) {
				if pledge, ok := ins.(Pledge); ok {
					storedPledge := createStoredPledge(pledge)
					return storedPledge, nil
				}

				return nil, errors.New("the given entity is not a valid Pledge instance")

			},
			ToEntity: func(storable interface{}) (entity.Entity, error) {
				if storablePledge, ok := storable.(*storedPledge); ok {
					pledge, pledgeErr := app.Populate(storablePledge)
					if pledgeErr != nil {
						str := fmt.Sprintf("there was an error while converting a storable Pledge instance to a Pledge instance: %s", pledgeErr.Error())
						return nil, errors.New(str)
					}

					return pledge, nil
				}

				return nil, errors.New("the given storable instance is not a valid storable Pledge instance")
			},
			EmptyStorable: new(storedPledge),
		}),
		DS: app.store,
	})
}
