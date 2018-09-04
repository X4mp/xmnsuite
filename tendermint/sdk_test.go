package tendermint

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	applications "github.com/XMNBlockchain/datamint/applications"
	datastore "github.com/XMNBlockchain/datamint/datastore"
	objects "github.com/XMNBlockchain/datamint/objects"
	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
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

	namespace := "testapp"
	name := "MyTestApp"
	id := uuid.NewV4()
	version := "2018.04.29"
	privKey := ed25519.GenPrivKey()
	store := datastore.SDKFunc.Create()
	fromPrivKey := ed25519.GenPrivKey()
	fromPubKey := fromPrivKey.PubKey()

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
		RouterParams: applications.CreateRouterParams{
			DataStore: routerDS,
			RoleKey:   routerRoleKey,
			RtesParams: []applications.CreateRouteParams{
				applications.CreateRouteParams{
					Pattern: "/messages",
					SaveTrx: func(from crypto.PubKey, path string, params map[string]string, data []byte, sig []byte) (applications.TransactionResponse, error) {

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

						resp := applications.SDKFunc.CreateTransactionResponse(applications.CreateTransactionResponseParams{
							Code:    applications.IsSuccessful,
							Log:     "success",
							GazUsed: 1205,
							Tags: map[string][]byte{
								msgPath: data,
							},
						})

						return resp, nil
					},
				},
				applications.CreateRouteParams{
					Pattern: "/messages/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
					QueryTrx: func(from crypto.PubKey, path string, params map[string]string, sig []byte) (applications.QueryResponse, error) {
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

						resp := applications.SDKFunc.CreateQueryResponse(applications.CreateQueryResponseParams{
							Code:  applications.IsSuccessful,
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
	appService := SDKFunc.CreateApplicationService(CreateApplicationServiceParams{
		RootDir:  rootDir,
		BlkChain: blkChain,
		Apps:     apps,
	})

	// spawn the node:
	node, nodeErr := appService.Spawn()
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

	// start the client:
	startErr := client.Start()
	if startErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", startErr.Error())
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
	firstRes := applications.SDKFunc.CreateResource(applications.CreateResourceParams{
		ResPtr: applications.SDKFunc.CreateResourcePointer(applications.CreateResourcePointerParams{
			From: fromPubKey,
			Path: "/messages",
		}),
		Data: jsFirstMsg,
	})

	// sign the resource:
	firstSig, firstSigErr := fromPrivKey.Sign(firstRes.Hash())
	if firstSigErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSigErr.Error())
		return
	}

	// save the message:
	trxResp, trxRespErr := client.Transact(applications.SDKFunc.CreateTransactionRequest(applications.CreateTransactionRequestParams{
		Res: firstRes,
		Sig: firstSig,
	}))

	if trxRespErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", trxRespErr.Error())
		return
	}

	retTrxCode := trxResp.Transaction().Code()
	if retTrxCode != applications.IsSuccessful {
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
	queryResPtr := applications.SDKFunc.CreateResourcePointer(applications.CreateResourcePointerParams{
		From: fromPubKey,
		Path: queryPath,
	})

	// create the signature:
	querySig, querySigErr := fromPrivKey.Sign(queryResPtr.Hash())
	if querySigErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", querySigErr.Error())
		return
	}

	// execute a query:
	queryResp, queryRespErr := client.Query(applications.SDKFunc.CreateQueryRequest(applications.CreateQueryRequestParams{
		Ptr: queryResPtr,
		Sig: querySig,
	}))

	if queryRespErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", queryRespErr.Error())
		return
	}

	retQueryCode := queryResp.Code()
	if retQueryCode != applications.IsSuccessful {
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
