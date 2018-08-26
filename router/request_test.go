package router

import (
	"math/rand"
	"reflect"
	"testing"

	tests "github.com/XMNBlockchain/datamint/tests"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

func TestCreateRequest_Success(t *testing.T) {
	//variables:
	pk := ed25519.GenPrivKey()
	from := pk.PubKey()
	path := "/this/is/a/path"
	data := []byte("this is some data")

	str := requestSignedStruct{
		Path: path,
		Data: data,
	}

	js, _ := cdc.MarshalJSON(str)
	sig, _ := pk.Sign(js)

	//execute:
	req, reqErr := createRequest(from, path, data, sig)
	if reqErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", reqErr.Error())
		return
	}

	retFrom := req.From()
	retPath := req.Path()
	retData := req.Data()
	retSig := req.Signature()

	if !reflect.DeepEqual(from, retFrom) {
		t.Errorf("the returned from is invalid")
		return
	}

	if !reflect.DeepEqual(path, retPath) {
		t.Errorf("the returned path is invalid")
		return
	}

	if !reflect.DeepEqual(data, retData) {
		t.Errorf("the returned data is invalid")
		return
	}

	if !reflect.DeepEqual(sig, retSig) {
		t.Errorf("the returned data is invalid")
		return
	}

	empty := new(request)
	tests.ConvertToJSON(t, req, empty, cdc)

}

func TestCreateTrxChkRequest_Success(t *testing.T) {
	//variables:
	pk := ed25519.GenPrivKey()
	from := pk.PubKey()
	path := "/this/is/a/path"
	dataSizeInBytes := int64(rand.Int() % 20)

	str := requestTrxChkSignedStruct{
		Path:           path,
		DtaSizeInBytes: dataSizeInBytes,
	}

	js, _ := cdc.MarshalJSON(str)

	sig, _ := pk.Sign(js)

	//execute:
	req, reqErr := createTrxChkRequest(from, path, dataSizeInBytes, sig)
	if reqErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", reqErr.Error())
		return
	}

	retFrom := req.From()
	retPath := req.Path()
	retDataSizeInBytes := req.DataSizeInBytes()
	retSig := req.Signature()

	if !reflect.DeepEqual(from, retFrom) {
		t.Errorf("the returned from is invalid")
		return
	}

	if !reflect.DeepEqual(path, retPath) {
		t.Errorf("the returned path is invalid")
		return
	}

	if !reflect.DeepEqual(dataSizeInBytes, retDataSizeInBytes) {
		t.Errorf("the returned dataSizeInBytes is invalid")
		return
	}

	if !reflect.DeepEqual(sig, retSig) {
		t.Errorf("the returned data is invalid")
		return
	}

	empty := new(trxChkRequest)
	tests.ConvertToJSON(t, req, empty, cdc)
}
