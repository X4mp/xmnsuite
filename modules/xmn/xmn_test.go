package xmn

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/XMNBlockchain/datamint/applications"
	datastore "github.com/XMNBlockchain/datamint/datastore"
	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

type messageForTest struct {
	ID          *uuid.UUID `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

func TestXMN_Success(t *testing.T) {

	// variables:
	instanceID := uuid.NewV4()
	dbPath := "./test_files"
	defer func() {
		os.RemoveAll(dbPath)
	}()

	nodePK := ed25519.GenPrivKey()
	fromPrivKey := ed25519.GenPrivKey()
	fromPubKey := fromPrivKey.PubKey()

	// create the initial write keys:
	firstPrivKey := ed25519.GenPrivKey()
	secondPrivKey := ed25519.GenPrivKey()
	thirdPrivKey := ed25519.GenPrivKey()
	rootPubKeys := []crypto.PubKey{
		firstPrivKey.PubKey(),
		secondPrivKey.PubKey(),
		thirdPrivKey.PubKey(),
		fromPubKey,
	}

	//create lua state:
	l := createLuaState()
	defer l.Close()

	// create the datastore:
	ds := datastore.SDKFunc.Create()

	// create XMN:
	xmn := createXMN(ds)

	// register:
	xmn.register(l)

	// execute:
	node, nodeErr := xmn.execute(l, dbPath, &instanceID, rootPubKeys, nodePK, "lua/core_test.lua")
	if nodeErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", nodeErr.Error())
		return
	}
	defer node.Stop()

	// retrieve the client:
	client, clientErr := node.GetClient()
	if clientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", clientErr.Error())
		return
	}

	startClientErr := client.Start()
	if startClientErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", startClientErr.Error())
		return
	}

	// create a new message:
	firstID := uuid.NewV4()
	firstMsg := messageForTest{
		ID:          &firstID,
		Title:       "this is a title",
		Description: "this is a description",
	}

	// save the message:
	saveMessage(t, client, fromPrivKey, &firstMsg)

	// retrieve the message:
	retrieveMessage(t, client, fromPrivKey, &firstMsg)

	// delete the message:
	deleteMessage(t, client, fromPrivKey, &firstMsg)

	// retrieve the message, not found:
	retrieveMessageNotFound(t, client, fromPrivKey, &firstMsg)

}

func TestChain_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	// create XMN:
	xmn := createXMN(datastore.SDKFunc.Create())

	// register:
	xmn.register(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/chain_test.lua")
}

func TestPrivKey_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	// create XMN:
	xmn := createXMN(datastore.SDKFunc.Create())

	// register:
	xmn.register(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/privkey_test.lua")
}

func TestKeys_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	// create XMN:
	xmn := createXMN(datastore.SDKFunc.Create())

	// register:
	xmn.register(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/keys_test.lua")
}

func TestRoles_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	// create XMN:
	xmn := createXMN(datastore.SDKFunc.Create())

	// register:
	xmn.register(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/roles_test.lua")
}

func TestTables_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	// create XMN:
	xmn := createXMN(datastore.SDKFunc.Create())

	// register:
	xmn.register(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/tables_test.lua")
}

func TestUsers_Success(t *testing.T) {

	//create lua state:
	l := createLuaState()
	defer l.Close()

	// create XMN:
	xmn := createXMN(datastore.SDKFunc.Create())

	// register:
	xmn.register(l)

	//execute the chunk:
	executeChunkForTests(l, "lua/users_test.lua")
}

func saveMessage(t *testing.T, client applications.Client, fromPrivKey crypto.PrivKey, firstMsg *messageForTest) {
	jsFirstMsg, jsFirstMsgErr := json.Marshal(firstMsg)
	if jsFirstMsgErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsFirstMsgErr.Error())
		return
	}

	// create the resource:
	firstRes := applications.SDKFunc.CreateResource(applications.CreateResourceParams{
		ResPtr: applications.SDKFunc.CreateResourcePointer(applications.CreateResourcePointerParams{
			From: fromPrivKey.PubKey(),
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

	tagKeyname := fmt.Sprintf("/messages/%s", firstMsg.ID.String())
	retTrxTags := trxResp.Transaction().Tags()
	if tagValueAsBytes, ok := retTrxTags[tagKeyname]; ok {
		decodedFirstMsg := new(messageForTest)
		decodedFirstMsgErr := json.Unmarshal(tagValueAsBytes, decodedFirstMsg)
		if decodedFirstMsgErr != nil {
			t.Errorf("the returned value was expected to be nil, error returned: %s", decodedFirstMsgErr.Error())
			return
		}

		if !reflect.DeepEqual(decodedFirstMsg, firstMsg) {
			t.Errorf("the value in the tags (keyname: %s) is invalid", tagKeyname)
			return
		}

	}

	if _, ok := retTrxTags[tagKeyname]; !ok {
		t.Errorf("the keyname (%s) was expected in the tags", tagKeyname)
		return
	}
}

func retrieveMessage(t *testing.T, client applications.Client, fromPrivKey crypto.PrivKey, firstMsg *messageForTest) {
	// create the resource pointer:
	queryPath := fmt.Sprintf("/messages/%s", firstMsg.ID.String())
	queryResPtr := applications.SDKFunc.CreateResourcePointer(applications.CreateResourcePointerParams{
		From: fromPrivKey.PubKey(),
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
		t.Errorf("the query was expected to be successful.  Expected: %d, Returned: %d", applications.IsSuccessful, retQueryCode)
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
	jsErr := json.Unmarshal(retQueryValue, queryMsg)
	if jsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", jsErr.Error())
		return
	}

	if !reflect.DeepEqual(queryMsg, firstMsg) {
		t.Errorf("the returned message is invalid")
		return
	}
}

func retrieveMessageNotFound(t *testing.T, client applications.Client, fromPrivKey crypto.PrivKey, firstMsg *messageForTest) {
	// create the resource pointer:
	queryPath := fmt.Sprintf("/messages/%s", firstMsg.ID.String())
	queryResPtr := applications.SDKFunc.CreateResourcePointer(applications.CreateResourcePointerParams{
		From: fromPrivKey.PubKey(),
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
	if retQueryCode != applications.NotFound {
		t.Errorf("the query was expected to be not found.  Expected: %d, Returned: %d", applications.NotFound, retQueryCode)
		return
	}

	retQueryKey := queryResp.Key()
	if retQueryKey != queryPath {
		t.Errorf("the returned key was expected to be: %s, returned: %s", queryPath, retQueryKey)
		return
	}

	retQueryLog := queryResp.Log()
	if retQueryLog != "not found" {
		t.Errorf("the returned log was expected to be: %s, returned: %s", "success", retQueryLog)
		return
	}

	retQueryValue := queryResp.Value()
	if retQueryValue != nil {
		t.Errorf("the value was expected to be nil")
		return
	}
}

func deleteMessage(t *testing.T, client applications.Client, fromPrivKey crypto.PrivKey, firstMsg *messageForTest) {
	// create the resource pointer:
	path := fmt.Sprintf("/messages/%s", firstMsg.ID.String())
	resPtr := applications.SDKFunc.CreateResourcePointer(applications.CreateResourcePointerParams{
		From: fromPrivKey.PubKey(),
		Path: path,
	})

	// sign the resource pointer:
	firstSig, firstSigErr := fromPrivKey.Sign(resPtr.Hash())
	if firstSigErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstSigErr.Error())
		return
	}

	// delete the message:
	trxResp, trxRespErr := client.Transact(applications.SDKFunc.CreateTransactionRequest(applications.CreateTransactionRequestParams{
		Ptr: resPtr,
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
	if retTrxGazUsed != 0 {
		t.Errorf("the returned gaz used was expected to be %d, returned: %d", 0, retTrxGazUsed)
		return
	}

	retTrxLog := trxResp.Transaction().Log()
	if retTrxLog != "success" {
		t.Errorf("the returned log was expected to be: %s, returned: %s", "success", retTrxLog)
		return
	}

	retTrxTags := trxResp.Transaction().Tags()
	expectedTags := map[string][]byte{}
	if !reflect.DeepEqual(retTrxTags, expectedTags) {
		t.Errorf("the returned tags are invalid")
		return
	}
}
