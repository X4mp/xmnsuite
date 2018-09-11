package roles

import (
	"reflect"
	"testing"

	"github.com/xmnservices/xmnsuite/lists"
	"github.com/xmnservices/xmnsuite/users"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

func TestSingle_save_thenExists_thenRetrieve_thenDelete_Success(t *testing.T) {
	//variables:

	firstPK := ed25519.GenPrivKey()
	secondPK := ed25519.GenPrivKey()
	thirdPK := ed25519.GenPrivKey()
	key := "first-role"
	writeOnKey := "i-want-to-write-on-this-key"
	writeOnKeyPattern := "[a-z-]+"

	//add users:
	users := users.SDKFunc.Create()
	users.Insert(firstPK.PubKey())
	users.Insert(secondPK.PubKey())
	users.Insert(thirdPK.PubKey())

	//create a list:
	setApp := lists.SDKFunc.CreateSet()
	setApp.Add(writeOnKey, []byte("some data"))

	//create roles:
	app := createConcreteRoles()

	//add users:
	retAmountAdded := app.Add(key, firstPK.PubKey(), secondPK.PubKey(), thirdPK.PubKey())
	if retAmountAdded != 3 {
		t.Errorf("the returned amount was expected to be 3, returned: %d", retAmountAdded)
		return
	}

	retAmountDeleted := app.Del(key, secondPK.PubKey())
	if retAmountDeleted != 1 {
		t.Errorf("the returned amount was expected to be 1, returned: %d", retAmountDeleted)
		return
	}

	//should not have write access:
	retKeys := app.HasWriteAccess(key, users.Key(firstPK.PubKey()), writeOnKey)
	if len(retKeys) != 0 {
		t.Errorf("there should be 0 keys that we have write access to, returned; %d", len(retKeys))
		return
	}

	//add the write access:
	retAmountEnabledKeys := app.EnableWriteAccess(key, users.Key(firstPK.PubKey()), writeOnKeyPattern)
	if retAmountEnabledKeys != 2 {
		t.Errorf("there should now be 2 keys where we have write access, returned: %d", retAmountEnabledKeys)
		return
	}

	//should now have write access:
	retValidWriteAccessKeys := app.HasWriteAccess(key, users.Key(firstPK.PubKey()), writeOnKey)
	expected := []string{users.Key(firstPK.PubKey()), writeOnKey}
	if !reflect.DeepEqual(retValidWriteAccessKeys, expected) {
		t.Errorf("the returned keys are invalid")
		return
	}

	//disable the write access on the user:
	retAmountDisabled := app.DisableWriteAccess(key, users.Key(firstPK.PubKey()))
	if retAmountDisabled != 1 {
		t.Errorf("there should now be 1 new disabled key, returned: %d", retAmountDisabled)
		return
	}

	//should now hav e write access to 1 key:
	retSingleWriteAccessKeys := app.HasWriteAccess(key, users.Key(firstPK.PubKey()), writeOnKey)
	singleExpected := []string{writeOnKey}
	if !reflect.DeepEqual(retSingleWriteAccessKeys, singleExpected) {
		t.Errorf("the returned keys are invalid")
		return
	}

}
