package roles

import (
	"testing"

	"github.com/xmnservices/xmnsuite/datastore/lists"
	"github.com/xmnservices/xmnsuite/helpers"
)

func TestSingle_save_thenExists_thenRetrieve_thenDelete_Success(t *testing.T) {
	//variables:
	writeOnKey := "i-want-to-write-on-this-key"

	//create a list:
	setApp := lists.SDKFunc.CreateSet()
	setApp.Add(writeOnKey, []byte("some data"))

	//create roles:
	app := createConcreteRoles()

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
