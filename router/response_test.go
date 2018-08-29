package router

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/XMNBlockchain/datamint/tests"
)

type someDataForTests struct {
	Title string `json:"title"`
	Desc  string `json:"description"`
}

func TestCreateQueryResponse_Success(t *testing.T) {
	//variables
	someData := someDataForTests{
		Title: "this is some data",
		Desc:  "This is the description of the some data",
	}

	data, _ := cdc.MarshalJSON(someData)
	gazUsed := int64(rand.Int() % 20)
	log := "success"

	//execute:
	response, responseErr := createQueryResponse(true, false, false, gazUsed, log, data)
	if responseErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", responseErr.Error())
		return
	}

	retIsSuccess := response.IsSuccess()
	retIsAuthorized := response.IsAuthorized()
	retHasNFS := response.HasInsufficientFunds()
	retGazUsed := response.GazUsed()
	retData := response.Data()
	retLog := response.Log()

	if !retIsSuccess {
		t.Errorf("the returned isSuccess is exepcted to be true, false returned")
		return
	}

	if retIsAuthorized {
		t.Errorf("the returned isAuthorized is exepcted to be false, true returned")
		return
	}

	if retHasNFS {
		t.Errorf("the returned hasNFS is exepcted to be false, true returned")
		return
	}

	if !reflect.DeepEqual(gazUsed, retGazUsed) {
		t.Errorf("the returned gazUsed is invalid.  Expected: %d, Returned: %d", gazUsed, retGazUsed)
		return
	}

	if !reflect.DeepEqual(data, retData) {
		t.Errorf("the returned data is invalid")
		return
	}

	if !reflect.DeepEqual(log, retLog) {
		t.Errorf("the returned log is invalid")
		return
	}

	firstEmpty := new(someDataForTests)
	firstEmptyErr := response.UnMarshal(firstEmpty)
	if firstEmptyErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", firstEmptyErr.Error())
		return
	}

	firstEmptyJS, _ := cdc.MarshalJSON(firstEmpty)
	if !reflect.DeepEqual(data, firstEmptyJS) {
		t.Errorf("the returned unmarshalling is invalid")
	}

	empty := new(queryResponse)
	tests.ConvertToJSON(t, response, empty, cdc)

}

func TestCreateTrxChkResponse_Success(t *testing.T) {
	//variables:
	gazWanted := int64(rand.Int() % 20)
	log := "success"

	//execute:
	response, responseErr := createTrxChkResponse(true, false, gazWanted, log)
	if responseErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", responseErr.Error())
		return
	}

	retCanBeExecuted := response.CanBeExecuted()
	retCanBeAuthorized := response.CanBeAuthorized()
	retGazWanted := response.GazWanted()
	retLog := response.Log()

	if !retCanBeExecuted {
		t.Errorf("the returned canBeExecuted is exepcted to be true, false returned")
		return
	}

	if retCanBeAuthorized {
		t.Errorf("the returned canBeAuthorized is exepcted to be false, true returned")
		return
	}

	if !reflect.DeepEqual(gazWanted, retGazWanted) {
		t.Errorf("the returned gaz wanted is invalid.  Expected: %d, Returned: %d", gazWanted, retGazWanted)
		return
	}

	if !reflect.DeepEqual(log, retLog) {
		t.Errorf("the returned log is invalid.")
		return
	}

	empty := new(trxChkResponse)
	tests.ConvertToJSON(t, response, empty, cdc)
}

func TestCreateTrxResponse_Success(t *testing.T) {
	//variables:
	tags := map[string][]byte{
		"some": []byte("this is the value of a tag"),
	}

	gazUsed := int64(rand.Int() % 20)
	log := "success"

	//execute:
	response, responseErr := createTrxResponse(true, false, false, tags, gazUsed, log)
	if responseErr != nil {
		t.Errorf("the retured error was expected to be nil, error returned: %s", responseErr.Error())
		return
	}

	retIsSuccess := response.IsSuccess()
	retIsAuthorized := response.IsAuthorized()
	retHasNFS := response.HasInsufficientFunds()
	retTags := response.Tags()
	retGazUsed := response.GazUsed()
	retLog := response.Log()

	if !retIsSuccess {
		t.Errorf("the returned isSuccess is exepcted to be true, false returned")
		return
	}

	if retIsAuthorized {
		t.Errorf("the returned isAuthorized is exepcted to be false, true returned")
		return
	}

	if retHasNFS {
		t.Errorf("the returned hasNFS is exepcted to be false, true returned")
		return
	}

	if !reflect.DeepEqual(tags, retTags) {
		t.Errorf("the returned tags is invalid")
		return
	}

	if !reflect.DeepEqual(gazUsed, retGazUsed) {
		t.Errorf("the returned gazUsed is invalid.  Expected: %d, Returned: %d", gazUsed, retGazUsed)
		return
	}

	if !reflect.DeepEqual(log, retLog) {
		t.Errorf("the returned log is invalid")
		return
	}

	empty := new(trxResponse)
	tests.ConvertToJSON(t, response, empty, cdc)
}
