package sell

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type repository struct {
	entityRepository entity.Repository
	sellMetaData     entity.MetaData
}

func createRepository(entityRepository entity.Repository, sellMetaData entity.MetaData) Repository {
	out := repository{
		entityRepository: entityRepository,
		sellMetaData:     sellMetaData,
	}

	return &out
}

// RetrieveMatch retrieves the best match that fit the wish
func (app *repository) RetrieveMatch(with Wish) (Sell, error) {
	return nil, nil
}

// RetrieveMatches retrieve the sell orders that matches the wish
func (app *repository) RetrieveMatches(wish Wish) (entity.PartialSet, error) {
	return nil, nil
}

// RetrieveSet retrieves a list of Sell orders
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllSellsKeyname()
	sellPS, sellPSErr := app.entityRepository.RetrieveSetByKeyname(app.sellMetaData, keyname, index, amount)
	if sellPSErr != nil {
		return nil, sellPSErr
	}

	return sellPS, nil
}
