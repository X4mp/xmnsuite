package keyname

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
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

// RetrieveByName returns a keyname by name
func (app *repository) RetrieveByName(name string) (Keyname, error) {
	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, []string{
		retrieveAllKeynamesKeyname(),
		retrieveKeynameByNameKeyname(name),
	})

	if insErr != nil {
		return nil, insErr
	}

	if kname, ok := ins.(Keyname); ok {
		return kname, nil
	}

	str := fmt.Sprintf("the given entity (ID: %s) is not a valid Keyname instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSetByGroup returns a keyname partial set
func (app *repository) RetrieveSetByGroup(grp group.Group, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllKeynamesKeyname(),
		retrieveKeynameByGroupKeyname(grp),
	}

	knamePS, knamePSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if knamePSErr != nil {
		return nil, knamePSErr
	}

	return knamePS, nil
}
