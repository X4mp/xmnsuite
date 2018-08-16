package objects

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSingle_save_thenExists_thenRetrieve_thenDelete_Success(t *testing.T) {

	//variables:
	obj := createObjForTests()
	key := fmt.Sprintf("entity:by_id:%s", obj.ID.String())

	//create the application:
	app := createObjects()

	//the lenght should be zero:
	lenIsZero := app.Keys().Len()
	if lenIsZero != 0 {
		t.Errorf("the length was expected to be 0: returned: %d", lenIsZero)
		return
	}

	//retrieve the head:
	head := app.Keys().Head()
	if head.Length() != 2 {
		t.Errorf("there was supposed to be 2 elements in the head hashtree, returned: %d", head.Length())
		return
	}

	//retrieve the ht, should ne nil:
	htIsNil := app.Keys().HashTree(key)
	if htIsNil != nil {
		t.Errorf("the returned hashtree was expected to be nil")
		return
	}

	//the object does not exists:
	amountExistsIsZero := app.Keys().Exists(key)
	if amountExistsIsZero != 0 {
		t.Errorf("the amount exists was expected to be 0, %d returned", amountExistsIsZero)
		return
	}

	//save the object:
	amountSaved := app.Save(&ObjInKey{
		Key: key,
		Obj: obj,
	})

	if amountSaved != 1 {
		t.Errorf("the amount saved was expected to be 1, %d returned", amountSaved)
		return
	}

	//the object does not exists:
	amountExistsIsOne := app.Keys().Exists(key)
	if amountExistsIsOne != 1 {
		t.Errorf("the amount exists was expected to be 1, %d returned", amountExistsIsOne)
		return
	}

	//retrieve the object:
	retObjInKey := ObjInKey{
		Key: key,
		Obj: new(objFortests),
	}

	amountRet := app.Retrieve(&retObjInKey)
	if amountRet != 1 {
		t.Errorf("the amount retrieved was expected to be 1, %d returned", amountRet)
		return
	}

	if !reflect.DeepEqual(retObjInKey.Obj, obj) {
		t.Errorf("the retrieved object is invalid\n Expected: %v\n Returned: %v\n", obj, retObjInKey.Obj)
		return
	}

	//the lenght should be one:
	lenIsOne := app.Keys().Len()
	if lenIsOne != 1 {
		t.Errorf("the length was expected to be 1: returned: %d", lenIsOne)
		return
	}

	//retrieve the head again:
	headAgain := app.Keys().Head()
	if headAgain.Length() != 4 {
		t.Errorf("there was supposed to be 4 elements in the head hashtree, returned: %d", headAgain.Length())
		return
	}

	//retrieve the ht:
	htAgain := app.Keys().HashTree(key)

	//retrieve the ht from list:
	htList := app.Keys().HashTrees(key)
	if len(htList) != 1 {
		t.Errorf("there was supposed to be 1 hashtree in the list")
		return
	}

	if !reflect.DeepEqual(htAgain, htList[0]) {
		t.Errorf("the hashtrees are invalid")
		return
	}

	//delete the object:
	amountDel := app.Keys().Delete(key)
	if amountDel != 1 {
		t.Errorf("the amount deleted was expected to be 1, %d returned", amountDel)
		return
	}

	//the lenght should be zero again:
	lenIsZeroAgain := app.Keys().Len()
	if lenIsZeroAgain != 0 {
		t.Errorf("the length was expected to be 0: returned: %d", lenIsZeroAgain)
		return
	}

}

func TestSingle_save_withNonJSONObj_panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code was expected to panic")
		}
	}()

	//variables:
	obj := createNonJSONObjForTests()
	key := fmt.Sprintf("entity:by_id:%s", obj.ID.String())

	//create the application:
	app := createObjects()

	//save the object:
	app.Save(&ObjInKey{
		Key: key,
		Obj: obj,
	})
}

func TestSingle_saveInvalid_retrieve_panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code was expected to panic")
		}
	}()

	//variables:
	data := []byte("this is some data")
	obj := createObjForTests()
	key := fmt.Sprintf("entity:by_id:%s", obj.ID.String())

	//create the application:
	app := createObjects()

	//save data in the keys:
	app.Keys().Save(key, data)

	//try to retrieve it as an object:
	app.Retrieve(&ObjInKey{
		Key: key,
		Obj: new(objFortests),
	})
}
