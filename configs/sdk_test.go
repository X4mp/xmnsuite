package configs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCreate_ThenSave_ThenRetrieve_Success(t *testing.T) {
	// variables:
	dirPath := "test_files"
	filePath := filepath.Join(dirPath, "configs.xmn")
	password := "this-is-a-password"
	defer func() {
		os.RemoveAll(dirPath)
	}()

	// generate the configs:
	conf := SDKFunc.Generate()

	// create the repository and service:
	repository := createRepository()
	service := SDKFunc.CreateService()

	// save:
	saveErr := service.Save(conf, filePath, password, password)
	if saveErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", saveErr.Error())
		return
	}

	// retrieve:
	retConf, retConfErr := repository.Retrieve(filePath, password)
	if retConfErr != nil {
		t.Errorf("the returned error was expected to be nil, error returned: %s", retConfErr.Error())
		return
	}

	// convert to storable:
	storable := createStorableConfigs(conf)
	retStorable := createStorableConfigs(retConf)

	// compare:
	if !reflect.DeepEqual(storable, retStorable) {
		t.Errorf("the saved config file did not matched the original config file")
		return
	}
}
