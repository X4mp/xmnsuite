package keys

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSingle_save_thenExists_thenRetrieve_thenDelete_Success(t *testing.T) {

	//variables:
	data := []byte("this is some data")
	key := fmt.Sprintf("entity:by_slug:%s", "this-is-some-data")

	//create the application:
	app := createConcreteKeys()

	//the lenght should be zero:
	lenIsZero := app.Len()
	if lenIsZero != 0 {
		t.Errorf("the length was expected to be 0: returned: %d", lenIsZero)
		return
	}

	//retrieve the object, but its not stored yet:
	retValueNotThere := app.Retrieve(key)
	if retValueNotThere != nil {
		t.Errorf("the returned value was expected to be nil, value returned: %v", retValueNotThere)
		return
	}

	//retrieve the head:
	head := app.Head()
	if head.Length() != 2 {
		t.Errorf("there was supposed to be 1 element in the head hashtree, returned: %d", head.Length())
		return
	}

	//retrieve the ht, should ne nil:
	htIsNil := app.HashTree(key)
	if htIsNil != nil {
		t.Errorf("the returned hashtree was expected to be nil")
		return
	}

	//the object does not exists:
	amountExistsIsZero := app.Exists(key)
	if amountExistsIsZero != 0 {
		t.Errorf("the amount exists was expected to be 0, %d returned", amountExistsIsZero)
		return
	}

	//save the object:
	app.Save(key, data)

	//the object does not exists:
	amountExistsIsOne := app.Exists(key)
	if amountExistsIsOne != 1 {
		t.Errorf("the amount exists was expected to be 1, %d returned", amountExistsIsOne)
		return
	}

	//retrieve the object:
	retValue := app.Retrieve(key)
	if !reflect.DeepEqual(data, retValue) {
		t.Errorf("the returned data is invalid")
		return
	}

	//the lenght should be one:
	lenIsOne := app.Len()
	if lenIsOne != 1 {
		t.Errorf("the length was expected to be 1: returned: %d", lenIsOne)
		return
	}

	//retrieve the head again:
	headAgain := app.Head()
	if headAgain.Length() != 4 {
		t.Errorf("there was supposed to be 3 elements in the head hashtree, returned: %d", headAgain.Length())
		return
	}

	//retrieve the ht:
	htAgain := app.HashTree(key)

	//retrieve the ht from list:
	htList := app.HashTrees(key)
	if len(htList) != 1 {
		t.Errorf("there was supposed to be 1 hashtree in the list")
		return
	}

	if !reflect.DeepEqual(htAgain, htList[0]) {
		t.Errorf("the hashtrees are invalid")
		return
	}

	//delete the object:
	amountDel := app.Delete(key)
	if amountDel != 1 {
		t.Errorf("the amount deleted was expected to be 1, %d returned", amountDel)
		return
	}

	//the lenght should be zero again:
	lenIsZeroAgain := app.Len()
	if lenIsZeroAgain != 0 {
		t.Errorf("the length was expected to be 0: returned: %d", lenIsZeroAgain)
		return
	}

}

func TestSearch_Success(t *testing.T) {
	//variables:
	data := []byte("this is some data")
	first := fmt.Sprintf("entity:by_slug:%s", "this-is-some-data")
	second := "should-not-be-found"
	third := fmt.Sprintf("another:by_name:%s", "this-is-another-key")
	shouldFind := []string{
		third,
		first,
	}

	//create the application:
	app := createConcreteKeys()

	//save the instances on keys:
	app.Save(first, data)
	app.Save(second, data)
	app.Save(third, data)

	//search:
	results := app.Search("[a-z]+:[a-z_]+:[a-z-]+")
	if !reflect.DeepEqual(shouldFind, results) {
		t.Errorf("the returned keys are invalid.  \n\n Expected: %v, \n Returned: %v\n\n", shouldFind, results)
		return
	}

}

func TestSearch_patternIsInvalid_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code was expected to panic")
		}
	}()

	//create the application:
	app := createConcreteKeys()

	//search:
	app.Search("\\K")
}
