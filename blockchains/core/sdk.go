package core

import (
	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// CreateParams represents the Create params
type CreateParams struct {
	Namespace                     string
	Name                          string
	ID                            *uuid.UUID
	Port                          int
	NodePK                        tcrypto.PrivKey
	RootDir                       string
	RoutePrefix                   string
	RouterRoleKey                 string
	Store                         datastore.StoredDataStore
	Meta                          meta.Meta
	RootPubKey                    crypto.PublicKey
	MaxAmountOfEntitiesToRetrieve int
}

// SDKFunc represents the core SDK func
var SDKFunc = struct {
	Create func(params CreateParams) applications.Applications
}{
	Create: func(params CreateParams) applications.Applications {

		if params.RootPubKey == nil {
			return createApplications(
				params.Namespace,
				params.Name,
				params.ID,
				params.RootDir,
				params.RoutePrefix,
				params.RouterRoleKey,
				params.Store,
				params.Meta,
				params.MaxAmountOfEntitiesToRetrieve,
			)
		}

		return createApplicationsWithRootPubKey(
			params.Namespace,
			params.Name,
			params.ID,
			params.RootDir,
			params.RoutePrefix,
			params.RouterRoleKey,
			params.Store,
			params.Meta,
			params.RootPubKey,
			params.MaxAmountOfEntitiesToRetrieve,
		)
	},
}
