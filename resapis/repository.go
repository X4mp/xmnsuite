package restapis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type repository struct {
	dbPath string
}

func createRepository(dbPath string) *repository {
	out := repository{
		dbPath: dbPath,
	}

	return &out
}

// Exists returns true of the file exists, false otherwise
func (app *repository) Exists(fileName string) bool {
	fileNameWithExt := fmt.Sprintf("%s.json", fileName)
	filePath := filepath.Join(app.dbPath, fileNameWithExt)
	js, jsErr := ioutil.ReadFile(filePath)
	if jsErr != nil {
		return false
	}

	if len(js) <= 0 {
		return false
	}

	return true
}

// Retrieve retrieves a file
func (app *repository) Retrieve(fileName string, ptr interface{}) error {

	if !app.Exists(fileName) {
		str := fmt.Sprintf("the fileName (%s) does not exists", fileName)
		return errors.New(str)
	}

	fileNameWithExt := fmt.Sprintf("%s.json", fileName)
	filePath := filepath.Join(app.dbPath, fileNameWithExt)
	js, jsErr := ioutil.ReadFile(filePath)
	if jsErr != nil {
		return jsErr
	}

	unErr := json.Unmarshal(js, ptr)
	if unErr != nil {
		return unErr
	}

	return nil
}

// RetrieveNames retrieves the file names in the directory
func (app *repository) RetrieveNames() ([]string, error) {
	filesInfo, fileErr := ioutil.ReadDir(app.dbPath)
	if fileErr != nil {
		return nil, fileErr
	}

	out := []string{}
	for _, oneFileInfo := range filesInfo {
		if oneFileInfo.IsDir() {
			continue
		}

		name := oneFileInfo.Name()
		ext := filepath.Ext(name)
		nameWithoutExt := name[0 : len(name)-len(ext)]
		out = append(out, nameWithoutExt)
	}

	return out, nil
}
