package users

import (
	"testing"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/helpers"
)

func TestSave_thenRetrieve_Success(t *testing.T) {
	//variables:
	pubKey := crypto.SDKFunc.GenPK().PublicKey()

	//create app:
	app := createConcreteUsers()

	//retrieve the head:
	head := app.Objects().Keys().Head()
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
	lenIsZero := app.Objects().Keys().Len()
	if lenIsZero != 0 {
		t.Errorf("the length was expected to be 0: returned: %d", lenIsZero)
		return
	}

	//delete should fail:
	isDeleted := app.Delete(pubKey)
	if isDeleted {
		t.Errorf("the returned bool was expected to be false, true returned")
		return
	}

	//insert the user:
	isInserted := app.Insert(pubKey)
	if !isInserted {
		t.Errorf("the returned bool was expected to be true, false returned")
		return
	}

	// convert with GOB:
	gobData, gobDataErr := helpers.GetBytes(app)
	if gobDataErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", gobDataErr.Error())
		return
	}

	ptr := new(concreteUsers)
	gobErr := helpers.Marshal(gobData, ptr)
	if gobErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", gobErr.Error())
		return
	}

	if !app.Objects().Keys().Head().Head().Compare(ptr.Objects().Keys().Head().Head()) {
		t.Errorf("there was an error while converting the hashtree backandforth using gob")
		return
	}

	//user should exists:
	if !app.Exists(pubKey) {
		t.Errorf("the given user should exists")
		return
	}

	//insert the user again:
	isInsertedAgain := app.Insert(pubKey)
	if isInsertedAgain {
		t.Errorf("the returned bool was expected to be false, true returned")
		return
	}

	//get the head again:
	againHead := app.Objects().Keys().Head()
	if againHead.Length() != 4 {
		t.Errorf("there was supposed to be 4 elements in the head hashtree, returned: %d", againHead.Length())
		return
	}

	//the lenght should be one:
	lenIsOne := app.Objects().Keys().Len()
	if lenIsOne != 1 {
		t.Errorf("the length was expected to be 1: returned: %d", lenIsOne)
		return
	}

	//delete the user:
	isDelSuccess := app.Delete(pubKey)
	if !isDelSuccess {
		t.Errorf("the returned bool was expected to be true, false returned")
		return
	}
}
