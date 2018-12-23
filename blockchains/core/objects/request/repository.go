package request

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
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

// RetrieveByID retrieves a request by ID
func (app *repository) RetrieveByID(id *uuid.UUID) (Request, error) {
	ins, insErr := app.entityRepository.RetrieveByID(app.metaData, id)
	if insErr != nil {
		return nil, insErr
	}

	if req, ok := ins.(Request); ok {
		return req, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Request instance", ins.ID().String())
	return nil, errors.New(str)
}

// RetrieveSet retrieves a request set
func (app *repository) RetrieveSet(index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllRequestsKeyname(),
	}
	reqPS, reqPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if reqPSErr != nil {
		return nil, reqPSErr
	}

	return reqPS, nil
}

// RetrieveSetByFromUser retrieves a request set from user
func (app *repository) RetrieveSetByFromUser(usr user.User, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllRequestsKeyname(),
		retrieveAllRequestsFromUserKeyname(usr),
	}
	reqPS, reqPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if reqPSErr != nil {
		return nil, reqPSErr
	}

	return reqPS, nil
}

// RetrieveSetByKeyname retrieves a request set by keyname
func (app *repository) RetrieveSetByKeyname(kname keyname.Keyname, index int, amount int) (entity.PartialSet, error) {
	keynames := []string{
		retrieveAllRequestsKeyname(),
		retrieveAllRequestsByKeynameKeyname(kname),
	}
	reqPS, reqPSErr := app.entityRepository.RetrieveSetByIntersectKeynames(app.metaData, keynames, index, amount)
	if reqPSErr != nil {
		return nil, reqPSErr
	}

	return reqPS, nil
}
