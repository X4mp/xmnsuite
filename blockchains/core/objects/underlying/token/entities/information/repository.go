package information

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type repository struct {
	repository  entity.Repository
	infoMetaData entity.MetaData
}

func createRepository(rep entity.Repository, infoMetaData entity.MetaData) Repository {
	out := repository{
		repository:  rep,
		infoMetaData: infoMetaData,
	}

	return &out
}

// Retrieve retrieves the information instance
func (app *repository) Retrieve() (Information, error) {
	// if there is already a Information instance, return an error:
	retInfo, retInfoErr := app.repository.RetrieveByIntersectKeynames(app.infoMetaData, []string{keyname()})
	if retInfoErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the Information instance: %s", retInfoErr.Error())
		return nil, errors.New(str)
	}

	if info, ok := retInfo.(Information); ok {
		return info, nil
	}

	str := fmt.Sprintf("the returned entity (ID: %s) is not a valid Information instance", retInfo.ID().String())
	return nil, errors.New(str)
}
