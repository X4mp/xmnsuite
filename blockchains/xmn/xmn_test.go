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
	routerDS := datastore.SDKFunc.Create()

	// enable our user to write on the right routes:
	routerDS.Users().Insert(privKey.PublicKey())
	routerDS.Roles().Add(routerRoleKey, privKey.PublicKey())
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/")
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/wallets")
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/wallets/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>")
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/user-requests")
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/user-request-votes")
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/users/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>")

	// create application:
	xmn := createXMN(namespace, name, id, fromBlockIndex, toBlockIndex, version, rootPath, routerDS, routerRoleKey)
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

	// save the instance:
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

func verifyClientTransactionResponseIsSuccessful(t *testing.T, path string, ins interface{}, trxResp applications.ClientTransactionResponse, gazPricePerKb int) {
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
	expectedTags := map[string][]byte{fmt.Sprintf(path): js}
	if !reflect.DeepEqual(retTrxTags, expectedTags) {
		t.Errorf("the returned tags are invalid")
		return
	}
}

func verifyQueryResponseIsSuccessful(t *testing.T, empty interface{}, queryPath string, queryResp routers.QueryResponse, ins interface{}) {
	// convert to json:
	js, jsErr := cdc.MarshalJSON(ins)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return
	}

	retQueryCode := queryResp.Code()
	if retQueryCode != routers.IsSuccessful {
		t.Errorf("the query was expected to be successful: %s", queryResp.Log())
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
	unJsErr := cdc.UnmarshalJSON(retQueryValue, empty)
	if unJsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", unJsErr.Error())
		return
	}

	if !reflect.DeepEqual(retQueryValue, js) {
		t.Errorf("the returned instance is invalid")
		return
	}
}

func saveGenesisIsSuccessful(t *testing.T, gen Genesis, privKey crypto.PrivateKey, client applications.Client) {
	// retrieve the wallets, should have 0 wallet:
	queryWithZeroWallets := executeQuery(t, privKey, "/wallets", client)
	if queryWithZeroWallets == nil {
		return
	}

	zeoWallets := new(walletPartialSet)
	zeroWalletsErr := cdc.UnmarshalJSON(queryWithZeroWallets.Value(), zeoWallets)
	if zeroWalletsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", zeroWalletsErr.Error())
		return
	}

	if zeoWallets.TotalAmount() != 0 {
		t.Errorf("the total amount was expected to be 0, %d returned", zeoWallets.TotalAmount())
		return
	}

	// execute the query, expects an error:
	queryRespWithError := executeQuery(t, privKey, "/", client)
	if queryRespWithError.Code() != routers.InvalidRequest {
		t.Errorf("the genesis instance retrieval was expected to be invalid")
		return
	}

	// execute the transaction:
	trxResp := executeTransaction(t, client, gen, "/", privKey)
	if trxResp == nil {
		return
	}

	// verify that the transaction was successful:
	verifyClientTransactionResponseIsSuccessful(t, "/", gen, trxResp, gen.GazPricePerKb())

	// retrieve the wallets, should have the genesis wallets:
	query := executeQuery(t, privKey, "/wallets", client)
	if query == nil {
		return
	}

	wallets := new(walletPartialSet)
	walletsErr := cdc.UnmarshalJSON(query.Value(), wallets)
	if walletsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", walletsErr.Error())
		return
	}

	if wallets.TotalAmount() != 1 {
		t.Errorf("the total amount was expected to be 1, %d returned", wallets.TotalAmount())
		return
	}
}

func TestXMN_Genesis_Success(t *testing.T) {
	//variables:
	namespace := "xmn"
	name := "core"
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
	secondGen := createGenesisForTests()

	// save the first genesis instance successfully:
	saveGenesisIsSuccessful(t, firstGen, privKey, client)

	// execute the query:
	queryResp := executeQuery(t, privKey, "/", client)
	if queryResp == nil {
		return
	}

	// make sure the query was successful:
	verifyQueryResponseIsSuccessful(t, new(genesis), "/", queryResp, firstGen)

	// execute the transaction again, expects an error:
	trxRespWithError := executeTransaction(t, client, secondGen, "/", privKey)
	if trxRespWithError.Check().Code() != routers.InvalidRequest {
		t.Errorf("the genesis transaction was expected to fail because it already exists")
		return
	}
}

