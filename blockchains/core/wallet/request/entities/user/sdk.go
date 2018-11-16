package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// User represents a user
type User interface {
	ID() *uuid.UUID
	PubKey() crypto.PublicKey
	Shares() int
	Wallet() wallet.Wallet
}

// Normalized represents a normalized user
type Normalized interface {
}

// Repository represents the user repository
type Repository interface {
	RetrieveByPubKeyAndWallet(pubKey crypto.PublicKey, wal wallet.Wallet) (User, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	PubKey crypto.PublicKey
	Shares int
	Wallet wallet.Wallet
}

// SDKFunc represents the User SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) User
	CreateRepository     func(ds datastore.DataStore) Repository
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) User {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createUser(params.ID, params.PubKey, params.Shares, params.Wallet)
		return out
	},
	CreateRepository: func(ds datastore.DataStore) Repository {
		entityRepository := entity.SDKFunc.CreateRepository(ds)
		userMetaData := createMetaData()
		out := createRepository(userMetaData, entityRepository)
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
					return []string{
						retrieveAllUserKeyname(),
						retrieveUserByPubKeyKeyname(usr.PubKey()),
						retrieveUserByWalletIDKeyname(usr.Wallet().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
				return nil, errors.New(str)

			},
		})
	},
}
