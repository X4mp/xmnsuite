package wallet

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Wallet represents a wallet
type Wallet interface {
	ID() *uuid.UUID
	Creator() crypto.PublicKey
	ConcensusNeeded() int
}

// Repository represents the wallet repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Wallet, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByCreatorPublicKey(pubKey crypto.PublicKey, index int, amount int) (entity.PartialSet, error)
}

// Normalized represents a normalized wallet
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID              *uuid.UUID
	Creator         crypto.PublicKey
	ConcensusNeeded int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Wallet SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Wallet
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Wallet {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createWallet(params.ID, params.Creator, params.ConcensusNeeded)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met:        createMetaData(),
			ToStorable: toData,
			Keynames: func(ins entity.Entity) ([]string, error) {
				if wal, ok := ins.(Wallet); ok {
					return []string{
						retrieveAllWalletKeyname(),
						retrieveByPublicKeyWalletKeyname(wal.Creator()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", ins.ID().String())
				return nil, errors.New(str)

			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)

				// create the metadata:
				metaData := createMetaData()

				if wal, ok := ins.(Wallet); ok {
					// if the wallet already exists:
					retWalletIns, retWalletInsErr := repository.RetrieveByID(metaData, wal.ID())
					if retWalletInsErr == nil {
						// cast the returned wallet:
						if retWal, ok := retWalletIns.(Wallet); ok {
							// make sure the creator is the same:
							if !retWal.Creator().Equals(wal.Creator()) {
								str := fmt.Sprintf("the Wallet (ID: %s) already existed but the creator pubKey does not match.  Expected: %s, Received: %s", retWal.ID().String(), wal.Creator().String(), retWal.Creator().String())
								return errors.New(str)
							}

							// everything is fine, it will update the wallet:
							return nil
						}
					}

					// the wallet doesnt exists, so everything is fine:
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		rep := createRepository(metaData, params.EntityRepository)
		return rep
	},
}
