package roles

import (
	"reflect"
	"testing"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore/lists"
	"github.com/xmnservices/xmnsuite/datastore/users"
	"github.com/xmnservices/xmnsuite/helpers"
)

func TestSingle_save_thenExists_thenRetrieve_thenDelete_Success(t *testing.T) {
	//variables:

	firstPK := crypto.SDKFunc.GenPK()
	secondPK := crypto.SDKFunc.GenPK()
	thirdPK := crypto.SDKFunc.GenPK()
	key := "first-role"
	writeOnKey := "i-want-to-write-on-this-key"
	writeOnKeyPattern := "[a-z-]+"

	//add users:
	users := users.SDKFunc.Create()
	users.Insert(firstPK.PublicKey())
	users.Insert(secondPK.PublicKey())
	users.Insert(thirdPK.PublicKey())

	//create a list:
	setApp := lists.SDKFunc.CreateSet()
	setApp.Add(writeOnKey, []byte("some data"))

	//create roles:
	app := createConcreteRoles()

	//add users:
	retAmountAdded := app.Add(key, firstPK.PublicKey(), secondPK.PublicKey(), thirdPK.PublicKey())
	if retAmountAdded != 3 {
		t.Errorf("the returned amount was expected to be 3, returned: %d", retAmountAdded)
		return
	}

	retAmountDeleted := app.Del(key, secondPK.PublicKey())
	if retAmountDeleted != 1 {
		t.Errorf("the returned amount was expected to be 1, returned: %d", retAmountDeleted)
		return
	}

	//should not have write access:
	retKeys := app.HasWriteAccess(key, users.Key(firstPK.PublicKey()), writeOnKey)
	if len(retKeys) != 0 {
		t.Errorf("there should be 0 keys that we have write access to, returned; %d", len(retKeys))
		return
	}

	//add the write access:
	retAmountEnabledKeys := app.EnableWriteAccess(key, users.Key(firstPK.PublicKey()), writeOnKeyPattern)
	if retAmountEnabledKeys != 2 {
		t.Errorf("there should now be 2 keys where we have write access, returned: %d", retAmountEnabledKeys)
		return
	}

	//should now have write access:
	retValidWriteAccessKeys := app.HasWriteAccess(key, users.Key(firstPK.PublicKey()), writeOnKey)
	expected := []string{users.Key(firstPK.PublicKey()), writeOnKey}
	if !reflect.DeepEqual(retValidWriteAccessKeys, expected) {
		t.Errorf("the returned keys are invalid")
		return
	}

	//disable the write access on the user:
	retAmountDisabled := app.DisableWriteAccess(key, users.Key(firstPK.PublicKey()))
	if retAmountDisabled != 1 {
		t.Errorf("there should now be 1 new disabled key, returned: %d", retAmountDisabled)
		return
	}

	//should now hav e write access to 1 key:
	retSingleWriteAccessKeys := app.HasWriteAccess(key, users.Key(firstPK.PublicKey()), writeOnKey)
	singleExpected := []string{writeOnKey}
	if !reflect.DeepEqual(retSingleWriteAccessKeys, singleExpected) {
		t.Errorf("the returned keys are invalid")
		return
	}

	// convert with GOB:
	gobData, gobDataErr := helpers.GetBytes(app)
	if gobDataErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", gobDataErr.Error())
		return
	}

	ptr := new(concreteRoles)
	gobErr := helpers.Marshal(gobData, ptr)
	if gobErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", gobErr.Error())
		return
	}

	if !app.Lists().Objects().Keys().Head().Head().Compare(ptr.Lists().Objects().Keys().Head().Head()) {
		t.Errorf("there was an error while converting the hashtree backandforth using gob")
		return
	}

}
