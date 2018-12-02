package core

import (
	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// CreateParams represents the Create params
type CreateParams struct {
	Namespace     string
	Name          string
	ID            *uuid.UUID
	Port          int
	NodePK        tcrypto.PrivKey
	RootDir       string
	RouterRoleKey string
	RootPubKey    crypto.PublicKey
	Store         datastore.StoredDataStore
}

var SDKFunc = struct {
	Create func(params CreateParams) applications.Applications
}{
	Create: func(params CreateParams) applications.Applications {
		apps := createApplications(params.Namespace, params.Name, params.ID, params.RootDir, params.RouterRoleKey, params.RootPubKey, params.Store)
		return apps
	},
}
