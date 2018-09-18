package routers

import (
	"reflect"
	"testing"

	crypto "github.com/xmnservices/xmnsuite/crypto"
	tests "github.com/xmnservices/xmnsuite/tests"
)

func createResourcePointerForTests() (ResourcePointer, crypto.PublicKey, string) {
	from := crypto.SDKFunc.GenPK().PublicKey()
	path := "/this/is/a/path"
	res := createResourcePointer(from, path)
	return res, from, path
}

func createResourcerForTests() (Resource, ResourcePointer, []byte) {
	ptr, _, _ := createResourcePointerForTests()
	data := []byte("this is some data")
	res := createResource(ptr, data)
	return res, ptr, data
}

func TestCreateResourcePointer_Success(t *testing.T) {
	//execute:
	res, from, path := createResourcePointerForTests()
	retFrom := res.From()
	retPath := res.Path()
	retHash := res.Hash()

	if !reflect.DeepEqual(retFrom, from) {
		t.Errorf("the returned from public key is invalid")
		return
	}

	if !reflect.DeepEqual(retPath, path) {
		t.Errorf("the returned path is invalid")
		return
	}

	hsh := createResourceHash(res)
	if !reflect.DeepEqual(retHash, hsh) {
		t.Errorf("the returned hash is invalid")
		return
	}

	// convert back and forth to json:
	empty := new(resourcePointer)
	tests.ConvertToJSON(t, res, empty, cdc)
}

func TestCreateResource_Success(t *testing.T) {
	//execute:
	res, resPtr, data := createResourcerForTests()
	retPtr := res.Pointer()
	retData := res.Data()
	retHash := res.Hash()

	if !reflect.DeepEqual(retPtr, resPtr) {
		t.Errorf("the returned resource pointer is invalid")
		return
	}

	if !reflect.DeepEqual(retData, data) {
		t.Errorf("the returned data is invalid")
		return
	}

	hsh := createResourceHash(res)
	if !reflect.DeepEqual(retHash, hsh) {
		t.Errorf("the returned hash is invalid")
		return
	}

	// convert back and forth to json:
	empty := new(resource)
	tests.ConvertToJSON(t, res, empty, cdc)
}
