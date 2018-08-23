package tendermint

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	router "github.com/XMNBlockchain/datamint/router"
	uuid "github.com/satori/go.uuid"
)

func TestCreateBlockchain_thenSpawn_Success(t *testing.T) {

	//variables:
	namespace := "xmnsuite"
	name := "users"
	id := uuid.NewV4()
	rootPath := filepath.Join("./test_files")
	simpleApp := createSimpleTestApplication()

	//delete the files at the end of the tests:
	defer func(dirPath string) {
		os.RemoveAll(dirPath)
	}(filepath.Join(rootPath, namespace))

	//create the blockchain service:
	blkChainService := SDKFunc.CreateBlockchainService(CreateBlockchainServiceParams{
		RootDirPath: rootPath,
	})

	//generate the blockchain:
	blkChain := SDKFunc.CreateBlockchain(CreateBlockchainParams{
		ID:        &id,
		Namespace: namespace,
		Name:      name,
	})

	//save the blockchain:
	saveErr := blkChainService.Save(blkChain)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	//create the router service:
	routerService := SDKFunc.CreateRouterService(CreateRouterServiceParams{
		RootDir:  rootPath,
		BlkChain: blkChain,
		Router:   simpleApp,
	})

	//spawn the service:
	serv, client, spawnErr := routerService.Spawn()
	if spawnErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", spawnErr.Error())
		return
	}
	defer serv.Stop()
	defer client.Stop()

	//start the client:
	startErr := client.Start()
	if startErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", startErr.Error())
		return
	}

	//execute a transaction:
	trsResponse := client.Transact(router.SDKFunc.CreateRequest(
		router.CreateRequestParams{
			Path: "/some/resource",
			Data: []byte("works!"),
		},
	))

	fmt.Printf("->-> %s\n\n", trsResponse.Log())

	fmt.Printf("->-> %v\n", trsResponse)
}
