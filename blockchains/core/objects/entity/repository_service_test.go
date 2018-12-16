package entity

import (
	"testing"

	"github.com/xmnservices/xmnsuite/datastore"
)

func compareEntityForTests(t *testing.T, first Entity, second Entity) {
	if first.ID().String() != second.ID().String() {
		t.Errorf("the returned ID is invalid.  Expected: %s, returned: %s", first.ID().String(), second.ID().String())
	}
}

func TestSave_thenRetrieve_Success(t *testing.T) {
	// variables:
	ins := createTestEntityForTests()
	rep := CreateRepresentationForTests()
	met := rep.MetaData()
	store := datastore.SDKFunc.Create()

	// service + repository:
	repository := createRepository(store)
	service := createService(store, repository)

	// save:
	saveErr := service.Save(ins, rep)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	// retireve by ID:
	retIns, retInsErr := repository.RetrieveByID(met, ins.ID())
	if retInsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retInsErr.Error())
		return
	}

	// retrieving by intersecting keynames:
	keynames, _ := rep.Keynames()(ins)
	retIntersectIns, retIntersectInsErr := repository.RetrieveByIntersectKeynames(met, keynames)
	if retIntersectInsErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retIntersectInsErr.Error())
		return
	}

	// compare:
	compareEntityForTests(t, ins, retIns)
	compareEntityForTests(t, retIntersectIns, retIns)

}
