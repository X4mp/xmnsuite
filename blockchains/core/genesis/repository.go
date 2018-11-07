package genesis

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
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
	// if there is already a gensis instance, return an error:
	retGen, retGenErr := app.repository.RetrieveByIntersectKeynames(app.genMetaData, []string{keyname()})
	if retGenErr == nil {
		str := fmt.Sprintf("an genesis instance has already been created (ID: %s)", retGen.ID().String())
		return nil, errors.New(str)
	}

	if gen, ok := retGen.(Genesis); ok {
		return gen, nil
	}

	str := fmt.Sprintf("the returned entity (ID: %s) is not a valid Genesis instance", retGen.ID().String())
	return nil, errors.New(str)
}
