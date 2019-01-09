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
	routerRoleKey string,
	ds datastore.StoredDataStore,
	met meta.Meta,
) applications.Applications {

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: []applications.Application{
			create20181106(namespace, name, id, 0, -1, rootDir, routerRoleKey, ds, met),
		},
	})

	return apps
}

func createApplicationsWithRootPubKey(
	namespace string,
	name string,
	id *uuid.UUID,
	rootDir string,
	routerRoleKey string,
	ds datastore.StoredDataStore,
	met meta.Meta,
	rootPubKey crypto.PublicKey,
) applications.Applications {

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: []applications.Application{
			create20181106WithRootPubKey(namespace, name, id, 0, -1, rootDir, routerRoleKey, ds, met, rootPubKey),
		},
	})

	return apps
}
