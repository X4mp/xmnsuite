package core

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/datastore"
)

func createApplications(
	namespace string,
	name string,
	id *uuid.UUID,
	rootDir string,
	routerDS datastore.DataStore,
	routerRoleKey string,
	ds datastore.DataStore,
) []applications.Application {

	// create the first application:
	apps := []applications.Application{
		create20181106(namespace, name, id, 0, -1, rootDir, routerDS, routerRoleKey, ds),
	}

	return apps
}
