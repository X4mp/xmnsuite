package configs

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type service struct {
}

func createService() Service {
	out := service{}
	return &out
}

// Save save an encrypted Config file on disk
func (app *service) Save(ins Configs, filePath string, password string, retypedPassword string) error {
	encrypted, encryptedErr := encrypt(ins, password, retypedPassword)
	if encryptedErr != nil {
		return encryptedErr
	}

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
