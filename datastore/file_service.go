package datastore

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xmnservices/xmnsuite/helpers"
)

type fileService struct {
	dirPath string
}

func createFileService(dirPath string) Service {
	out := fileService{
		dirPath: dirPath,
	}

	return &out
}

// Save saves the datastore on disk
func (app *fileService) Save(ds DataStore, filePath string) error {
	newFilePath := filepath.Join(app.dirPath, filePath)
	data, dataErr := helpers.GetBytes(ds)
	if dataErr != nil {
		return dataErr
	}

	newDirPath := filepath.Dir(newFilePath)
	if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
		os.MkdirAll(newDirPath, os.ModePerm)
	}

	writeErr := ioutil.WriteFile(newFilePath, data, 0777)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

// Retrieve retrieves a datastore stored on disk
func (app *fileService) Retrieve(filePath string) (DataStore, error) {
	newFilePath := filepath.Join(app.dirPath, filePath)
	data, dataErr := ioutil.ReadFile(newFilePath)
	if dataErr != nil {
		return nil, dataErr
	}

	ptr := new(concreteDataStore)
	maErr := helpers.Marshal(data, ptr)
	if maErr != nil {
		return nil, maErr
	}

	return ptr, nil
}
