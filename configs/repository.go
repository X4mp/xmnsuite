package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	crypto "github.com/xmnservices/xmnsuite/crypto"
)

type repository struct {
}

func createRepository() Repository {
	out := repository{}
	return &out
}

// Retrieve retrieves a configs file
func (app *repository) Retrieve(filePath string, password string) (Configs, error) {
	data, dataErr := ioutil.ReadFile(filePath)
	if dataErr != nil {
		str := fmt.Sprintf("there was an error while reading the Configs file: %s", dataErr.Error())
		return nil, errors.New(str)
	}

	jsonAsBytes := crypto.SDKFunc.Decrypt(crypto.DecryptParams{
		Pass:         []byte(password),
		EncryptedMsg: string(data),
	})

	ptr := new(storableConfigs)
	jsErr := json.Unmarshal(jsonAsBytes, ptr)
	if jsErr != nil {
		str := fmt.Sprintf("the file does not contain an encrypted Configs file: %s", jsErr.Error())
		return nil, errors.New(str)
	}

	conf, confErr := fromStorableToConfigs(ptr)
	if confErr != nil {
		str := fmt.Sprintf("there was an error while converting a storable config instance to a configs instance: %s", confErr.Error())
		return nil, errors.New(str)
	}

	return conf, nil
}
