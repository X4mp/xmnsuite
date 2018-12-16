package genesis

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type repository struct {
	repository  entity.Repository
	genMetaData entity.MetaData
}

func createRepository(rep entity.Repository, genMetaData entity.MetaData) Repository {
	out := repository{
		repository:  rep,
		genMetaData: genMetaData,
	}

	return &out
}

// Retrieve retrieves the genesis instance
func (app *repository) Retrieve() (Genesis, error) {
	// if there is already a Genesis instance, return an error:
	retGen, retGenErr := app.repository.RetrieveByIntersectKeynames(app.genMetaData, []string{keyname()})
	if retGenErr != nil {
		str := fmt.Sprintf("there was an error while retrieving the Genesis instance: %s", retGenErr.Error())
		return nil, errors.New(str)
	}

	if gen, ok := retGen.(Genesis); ok {
		return gen, nil
	}

	str := fmt.Sprintf("the returned entity (ID: %s) is not a valid Genesis instance", retGen.ID().String())
	return nil, errors.New(str)
}
