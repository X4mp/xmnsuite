package deposit

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Deposit represents the initial deposit
type Deposit interface {
	ID() *uuid.UUID
	To() wallet.Wallet
	Token() token.Token
	Amount() int
}

// Normalized represents the normalized deposit
type Normalized interface {
}

// Repository represents the deposit Repository
type Repository interface {
	RetrieveSetByToWalletAndToken(wal wallet.Wallet, tok token.Token) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	To     wallet.Wallet
	Token  token.Token
	Amount int
}

// CreateRepositoryParams represents a CreateRepository params
type CreateRepositoryParams struct {
	Datastore        datastore.DataStore
	EntityRepository entity.Repository
}

// SDKFunc represents the Deposit SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Deposit
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Deposit {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createDeposit(params.ID, params.To, params.Token, params.Amount)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if deposit, ok := ins.(Deposit); ok {
					out := createStorableDeposit(deposit)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if deposit, ok := ins.(Deposit); ok {
					base := retrieveAllDepositsKeyname()

					return []string{
						base,
						retrieveDepositsByToWalletIDKeyname(deposit.To().ID()),
						retrieveDepositsByTokenIDKeyname(deposit.Token().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return nil, errors.New(str)

			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {

				// create the entity repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)
				service := entity.SDKFunc.CreateService(ds)

				walletRepresentation := wallet.SDKFunc.CreateRepresentation()
				tokRepresentation := token.SDKFunc.CreateRepresentation()

				if deposit, ok := ins.(Deposit); ok {
					// try to retrieve the wallet:
					toWallet := deposit.To()
					_, retToWalletErr := repository.RetrieveByID(walletRepresentation.MetaData(), toWallet.ID())
					if retToWalletErr != nil {
						// save the wallet:
						saveErr := service.Save(toWallet, walletRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}

					// try to retrieve the token:
					tok := deposit.Token()
					_, retTokErr := repository.RetrieveByID(tokRepresentation.MetaData(), tok.ID())
					if retTokErr != nil {
						// save the token:
						saveErr := service.Save(tok, tokRepresentation)
						if saveErr != nil {
							return saveErr
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Deposit instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		met := createMetaData()
		if params.Datastore != nil {
			entityRepository := entity.SDKFunc.CreateRepository(params.Datastore)
			out := createRepository(entityRepository, met)
			return out
		}

		out := createRepository(params.EntityRepository, met)
		return out
	},
}
