package core

import (
	"fmt"
	"net"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

func connectToBlockchain(ip net.IP, port int) (applications.Client, error) {
	// create the service:
	appService := tendermint.SDKFunc.CreateApplicationService()

	// create the address:
	address := fmt.Sprintf("tcp://%s:%d", ip.String(), port)

	// connect:
	cl, clErr := appService.Connect(address)
	if clErr != nil {
		return nil, clErr
	}

	return cl, nil
}

func spawnBlockchain(
	namespace string,
	name string,
	id *uuid.UUID,
	seeds []string,
	rootDirPath string,
	port int,
	met meta.Meta,
) (applications.Node, error) {
	service := tendermint.SDKFunc.CreateBlockchainService(tendermint.CreateBlockchainServiceParams{
		RootDirPath: rootDirPath,
	})

	blkchain, blkchainErr := service.Retrieve(tendermint.SDKFunc.CreatePath(tendermint.CreatePathParams{
		Namespace: namespace,
		Name:      name,
		ID:        id,
	}))

	if blkchainErr != nil {
		return nil, blkchainErr
	}

	// create the datastore:
	store := datastore.SDKFunc.CreateStoredDataStore(datastore.StoredDataStoreParams{
		FilePath: filepath.Join(rootDirPath, "db.xmn"),
	})

	// create the applications:
	routerRoleKey := "router-role"
	apps := createApplications(namespace, name, id, rootDirPath, routerRoleKey, store, met)

	// create the application service:
	appService := tendermint.SDKFunc.CreateApplicationService()

	// spawn the node:
	node, nodeErr := appService.Spawn(port, seeds, rootDirPath, blkchain, apps)
	if nodeErr != nil {
		return nil, nodeErr
	}

	return node, nil
}

func saveThenSpawnBlockchain(
	namespace string,
	name string,
	id *uuid.UUID,
	seeds []string,
	rootDirPath string,
	port int,
	pk tcrypto.PrivKey,
	rootPubKey crypto.PublicKey,
	met meta.Meta,
) (applications.Node, error) {
	blkchain := tendermint.SDKFunc.CreateBlockchain(tendermint.CreateBlockchainParams{
		Namespace: namespace,
		Name:      name,
		ID:        id,
		PrivKey:   pk,
	})

	service := tendermint.SDKFunc.CreateBlockchainService(tendermint.CreateBlockchainServiceParams{
		RootDirPath: rootDirPath,
	})

	saveErr := service.Save(blkchain)
	if saveErr != nil {
		return nil, saveErr
	}

	// create the datastore:
	store := datastore.SDKFunc.CreateStoredDataStore(datastore.StoredDataStoreParams{
		FilePath: filepath.Join(rootDirPath, "db.xmn"),
	})

	// create the applications:
	routerRoleKey := "router-role"
	apps := createApplicationsWithRootPubKey(namespace, name, id, rootDirPath, routerRoleKey, store, met, rootPubKey)

	// create the application service:
	appService := tendermint.SDKFunc.CreateApplicationService()

	// spawn the node:
	node, nodeErr := appService.Spawn(port, seeds, rootDirPath, blkchain, apps)
	if nodeErr != nil {
		return nil, nodeErr
	}

	return node, nil
}
