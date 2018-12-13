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
	ds datastore.StoredDataStore,
	met meta.Meta,
	maxAmountOfEntitiesToRetrieve int,
) applications.Applications {

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: []applications.Application{
			create20181106(namespace, name, id, 0, -1, rootDir, routePrefix, routerRoleKey, ds, met, maxAmountOfEntitiesToRetrieve),
		},
	})

	return apps
}

func createApplicationsWithRootPubKey(
	namespace string,
	name string,
	id *uuid.UUID,
	rootDir string,
	routePrefix string,
	routerRoleKey string,
	ds datastore.StoredDataStore,
	met meta.Meta,
	rootPubKey crypto.PublicKey,
	maxAmountOfEntitiesToRetrieve int,
) applications.Applications {

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: []applications.Application{
			create20181106WithRootPubKey(namespace, name, id, 0, -1, rootDir, routePrefix, routerRoleKey, ds, met, rootPubKey, maxAmountOfEntitiesToRetrieve),
		},
	})

	return apps
}
