package lists

import (
	"reflect"
	"testing"

	"github.com/XMNBlockchain/datamint"
)

func TestAdd_thenRetrieve_notUnique_Success(t *testing.T) {
	//variables:
	firstElement := []byte("this is the element")
	secondElement := []byte("this is the element")
	thirdElement := []byte("this is the third element")
	key := "this-is-a-key"
	emptyKey := "this-is-an-empty-key"

	//create the app:
	app := createConcreteLists(false)

	//create an empty key:
	emptyData, _ := datamint.GetBytes([]interface{}{})
	app.Objects().Keys().Save(emptyKey, emptyData)

	//verify the length:
	firstLength := app.Len(emptyKey)
	if firstLength != 0 {
		t.Errorf("tre returned length was expected to be 0, returned: %d", firstLength)
		return
	}

	//add the list:
	firstAddAmount := app.Add(key, firstElement)
	if firstAddAmount != 1 {
		t.Errorf("the returned amount was expected to be 1, returned: %d", firstAddAmount)
		return
	}

	//verify the length:
	secondLength := app.Len(key)
	if secondLength != 1 {
		t.Errorf("tre returned length was expected to be 1, returned: %d", secondLength)
		return
	}

	//add again:
	secondAmount := app.Add(key, secondElement, thirdElement)
	if secondAmount != 2 {
		t.Errorf("the returned amount was expected to be 2, returned: %d", secondAmount)
		return
	}

	//verify the length:
	thirdLength := app.Len(key)
	if thirdLength != 3 {
		t.Errorf("tre returned length was expected to be 3, returned: %d", thirdLength)
		return
	}

	//retrieve the last element:
	firstRetKeys := app.Retrieve(key, 2, 1)
	if !reflect.DeepEqual(firstRetKeys, []interface{}{thirdElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", firstRetKeys, []interface{}{thirdElement})
		return
	}

	//retrieve the first element:
	secondRetKeys := app.Retrieve(key, 0, 1)
	if !reflect.DeepEqual(secondRetKeys, []interface{}{firstElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", secondRetKeys, []interface{}{firstElement})
		return
	}

	//retrieve the first two elements:
	thirdRetKeys := app.Retrieve(key, 0, 2)
	if !reflect.DeepEqual(thirdRetKeys, []interface{}{firstElement, secondElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", thirdRetKeys, []interface{}{firstElement, secondElement})
		return
	}

	//retrieve the last two elements:
	fourthRetKeys := app.Retrieve(key, 1, 2)
	if !reflect.DeepEqual(fourthRetKeys, []interface{}{secondElement, thirdElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", fourthRetKeys, []interface{}{secondElement, thirdElement})
		return
	}

	//retrieve all elements:
	fifthRetKeys := app.Retrieve(key, 0, -1)
	if !reflect.DeepEqual(fifthRetKeys, []interface{}{firstElement, secondElement, thirdElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", fifthRetKeys, []interface{}{firstElement, secondElement, thirdElement})
		return
	}

	//retrieve all elements with an amount way too large:
	sixthRetKeys := app.Retrieve(key, 0, 25)
	if !reflect.DeepEqual(sixthRetKeys, []interface{}{firstElement, secondElement, thirdElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", sixthRetKeys, []interface{}{firstElement, secondElement, thirdElement})
		return
	}

	//retrieve elements with an index way too high:
	seventhRetKeys := app.Retrieve(key, 3, 2)
	if !reflect.DeepEqual(seventhRetKeys, []interface{}{}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", seventhRetKeys, []interface{}{})
		return
	}

	//retrieve from an empty key:
	retFromEmptyKey := app.Retrieve(emptyKey, 0, 2)
	if !reflect.DeepEqual(retFromEmptyKey, []interface{}{}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", retFromEmptyKey, []interface{}{})
		return
	}

	//retrieve with an invalid index:
	firstIsNil := app.Retrieve(key, -4, 1)
	if firstIsNil != nil {
		t.Errorf("the returned value was expected to be nil, value returned: %v", firstIsNil)
		return
	}

	//retrieve with an invalid amount:
	secondIsNil := app.Retrieve(key, 0, -4)
	if secondIsNil != nil {
		t.Errorf("the returned value was expected to be nil, value returned: %v", secondIsNil)
		return
	}

	//retrieve from an invalid key:
	thirdIsNil := app.Retrieve("this-is-invalid-yes", 0, 2)
	if thirdIsNil != nil {
		t.Errorf("the returned value was expected to be nil, value returned: %v", thirdIsNil)
		return
	}
}

func TestAdd_thenRetrieve_isUnique_Success(t *testing.T) {
	//variables:
	firstElement := []byte("this is the element")
	secondElement := []byte("this is the element")
	thirdElement := []byte("this is the third element")
	key := "this-is-a-key"
	emptyKey := "this-is-an-empty-key"

	//create the app:
	app := createConcreteLists(true)

	//verify the length:
	invalidLength := app.Len(key)
	if invalidLength != 0 {
		t.Errorf("tre returned length was expected to be 0, returned: %d", invalidLength)
		return
	}

	//create an empty key:
	emptyData, _ := datamint.GetBytes([]interface{}{})
	app.Objects().Keys().Save(emptyKey, emptyData)

	//add the list:
	firstAddAmount := app.Add(key, firstElement)
	if firstAddAmount != 1 {
		t.Errorf("the returned amount was expected to be 1, returned: %d", firstAddAmount)
		return
	}

	//verify the length:
	firstLength := app.Len(key)
	if firstLength != 1 {
		t.Errorf("tre returned length was expected to be 1, returned: %d", firstLength)
		return
	}

	//add again:
	secondAmount := app.Add(key, secondElement, thirdElement)
	if secondAmount != 1 {
		t.Errorf("the returned amount was expected to be 1, returned: %d", secondAmount)
		return
	}

	//verify the length:
	secondLength := app.Len(key)
	if secondLength != 2 {
		t.Errorf("tre returned length was expected to be 2, returned: %d", secondLength)
		return
	}

	//retrieve the last element:
	firstRetKeys := app.Retrieve(key, 1, 1)
	if !reflect.DeepEqual(firstRetKeys, []interface{}{thirdElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", firstRetKeys, []interface{}{thirdElement})
		return
	}

	//retrieve the first element:
	secondRetKeys := app.Retrieve(key, 0, 1)
	if !reflect.DeepEqual(secondRetKeys, []interface{}{firstElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", secondRetKeys, []interface{}{firstElement})
		return
	}

	//retrieve all elements:
	fifthRetKeys := app.Retrieve(key, 0, -1)
	if !reflect.DeepEqual(fifthRetKeys, []interface{}{firstElement, thirdElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", fifthRetKeys, []interface{}{firstElement, thirdElement})
		return
	}

	//retrieve all elements with an amount way too large:
	sixthRetKeys := app.Retrieve(key, 0, 25)
	if !reflect.DeepEqual(sixthRetKeys, []interface{}{firstElement, thirdElement}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", sixthRetKeys, []interface{}{firstElement, thirdElement})
		return
	}

	//retrieve elements with an index way too high:
	seventhRetKeys := app.Retrieve(key, 2, 2)
	if !reflect.DeepEqual(seventhRetKeys, []interface{}{}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", seventhRetKeys, []interface{}{})
		return
	}

	//retrieve from en empty key:
	retFromEmptyKey := app.Retrieve(emptyKey, 0, 2)
	if !reflect.DeepEqual(retFromEmptyKey, []interface{}{}) {
		t.Errorf("the returned keys are invalid: \nExpected: %v\nReturned: %v\n\n", retFromEmptyKey, []interface{}{})
		return
	}

	//retrieve with an invalid index:
	firstIsNil := app.Retrieve(key, -4, 1)
	if firstIsNil != nil {
		t.Errorf("the returned value was expected to be nil, value returned: %v", firstIsNil)
		return
	}

	//retrieve with an invalid amount:
	secondIsNil := app.Retrieve(key, 0, -4)
	if secondIsNil != nil {
		t.Errorf("the returned value was expected to be nil, value returned: %v", secondIsNil)
		return
	}

	//retrieve from an invalid key:
	thirdIsNil := app.Retrieve("this-is-invalid-yes", 0, 2)
	if thirdIsNil != nil {
		t.Errorf("the returned value was expected to be nil, value returned: %v", thirdIsNil)
		return
	}
}

func TestAdd_thenUnion_notUnique_Success(t *testing.T) {
	//variables:
	firstElement := []byte("this is the element")
	secondElement := []byte("this is the third element")
	firstKey := "this-is-a-first-key"
	secondKey := "this-is-a-second-key"
	thirdKey := "this-is-a-third-key"

	//create the app:
	app := createConcreteLists(false)

	//add the elements in the keys:
	app.Add(firstKey, firstElement, secondElement)
	app.Add(secondKey, firstElement)
	app.Add(thirdKey, firstElement, secondElement)

	//union:
	retUnion := app.Union(firstKey, secondKey, thirdKey)
	expected := []interface{}{firstElement, secondElement, firstElement, firstElement, secondElement}
	if !reflect.DeepEqual(retUnion, expected) {
		t.Errorf("the retrieved union is invalid.  \n Expected: %v \n Returned: %v\n", expected, retUnion)
		return
	}

}

func TestAdd_thenUnion_unique_Success(t *testing.T) {
	//variables:
	firstElement := []byte("this is the element")
	secondElement := []byte("this is the third element")
	thirdElement := []byte("this is the last element")
	firstKey := "this-is-a-first-key"
	secondKey := "this-is-a-second-key"
	thirdKey := "this-is-a-third-key"

	//create the app:
	app := createConcreteLists(true)

	//add the elements in the keys:
	app.Add(firstKey, firstElement, secondElement)
	app.Add(secondKey, firstElement, thirdElement)
	app.Add(thirdKey, firstElement, secondElement)

	//union:
	retUnion := app.Union(firstKey, secondKey, thirdKey, "this-is-an-invalid-key")
	expected := []interface{}{firstElement, secondElement, thirdElement}
	if !reflect.DeepEqual(retUnion, expected) {
		t.Errorf("the retrieved union is invalid.  \n Expected: %v \n Returned: %v\n", expected, retUnion)
		return
	}

}

func TestAdd_thenUnionStore_unique_Success(t *testing.T) {
	//variables:
	firstElement := []byte("this is the element")
	secondElement := []byte("this is the third element")
	thirdElement := []byte("this is the last element")
	firstKey := "this-is-a-first-key"
	secondKey := "this-is-a-second-key"
	thirdKey := "this-is-a-third-key"
	destination := "this-is-the-destination-key"

	//create the app:
	app := createConcreteLists(true)

	//add the elements in the keys:
	app.Add(firstKey, firstElement, secondElement)
	app.Add(secondKey, firstElement, thirdElement)
	app.Add(thirdKey, firstElement, secondElement)

	//unionstore:
	retAmount := app.UnionStore(destination, firstKey, secondKey, thirdKey, "this-is-an-invalid-key")
	if retAmount != 3 {
		t.Errorf("the returned amount was expected to be 3, returned: %d", retAmount)
		return
	}

	//union:
	retUnion := app.Retrieve(destination, 0, -1)
	expected := []interface{}{firstElement, secondElement, thirdElement}
	if !reflect.DeepEqual(retUnion, expected) {
		t.Errorf("the retrieved union is invalid.  \n Expected: %v \n Returned: %v\n", expected, retUnion)
		return
	}

}
