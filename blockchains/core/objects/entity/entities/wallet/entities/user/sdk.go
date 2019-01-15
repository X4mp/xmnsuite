package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// User represents a user
type User interface {
	ID() *uuid.UUID
	PubKey() crypto.PublicKey
	Name() string
	Shares() int
	Wallet() wallet.Wallet
	HasBeenReferred() bool
	Referral() wallet.Wallet
}

// Normalized represents a normalized user
type Normalized interface {
}

// Repository represents the user repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (User, error)
	RetrieveByName(name string) (User, error)
	RetrieveByPubKey(pubKey crypto.PublicKey) (User, error)
	RetrieveSetByWallet(wal wallet.Wallet, index int, amount int) (entity.PartialSet, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID       *uuid.UUID
	Name     string
	PubKey   crypto.PublicKey
	Shares   int
	Wallet   wallet.Wallet
	Referral wallet.Wallet
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	Store            datastore.DataStore
	EntityRepository entity.Repository
}

// SDKFunc represents the User SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) User
	CreateRepository     func(params CreateRepositoryParams) Repository
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) User {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		if params.Referral != nil {
			out, outErr := createUserWithReferral(params.ID, params.Name, params.PubKey, params.Shares, params.Wallet, params.Referral)
			if outErr != nil {
				panic(outErr)
			}

			return out
		}

		out, outErr := createUser(params.ID, params.Name, params.PubKey, params.Shares, params.Wallet)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		if params.Store != nil {
			params.EntityRepository = entity.SDKFunc.CreateRepository(params.Store)
		}

		userMetaData := createMetaData()
		out := createRepository(userMetaData, params.EntityRepository)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if usr, ok := ins.(User); ok {
					out := createStorableUser(usr)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if usr, ok := ins.(User); ok {
					keynames := []string{
						retrieveAllUserKeyname(),
						retrieveUserByPubKeyKeyname(usr.PubKey()),
						retrieveUserByWalletIDKeyname(usr.Wallet().ID()),
						retrieveUserByNameKeyname(usr.Name()),
					}

					if usr.HasBeenReferred() {
						keynames = append(keynames, retrieveUserByWalletReferralKeyname(usr.Referral().ID()))
					}

					return keynames, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
				return nil, errors.New(str)

			},
			OnSave: func(ds datastore.DataStore, ins entity.Entity) error {

				if usr, ok := ins.(User); ok {

					//create the metadata and representation:
					metaData := createMetaData()
					walRepresentation := wallet.SDKFunc.CreateRepresentation()

					// create the repositories and services:
					entityRepository := entity.SDKFunc.CreateRepository(ds)
					repository := createRepository(metaData, entityRepository)
					entityService := entity.SDKFunc.CreateService(ds)

					// make sure there is no other user with the given public key:
					_, retUserErr := repository.RetrieveByPubKey(usr.PubKey())
					if retUserErr == nil {
						str := fmt.Sprintf("the User instance (PubKey: %s) already exists", usr.PubKey().String())
						return errors.New(str)
					}

					// make sure there is no other user with the same name:
					_, retUserByNameErr := repository.RetrieveByName(usr.Name())
					if retUserByNameErr == nil {
						str := fmt.Sprintf("the User instance (Name: %s) already exists", usr.Name())
						return errors.New(str)
					}

					// if the wallet does not exists, create it:
					wal := usr.Wallet()
					_, retWalErr := entityRepository.RetrieveByID(walRepresentation.MetaData(), wal.ID())
					if retWalErr != nil {
						saveWalErr := entityService.Save(wal, walRepresentation)
						if saveWalErr != nil {
							return saveWalErr
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
}
