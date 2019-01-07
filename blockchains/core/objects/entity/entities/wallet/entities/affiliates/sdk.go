package affiliates

import (
	"net/url"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

// Affiliate represents an affiliate that hosts a web service that creates new users (and refer them)
type Affiliate interface {
	ID() *uuid.UUID
	Owner() wallet.Wallet
	URL() *url.URL
}

// Normalized represents a normalized host
type Normalized interface {
}

// Repository represents the affiliate repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Affiliate, error)
	RetrieveByWallet(wal wallet.Wallet) (Affiliate, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID    *uuid.UUID
	Owner wallet.Wallet
	URL   *url.URL
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// SDKFunc represents the Affiliate SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Affiliate
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
}{
	Create: func(params CreateParams) Affiliate {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createAffiliate(params.ID, params.Owner, params.URL)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return createRepresentation()
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(metaData, params.EntityRepository)
		return out
	},
}
