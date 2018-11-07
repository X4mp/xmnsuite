package entity

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

type storableHuman struct {
	UUID   string `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
}

type human struct {
	UUID   *uuid.UUID `json:"id"`
	Name   string     `json:"name"`
	Height int        `json:"height"`
}

// ID returns the ID
func (obj *human) ID() *uuid.UUID {
	return obj.UUID
}

func createHumanForTests() Entity {
	names := []string{
		"Steve",
		"Roger",
		"John",
	}

	id := uuid.NewV4()
	out := human{
		UUID:   &id,
		Name:   names[rand.Int()%len(names)],
		Height: (rand.Int() % 100) + 100,
	}

	return &out
}

func startBlockchain(
	t *testing.T,
	app applications.Application,
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
			app,
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

func createApplicationForTests(
	routerRoleKey string,
	ctrls Controllers,
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

	// create the filepath:
	fileName := fmt.Sprintf("%s.%s", version, "xmndb")
	filePath := filepath.Join(rootPath, namespace, name, id.String(), "application", fileName)

	// create the stored ds:
	st := datastore.SDKFunc.CreateStoredDataStore(datastore.StoredDataStoreParams{
		FilePath: filePath,
	})

	// enable our user to write on the right routes:
	st.DataStore().Users().Insert(privKey.PublicKey())
	st.DataStore().Roles().Add(routerRoleKey, privKey.PublicKey())
	st.DataStore().Roles().EnableWriteAccess(routerRoleKey, "/")

	// create application:
	return applications.SDKFunc.CreateApplication(applications.CreateApplicationParams{
		Namespace:      namespace,
		Name:           name,
		ID:             id,
		FromBlockIndex: fromBlockIndex,
		ToBlockIndex:   toBlockIndex,
		Version:        version,
		DirPath:        rootPath,
		Store:          st,
		RouterParams: routers.CreateRouterParams{
			DataStore: st.DataStore(),
			RoleKey:   routerRoleKey,
			RtesParams: []routers.CreateRouteParams{
				routers.CreateRouteParams{
					Pattern: "/",
					SaveTrx: ctrls.Save(),
				},
				routers.CreateRouteParams{
					Pattern: "/id:<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
					DelTrx:  ctrls.Delete(),
				},
				routers.CreateRouteParams{
					Pattern:  "/id:<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
					QueryTrx: ctrls.RetrieveByID(),
				},
				routers.CreateRouteParams{
					Pattern:  "/keynames:<keynames|[a-z-,]+>",
					QueryTrx: ctrls.RetrieveByIntersectKeynames(),
				},
				routers.CreateRouteParams{
					Pattern:  "/set/keynames:<keynames|[a-z-,]+>",
					QueryTrx: ctrls.RetrieveSetByIntersectKeynames(),
				},
				routers.CreateRouteParams{
					Pattern:  "/set/keyname:<keyname|[a-z]+>",
					QueryTrx: ctrls.RetrieveSetByKeyname(),
				},
			},
		},
	})
}

func executeSaveTransaction(
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

func executeDeleteTransaction(
	t *testing.T,
	client applications.Client,
	path string,
	privKey crypto.PrivateKey,
) applications.ClientTransactionResponse {
	// create the resource:
	resPtr := routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
		From: privKey.PublicKey(),
		Path: path,
	})

	// sign the resource:
	sig := privKey.Sign(resPtr.Hash())

	// save the instance:
	trxResp, trxRespErr := client.Transact(routers.SDKFunc.CreateTransactionRequest(routers.CreateTransactionRequestParams{
		Ptr: resPtr,
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

func verifyClientTransactionResponseIsSuccessful(t *testing.T, path string, ins Entity, trxResp applications.ClientTransactionResponse, gazPricePerKb int) {
	// convert to json:
	js, jsErr := cdc.MarshalJSON(ins)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return
	}

	// verify the code:
	retTrxCode := trxResp.Transaction().Code()
	if retTrxCode != routers.IsSuccessful {
		t.Errorf("the transaction was expected to be successful, the check log was: %s, the transaction log was: %s", trxResp.Check().Log(), trxResp.Transaction().Log())
		return
	}

	// verify the gaz used:
	retTrxGazUsed := trxResp.Transaction().GazUsed()
	gazUsed := int(unsafe.Sizeof(js)) * gazPricePerKb
	if retTrxGazUsed != int64(gazUsed) {
		t.Errorf("the returned gaz used was expected to be %d, returned: %d, the check log was: %s, the transaction log was: %s", int64(gazUsed), retTrxGazUsed, trxResp.Check().Log(), trxResp.Transaction().Log())
		return
	}

	// verify the log:
	retTrxLog := trxResp.Transaction().Log()
	if retTrxLog != "success" {
		t.Errorf("the returned log was expected to be: %s, returned: %s, the check log was: %s, the transaction log was: %s", "success", retTrxLog, trxResp.Check().Log(), trxResp.Transaction().Log())
		return
	}

	// verify the tags:
	retTrxTags := trxResp.Transaction().Tags()
	elementPath := fmt.Sprintf("%sid:%s", path, ins.ID().String())
	expectedTags := map[string][]byte{elementPath: js}
	if !reflect.DeepEqual(retTrxTags, expectedTags) {
		t.Errorf("the returned tags are invalid, the check log was: %s, the transaction log was: %s", trxResp.Check().Log(), trxResp.Transaction().Log())
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

func saveIsSuccessful(t *testing.T, routePath string, ins Entity, storableIns interface{}, privKey crypto.PrivateKey, client applications.Client, gazPricePerKb int) {
	// execute the transaction:
	trxResp := executeSaveTransaction(t, client, storableIns, routePath, privKey)
	if trxResp == nil {
		return
	}

	// verify that the transaction was successful:
	verifyClientTransactionResponseIsSuccessful(t, routePath, ins, trxResp, gazPricePerKb)
}

func TestSave_thenRetrieve_thenRetrieveByID_thenDelete_Success(t *testing.T) {
	// variables:
	routerRoleKey := "router-role-key"
	privKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
	namespace := "xmn"
	name := "human"
	id := uuid.NewV4()
	gazPricePerKb := rand.Int() % 20
	rootPath := filepath.Join("./test_files")
	defer func() {
		os.RemoveAll(rootPath)
	}()

	met := SDKFunc.CreateMetaData(CreateMetaDataParams{
		Name: "Human",
		ToEntity: func(rep Repository, data interface{}) (Entity, error) {
			fromStorableToHuman := func(storable *storableHuman) (*human, error) {
				insID, insIDErr := uuid.FromString(storable.UUID)
				if insIDErr != nil {
					return nil, insIDErr
				}

				return &human{
					UUID:   &insID,
					Name:   storable.Name,
					Height: storable.Height,
				}, nil
			}

			if storable, ok := data.(*storableHuman); ok {
				return fromStorableToHuman(storable)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(storableHuman)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return fromStorableToHuman(ptr)
			}

			return nil, errors.New("the given data is invalid and therefore cannot be converted to an Entity instance")
		},
		EmptyStorable: new(storableHuman),
	})

	rep := SDKFunc.CreateRepresentation(CreateRepresentationParams{
		Met: met,
		ToStorable: func(ins Entity) (interface{}, error) {
			if obj, ok := ins.(*human); ok {
				return &storableHuman{
					UUID:   obj.ID().String(),
					Name:   obj.Name,
					Height: obj.Height,
				}, nil
			}

			return nil, errors.New("the given Entity instance (ID: %s) cannot be converted to storable instance")

		},
		Keynames: func(ins Entity) ([]string, error) {
			sh := ins.(*human)
			return []string{
				fmt.Sprintf("human:by_height:%d", sh.Height),
			}, nil

		},
	})

	// create the controllers:
	controllers := SDKFunc.CreateControllers(CreateControllersParams{
		Met: met,
		Rep: rep,
		DefaultAmountOfElements:  20,
		GazPricePerKb:            gazPricePerKb,
		OverwriteIfAlreadyExists: false,
		RouterRoleKey:            routerRoleKey,
	})

	// start the blockchain:
	app := createApplicationForTests(routerRoleKey, controllers, namespace, name, &id, rootPath, privKey)
	node, client := startBlockchain(t, app, namespace, name, &id, rootPath)
	defer node.Stop()

	// create 2 entities:
	first := createHumanForTests()
	firstStorable, _ := rep.ToStorable()(first)
	second := createHumanForTests()
	secondStorable, _ := rep.ToStorable()(second)

	// route paths:
	retrieveByFirstID := fmt.Sprintf("/id:%s", first.ID().String())
	retrieveBySecondID := fmt.Sprintf("/id:%s", second.ID().String())

	// retrieve the first entity by ID before its inserted, returns error:
	firstQueryRespWithError := executeQuery(t, privKey, retrieveByFirstID, client)
	if firstQueryRespWithError.Code() != routers.InvalidRequest {
		t.Errorf("the entity instance retrieval was expected to be invalid")
		return
	}

	// delete the first entity, before its inserted:
	firstDelTrxRespWithError := executeDeleteTransaction(t, client, retrieveByFirstID, privKey)
	if firstDelTrxRespWithError == nil {
		return
	}

	// verify delete, returns error:
	if firstDelTrxRespWithError.Check().Code() == routers.IsSuccessful {
		t.Errorf("the returned transaction was not expected to succeed, log: %s", firstDelTrxRespWithError.Check().Log())
		return
	}

	// save the first entity instance successfully:
	saveIsSuccessful(t, "/", first, firstStorable, privKey, client, gazPricePerKb)

	// delete the first entity, success:
	firstDelTrxResp := executeDeleteTransaction(t, client, retrieveByFirstID, privKey)
	if firstDelTrxResp == nil {
		return
	}

	// verify delete, success:
	if firstDelTrxResp.Check().Code() != routers.IsSuccessful {
		t.Errorf("the returned transaction was expected to succeed, log: %s", firstDelTrxResp.Check().Log())
		return
	}

	// retrieve the first entity by ID, success:
	firstQueryResp := executeQuery(t, privKey, retrieveByFirstID, client)
	if firstQueryResp == nil {
		return
	}

	// retrieve the second entity by ID before its inserted, returns error:
	secondQueryRespWithError := executeQuery(t, privKey, retrieveBySecondID, client)
	if secondQueryRespWithError.Code() != routers.InvalidRequest {
		t.Errorf("the entity instance retrieval was expected to be invalid")
		return
	}

	// save the second entity instance successfully:
	saveIsSuccessful(t, "/", second, secondStorable, privKey, client, gazPricePerKb)

	// retrieve the first entity by ID, success:
	secondQueryResp := executeQuery(t, privKey, retrieveBySecondID, client)
	if secondQueryResp == nil {
		return
	}

}
