package datastore

import (
	"bytes"
	"os"
	"testing"
)

func TestSaveThenRetrieve_Success(t *testing.T) {
	//variables:
	dirPath := "test_files"
	relFilePath := "db.xmnds"
	defer func() {
		os.RemoveAll(dirPath)
	}()

	// create datastore:
	ds := createConcreteDataStore()

	// create service:
	service := createFileService(dirPath)

	// add some data:
	ds.Keys().Save("some_data", "this is some data")
	ds.Keys().Save("other_data", "this is some other data")

	// save:
	saveErr := service.Save(ds, relFilePath)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	// retrieve:
	retDS, retDSErr := service.Retrieve(relFilePath)
	if retDSErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retDSErr.Error())
		return
	}

	if bytes.Compare(ds.Head().Head().Get(), retDS.Head().Head().Get()) != 0 {
		t.Errorf("the retrieved datastore is invalid")
		return
	}

}
