package group

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type repository struct {
	entityRepository entity.Repository
	metaData         entity.MetaData
}

func createRepository(entityRepository entity.Repository, metaData entity.MetaData) Repository {
	out := repository{
		entityRepository: entityRepository,
		metaData:         metaData,
	}

	return &out
}

// RetrieveByName returns a group by name
func (app *repository) RetrieveByName(name string) (Group, error) {
	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, []string{
		retrieveAllGroupsKeyname(),
		retrieveGroupByNameKeyname(name),
	})

	if insErr != nil {
		return nil, insErr
	}

	if grp, ok := ins.(Group); ok {
		return grp, nil
	}

	str := fmt.Sprintf("the given entity (ID: %s) is not a valid Group instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet returns a group set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keyname := retrieveAllGroupsKeyname()
	insPS, insPSErr := app.entityRepository.RetrieveSetByKeyname(app.metaData, keyname, index, amount)
	if insPSErr != nil {
		return nil, insPSErr
	}

	return insPS, nil
}
