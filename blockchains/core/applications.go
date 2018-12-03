package core

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

func createApplications(
	namespace string,
	name string,
	id *uuid.UUID,
	rootDir string,
	routePrefix string,
	routerRoleKey string,
	rootPubKey crypto.PublicKey,
	ds datastore.StoredDataStore,
	met meta.Meta,
) applications.Applications {

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: []applications.Application{
			create20181106(namespace, name, id, 0, -1, rootDir, routePrefix, routerRoleKey, rootPubKey, ds, met),
		},
	})

	return apps
}
