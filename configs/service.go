package configs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xmnservices/xmnsuite/crypto"
)

type service struct {
}

func createService() Service {
	out := service{}
	return &out
}

// Save save an encrypted Config file on disk
func (app *service) Save(ins Configs, filePath string, password string, retypedPassword string) error {

	if len(password) < 6 {
		return errors.New("The password must contain at least 6 characters")
	}

	if password != retypedPassword {
		return errors.New("The passwords do not match")
	}

	storable := createStorableConfigs(ins)
	js, jsErr := json.Marshal(storable)
	if jsErr != nil {
		return jsErr
	}

	encrypted := crypto.SDKFunc.Encrypt(crypto.EncryptParams{
		Pass: []byte(password),
		Msg:  js,
	})

	dirPath := filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, os.ModePerm)
	}

	writeErr := ioutil.WriteFile(filePath, []byte(encrypted), os.ModePerm)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
