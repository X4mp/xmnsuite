package restapis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type service struct {
	dbPath string
}

func createService(dbPath string) *service {
	out := service{
		dbPath: dbPath,
	}

	return &out
}

// Save saves an instance on disk
func (app *service) Save(fileName string, v interface{}) error {
	js, jsErr := json.Marshal(v)
	if jsErr != nil {
		return jsErr
	}

	fileNameWithExt := fmt.Sprintf("%s.json", fileName)
	filePath := filepath.Join(app.dbPath, fileNameWithExt)
	fileDir := filepath.Dir(filePath)

	// if the directory doesnt exists, create it:
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, os.ModePerm)
	}

	// write to the file:
	writeErr := ioutil.WriteFile(filePath, js, os.ModePerm)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
