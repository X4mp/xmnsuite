package users

import (
	"testing"

	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

func TestSave_thenRetrieve_Success(t *testing.T) {
	//variables:
	pubKey := ed25519.GenPrivKey().PubKey()

	//create app:
	app := createConcreteUsers()

	//retrieve the head:
	head := app.Head()
	if head.Length() != 2 {
		t.Errorf("there was supposed to be 2 elements in the head hashtree, returned: %d", head.Length())
		return
	}

	//user should NOT exists:
	if app.Exists(pubKey) {
		t.Errorf("the given user should not exists")
		return
	}

	//the lenght should be zero:
	lenIsZero := app.Len()
	if lenIsZero != 0 {
		t.Errorf("the length was expected to be 0: returned: %d", lenIsZero)
		return
	}

	//delete should fail:
	delErr := app.Delete(pubKey)
	if delErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	//insert the user:
	insertErr := app.Insert(pubKey)
	if insertErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", insertErr.Error())
		return
	}

	//user should exists:
	if !app.Exists(pubKey) {
		t.Errorf("the given user should exists")
		return
	}

	//insert the user again:
	insertAgainErr := app.Insert(pubKey)
	if insertAgainErr == nil {
		t.Errorf("the returned error was expected to be an error, nil returned")
		return
	}

	//get the head again:
	againHead := app.Head()
	if againHead.Length() != 4 {
		t.Errorf("there was supposed to be 4 elements in the head hashtree, returned: %d", againHead.Length())
		return
	}

	//the lenght should be one:
	lenIsOne := app.Len()
	if lenIsOne != 1 {
		t.Errorf("the length was expected to be 1: returned: %d", lenIsOne)
		return
	}

	//delete the user:
	delSuccessErr := app.Delete(pubKey)
	if delSuccessErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", delSuccessErr.Error())
		return
	}
}
