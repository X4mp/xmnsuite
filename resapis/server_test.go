package restapis

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/xmnservices/xmnsuite/crypto"
)

func TestCreateAccount_Success(t *testing.T) {
	// variables:
	pk := crypto.SDKFunc.GenPK()
	dbPath := "./test_files"
	rter := mux.NewRouter()
	rep := createRepository(dbPath)
	serv := createService(dbPath)
	port := 8080
	gracefulTimeout := time.Second * 15
	clientURL, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	name := "steve-rodrigue"
	seedwords := []string{"je", "mange", "du", "spaghetti", "et", "je", "dois", "avouer", "que", "j'aime", "bien", "ca"}
	defer func() {
		os.RemoveAll(dbPath)
	}()

	// create server:
	server := createServer(pk, rter, rep, serv, port, gracefulTimeout)

	// start the server:
	startErr := server.Start()
	if startErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", startErr.Error())
		return
	}
	defer server.Stop()

	// create client:
	cl := createClient(pk, clientURL)

	// create the account:
	accountErr := cl.CreateAccount(name, seedwords)
	if accountErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", accountErr.Error())
		return
	}

	// retrieve the account by name:
	retAccount, retAccountErr := cl.RetrieveAccountByName(name)
	if retAccountErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retAccountErr.Error())
		return
	}

	if retAccount.Name() != name {
		t.Errorf("the account name was expected to be: %s, returned: %s", name, retAccount.Name())
		return
	}

	decryptedPK := retAccount.DecryptPK(seedwords)
	if decryptedPK == nil {
		t.Errorf("the decrypted PK was expected to be valid")
		return
	}

	// retrieve the accounts:
	retAccounts, retAccountsErr := cl.RetrieveAccounts()
	if retAccountsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retAccountsErr.Error())
		return
	}

	if len(retAccounts) != 1 {
		t.Errorf("1 account was expected, returned: %d", len(retAccounts))
		return
	}

	if !reflect.DeepEqual(retAccount, retAccounts[0]) {
		t.Errorf("the returned accounts are invalid")
		return
	}
}
