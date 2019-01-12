package completed

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	prev_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

type repository struct {
	metaData         entity.MetaData
	entityRepository entity.Repository
}

func createRepository(entityRepository entity.Repository, metaData entity.MetaData) Repository {
	out := repository{
		metaData:         metaData,
		entityRepository: entityRepository,
	}

	return &out
}

// RetrieveByID retrieves a completed request by ID
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

// RetrieveByRequest retrieves a completed request by request
func (app *repository) RetrieveByRequest(req prev_request.Request) (Request, error) {
	keynames := []string{
		retrieveAllRequestsKeyname(),
		retrieveRequestByRequestKeyname(req),
	}

	ins, insErr := app.entityRepository.RetrieveByIntersectKeynames(app.metaData, keynames)
	if insErr != nil {
		return nil, insErr
	}

	if req, ok := ins.(Request); ok {
		return req, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Request instance", ins.ID().String())
	return nil, errors.New(str)
}
