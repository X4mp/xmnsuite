package xmn

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

func createXMNForTests(
	namespace string,
	name string,
	id *uuid.UUID,
	rootPath string,
	privKey crypto.PrivateKey,
) applications.Application {
	//variables:
	fromBlockIndex := int64(0)
	toBlockIndex := int64(-1)
	version := "2018.10.15"
	routerRoleKey := "router-role-key"

	// datastore:
	store := datastore.SDKFunc.Create()
	routerDS := datastore.SDKFunc.Create()

	// enable our user to write on the right routes:
	routerDS.Users().Insert(privKey.PublicKey())
	routerDS.Roles().Add(routerRoleKey, privKey.PublicKey())
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/")

	// services:
	walletService := createWalletService(store)
	genService := createGenesisService(store, walletService)

	// create application:
	xmn := createXMN(genService, namespace, name, id, fromBlockIndex, toBlockIndex, version, rootPath, routerDS, routerRoleKey)
	return xmn
}

func startBlockchain(
	t *testing.T,
	xmn applications.Application,
	namespace string,
	name string,
	id *uuid.UUID,
	rootPath string,
) (applications.Node, applications.Client) {
	// variables:
	nodePrivKey := ed25519.GenPrivKey()
	port := rand.Int()%9000 + 1000

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: []applications.Application{
			xmn,
		},
	})

	// create the blockchain:
	blkChain := tendermint.SDKFunc.CreateBlockchain(tendermint.CreateBlockchainParams{
		Namespace: namespace,
		Name:      name,
		ID:        id,
		PrivKey:   nodePrivKey,
	})

	// create the blockchain service:
	blkChainService := tendermint.SDKFunc.CreateBlockchainService(tendermint.CreateBlockchainServiceParams{
		RootDirPath: rootPath,
	})

	// save the blockchain:
	saveBlkChainErr := blkChainService.Save(blkChain)
	if saveBlkChainErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveBlkChainErr.Error())
		return nil, nil
	}

	// create the application service:
	appService := tendermint.SDKFunc.CreateApplicationService()

	// spawn the node:
	node, nodeErr := appService.Spawn(port, rootPath, blkChain, apps)
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return nil, nil
	}

	address := node.GetAddress()
	client, clientErr := appService.Connect(address)
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return nil, nil
	}

	// start the node:
	startNodeErr := node.Start()
	if startNodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", startNodeErr.Error())
		return nil, nil
	}

	return node, client
}

func executeTransaction(
	t *testing.T,
	client applications.Client,
	ins interface{},
	path string,
	privKey crypto.PrivateKey,
) applications.ClientTransactionResponse {
	// convert to json:
	js, jsErr := cdc.MarshalJSON(ins)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return nil
	}

	// create the resource:
	res := routers.SDKFunc.CreateResource(routers.CreateResourceParams{
		ResPtr: routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
			From: privKey.PublicKey(),
			Path: path,
		}),
		Data: js,
	})

	// sign the resource:
	sig := privKey.Sign(res.Hash())

	// save the message:
	trxResp, trxRespErr := client.Transact(routers.SDKFunc.CreateTransactionRequest(routers.CreateTransactionRequestParams{
		Res: res,
		Sig: sig,
	}))

	if trxRespErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", trxRespErr.Error())
		return nil
	}

	return trxResp
}

func executeQuery(t *testing.T, privKey crypto.PrivateKey, queryPath string, client applications.Client) routers.QueryResponse {
	// create the resource pointer:
	queryResPtr := routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
		From: privKey.PublicKey(),
		Path: queryPath,
	})

	// create the signature:
	querySig := privKey.Sign(queryResPtr.Hash())

	// execute a query:
	queryResp, queryRespErr := client.Query(routers.SDKFunc.CreateQueryRequest(routers.CreateQueryRequestParams{
		Ptr: queryResPtr,
		Sig: querySig,
	}))

	if queryRespErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", queryRespErr.Error())
		return nil
	}

	return queryResp
}

func verifyClientTransactionResponseIsSuccessful(t *testing.T, ins interface{}, trxResp applications.ClientTransactionResponse, gazPricePerKb int) {
	// convert to json:
	js, jsErr := cdc.MarshalJSON(ins)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return
	}

	// verify the code:
	retTrxCode := trxResp.Transaction().Code()
	if retTrxCode != routers.IsSuccessful {
		t.Errorf("the transaction was expected to be successful")
		return
	}

	// verify the gaz used:
	retTrxGazUsed := trxResp.Transaction().GazUsed()
	gazUsed := int(unsafe.Sizeof(js)) * gazPricePerKb
	if retTrxGazUsed != int64(gazUsed) {
		t.Errorf("the returned gaz used was expected to be %d, returned: %d", int64(gazUsed), retTrxGazUsed)
		return
	}

	// verify the log:
	retTrxLog := trxResp.Transaction().Log()
	if retTrxLog != "success" {
		t.Errorf("the returned log was expected to be: %s, returned: %s", "success", retTrxLog)
		return
	}

	// verify the tags:
	retTrxTags := trxResp.Transaction().Tags()
	expectedTags := map[string][]byte{fmt.Sprintf("/"): js}
	if !reflect.DeepEqual(retTrxTags, expectedTags) {
		t.Errorf("the returned tags are invalid")
		return
	}
}

func verifyQueryResponseIsSuccessful(t *testing.T, queryPath string, queryResp routers.QueryResponse, ins interface{}) {
	// convert to json:
	js, jsErr := cdc.MarshalJSON(ins)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return
	}

	retQueryCode := queryResp.Code()
	if retQueryCode != routers.IsSuccessful {
		t.Errorf("the query was expected to be successful")
		return
	}

	retQueryKey := queryResp.Key()
	if retQueryKey != queryPath {
		t.Errorf("the returned key was expected to be: %s, returned: %s", queryPath, retQueryKey)
		return
	}

	retQueryLog := queryResp.Log()
	if retQueryLog != "success" {
		t.Errorf("the returned log was expected to be: %s, returned: %s", "success", retQueryLog)
		return
	}

	retQueryValue := queryResp.Value()
	queryGenesis := new(genesis)
	unJsErr := cdc.UnmarshalJSON(retQueryValue, queryGenesis)
	if unJsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", unJsErr.Error())
		return
	}

	if !reflect.DeepEqual(retQueryValue, js) {
		t.Errorf("the returned instance is invalid")
		return
	}
}

func TestXMN_Genesis_Success(t *testing.T) {
	//variables:
	namespace := "xsuite"
	name := "users"
	id := uuid.NewV4()
	rootPath := filepath.Join("./test_files")
	privKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	defer func() {
		os.RemoveAll(rootPath)
	}()

	// create application:
	xmn := createXMNForTests(namespace, name, &id, rootPath, privKey)

	// start the blockchain:
	node, client := startBlockchain(t, xmn, namespace, name, &id, rootPath)
	defer node.Stop()

	// create the genesis:
	firstGen := createGenesisForTests()

	// execute the transaction:
	trxResp := executeTransaction(t, client, firstGen, "/", privKey)

	// verify that the transaction was successful:
	verifyClientTransactionResponseIsSuccessful(t, firstGen, trxResp, firstGen.GazPricePerKb())

	// execute the query:
	queryResp := executeQuery(t, privKey, "/", client)

	// make sure the query was successful:
	verifyQueryResponseIsSuccessful(t, "/", queryResp, firstGen)
}
