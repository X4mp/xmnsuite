package link

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

type repository struct {
	rep      entity.Repository
	metaData entity.MetaData
}

func createRepository(rep entity.Repository, metaData entity.MetaData) Repository {
	out := repository{
		rep:      rep,
		metaData: metaData,
	}

	return &out
}

// RetrieveSet retrieves a set of links
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllLinksKeyname()
	retPartialSet, retPartialSetErr := app.rep.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if retPartialSetErr != nil {
		return nil, retPartialSetErr
	}

	return retPartialSet, nil
}
