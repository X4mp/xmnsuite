package tendermint

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	applications "github.com/xmnservices/xmnsuite/blockchains/applications"
	crypto "github.com/xmnservices/xmnsuite/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	objects "github.com/xmnservices/xmnsuite/objects"
	"github.com/xmnservices/xmnsuite/routers"
)

type messageForTest struct {
	ID          *uuid.UUID `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

func TestCreateBlockchainWithApplication_thenSpawn_Success(t *testing.T) {
	//variables:
	rootDir := "./test_files"
	defer func() {
		os.RemoveAll(rootDir)
	}()

	port := rand.Int()%9000 + 1000
	namespace := "testapp"
	name := "MyTestApp"
	id := uuid.NewV4()
	version := "2018.04.29"
	privKey := ed25519.GenPrivKey()
	store := datastore.SDKFunc.Create()
	fromPrivKey := crypto.SDKFunc.GenPK()
	fromPubKey := fromPrivKey.PublicKey()

	// enable our user to write on the right routes:
	routerRoleKey := "router-role-key"
	routerDS := datastore.SDKFunc.Create()
	routerDS.Users().Insert(fromPubKey)
	routerDS.Roles().Add(routerRoleKey, fromPubKey)
	routerDS.Roles().EnableWriteAccess(routerRoleKey, "/messages")

	// create application:
	app := applications.SDKFunc.CreateApplication(applications.CreateApplicationParams{
		FromBlockIndex: 0,
		ToBlockIndex:   -1,
		Version:        version,
		DataStore:      store,
		RouterParams: routers.CreateRouterParams{
			DataStore: routerDS,
			RoleKey:   routerRoleKey,
			RtesParams: []routers.CreateRouteParams{
				routers.CreateRouteParams{
					Pattern: "/messages",
					SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

						// unmarshal data:
						msg := new(messageForTest)
						jsErr := cdc.UnmarshalJSON(data, msg)
						if jsErr != nil {
							return nil, jsErr
						}

						// create the message path:
						msgPath := filepath.Join(path, msg.ID.String())

						// save the object:
						amountSaved := store.Objects().Save(&objects.ObjInKey{
							Key: msgPath,
							Obj: msg,
						})

						if amountSaved != 1 {
							str := fmt.Sprintf("there was an error while saving the message")
							return nil, errors.New(str)
						}

						resp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
							Code:    routers.IsSuccessful,
							Log:     "success",
							GazUsed: 1205,
							Tags: map[string][]byte{
								msgPath: data,
							},
						})

						return resp, nil
					},
				},
				routers.CreateRouteParams{
					Pattern: "/messages/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
					QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
						obj := objects.ObjInKey{
							Key: path,
							Obj: new(messageForTest),
						}

						amount := store.Objects().Retrieve(&obj)
						if amount != 1 {
							str := fmt.Sprintf("there is no message on path: %s - %v", path, store.Objects().Keys().Search("[a-z/]+"))
							return nil, errors.New(str)
						}

						js, jsErr := cdc.MarshalJSON(obj.Obj)
						if jsErr != nil {
							return nil, jsErr
						}

						resp := routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
							Code:  routers.IsSuccessful,
							Log:   "success",
							Key:   path,
							Value: js,
						})

						return resp, nil
					},
				},
			},
		},
	})

	// create the applications:
	apps := applications.SDKFunc.CreateApplications(applications.CreateApplicationsParams{
		Apps: []applications.Application{
			app,
		},
	})

	// create the blockchain:
	blkChain := SDKFunc.CreateBlockchain(CreateBlockchainParams{
		Namespace: namespace,
		Name:      name,
		ID:        &id,
		PrivKey:   privKey,
	})

	// create the blockchain service:
	blkChainService := SDKFunc.CreateBlockchainService(CreateBlockchainServiceParams{
		RootDirPath: rootDir,
	})

	// save the blockchain:
	saveBlkChainErr := blkChainService.Save(blkChain)
	if saveBlkChainErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveBlkChainErr.Error())
		return
	}

	// create the application service:
	appService := SDKFunc.CreateApplicationService()

	// spawn the node:
	node, nodeErr := appService.Spawn(port, rootDir, blkChain, apps)
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return
	}
	defer node.Stop()

	address := node.GetAddress()
	client, clientErr := appService.Connect(address)
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}

	// start the node:
	startNodeErr := node.Start()
	if startNodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", startNodeErr.Error())
		return
	}

	// create a new message:
	firstID := uuid.NewV4()
	firstMsg := messageForTest{
		ID:          &firstID,
		Title:       "this is a title",
		Description: "this is a description",
	}

	jsFirstMsg, jsFirstMsgErr := cdc.MarshalJSON(firstMsg)
	if jsFirstMsgErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsFirstMsgErr.Error())
		return
	}

	// create the resource:
	firstRes := routers.SDKFunc.CreateResource(routers.CreateResourceParams{
		ResPtr: routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
			From: fromPubKey,
			Path: "/messages",
		}),
		Data: jsFirstMsg,
	})

	// sign the resource:
	firstSig := fromPrivKey.Sign(firstRes.Hash())

	// save the message:
	trxResp, trxRespErr := client.Transact(routers.SDKFunc.CreateTransactionRequest(routers.CreateTransactionRequestParams{
		Res: firstRes,
		Sig: firstSig,
	}))

	if trxRespErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", trxRespErr.Error())
		return
	}

	retTrxCode := trxResp.Transaction().Code()
	if retTrxCode != routers.IsSuccessful {
		t.Errorf("the transaction was expected to be successful")
		return
	}

	retTrxGazUsed := trxResp.Transaction().GazUsed()
	if retTrxGazUsed != 1205 {
		t.Errorf("the returned gaz used was expected to be %d, returned: %d", 1205, retTrxGazUsed)
		return
	}

	retTrxLog := trxResp.Transaction().Log()
	if retTrxLog != "success" {
		t.Errorf("the returned log was expected to be: %s, returned: %s", "success", retTrxLog)
		return
	}

	retTrxTags := trxResp.Transaction().Tags()
	expectedTags := map[string][]byte{fmt.Sprintf("/messages/%s", firstID.String()): jsFirstMsg}
	if !reflect.DeepEqual(retTrxTags, expectedTags) {
		t.Errorf("the returned tags are invalid")
		return
	}

	// create the resource pointer:
	queryPath := fmt.Sprintf("/messages/%s", firstID.String())
	queryResPtr := routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
		From: fromPubKey,
		Path: queryPath,
	})

	// create the signature:
	querySig := fromPrivKey.Sign(queryResPtr.Hash())

	// execute a query:
	queryResp, queryRespErr := client.Query(routers.SDKFunc.CreateQueryRequest(routers.CreateQueryRequestParams{
		Ptr: queryResPtr,
		Sig: querySig,
	}))

	if queryRespErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", queryRespErr.Error())
		return
	}

	retQueryCode := queryResp.Code()
	if retQueryCode != routers.IsSuccessful {
		t.Errorf("the query ewas expected to be successful")
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
	queryMsg := new(messageForTest)
	jsErr := cdc.UnmarshalJSON(retQueryValue, queryMsg)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return
	}

	if !reflect.DeepEqual(queryMsg, &firstMsg) {
		t.Errorf("the returned message is invalid")
		return
	}
}