func TestXMN_wallet_Success(t *testing.T) {
	//variables:
	namespace := "xmn"
	name := "core"
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

	// create the wallets path:
	walletPath := "/wallets"

	// create the genesis:
	gen := createGenesisForTests()

	// save the first genesis instance successfully:
	saveGenesisIsSuccessful(t, gen, privKey, client)

	// create wallets:
	firstWallet := createWalletForTests()

	// save the first wallet:
	firstTrxResp := executeTransaction(t, client, firstWallet, walletPath, privKey)
	if firstTrxResp == nil {
		return
	}

	// verify that the transaction was successful:
	verifyClientTransactionResponseIsSuccessful(t, walletPath, firstWallet, firstTrxResp, gen.GazPricePerKb())

	// execute the query:
	firstWalletByIDPath := filepath.Join(walletPath, firstWallet.ID().String())
	firstQueryResp := executeQuery(t, privKey, firstWalletByIDPath, client)
	if firstQueryResp == nil {
		return
	}

	// make sure the query was successful:
	verifyQueryResponseIsSuccessful(t, new(wallet), firstWalletByIDPath, firstQueryResp, firstWallet)

	// retrieve the wallets, should have the new wallet in the list:
	queryAllRep := executeQuery(t, privKey, walletPath, client)
	if queryAllRep == nil {
		return
	}

	wallets := new(walletPartialSet)
	walletsErr := cdc.UnmarshalJSON(queryAllRep.Value(), wallets)
	if walletsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", walletsErr.Error())
		return
	}

	if wallets.TotalAmount() != 2 {
		t.Errorf("the total amount was expected to be 2, %d returned", wallets.TotalAmount())
		return
	}
}

func TestXMN_user_Success(t *testing.T) {
	//variables:
	namespace := "xmn"
	name := "core"
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

	// create the paths:
	walletPath := "/wallets"
	userRequestPath := "/user-requests"
	userRequestVotePath := "/user-request-votes"
	userPath := "/users"

	// create the genesis:
	shares := 20
	concensusNeeded := 20
	gen := createGenesisWithSharesAndConcensusForTests(shares, concensusNeeded)

	// save the first genesis instance successfully:
	saveGenesisIsSuccessful(t, gen, privKey, client)

	// create wallets:
	firstWallet := createWalletForTests()

	// save the first wallet:
	saveWalletTrxResp := executeTransaction(t, client, firstWallet, walletPath, privKey)
	if saveWalletTrxResp == nil {
		return
	}

	// verify that the transaction was successful:
	verifyClientTransactionResponseIsSuccessful(t, walletPath, firstWallet, saveWalletTrxResp, gen.GazPricePerKb())

	// create user request:
	requestUsr := createUserWithWalletForTests(gen.Deposit().To().Wallet())
	userReq := createUserRequest(requestUsr)
	storedUserReq := createStoredUserRequest(userReq.User())

	// save user request:
	saveUsrReqTrxResp := executeTransaction(t, client, storedUserReq, userRequestPath, privKey)
	if saveUsrReqTrxResp == nil {
		return
	}

	// verify that the transaction was successful:
	verifyClientTransactionResponseIsSuccessful(t, userRequestPath, userReq, saveUsrReqTrxResp, gen.GazPricePerKb())

	// create first user request vote:
	firstVote := createUserRequestVoteWithUserRequestAndVoterForTests(userReq, gen.Deposit().To(), true)
	storedFirstVote := createStoredUserRequestVote(firstVote)

	// save first vote:
	saveVoteTrxResp := executeTransaction(t, client, storedFirstVote, userRequestVotePath, privKey)
	if saveVoteTrxResp == nil {
		return
	}

	// verify that the transaction was successful:
	verifyClientTransactionResponseIsSuccessful(t, userRequestVotePath, firstVote, saveVoteTrxResp, gen.GazPricePerKb())

	// retrieve the user by ID:
	userByIDPath := fmt.Sprintf("%s/%s", userPath, firstVote.Request().User().ID().String())
	queryUserResp := executeQuery(t, privKey, userByIDPath, client)
	if queryUserResp == nil {
		return
	}

	ptr := new(user)
	unUsrErr := cdc.UnmarshalJSON(queryUserResp.Value(), ptr)
	if unUsrErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", unUsrErr.Error())
		return
	}

	// compare the users:
	compareUserForTests(t, firstVote.Request().User(), ptr)
}
