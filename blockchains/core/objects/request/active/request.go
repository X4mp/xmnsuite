package active

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	core_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

type request struct {
	UUID    *uuid.UUID           `json:"id"`
	Req     core_request.Request `json:"request"`
	CNeeded int                  `json:"concensus_needed"`
}

func createRequest(id *uuid.UUID, req core_request.Request, concensusNeeded int) (Request, error) {
	out := request{
		UUID:    id,
		Req:     req,
		CNeeded: concensusNeeded,
	}

	return &out, nil
}

func createRequestFromNormalized(normalized *normalizedRequest) (Request, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	reqIns, reqInsErr := core_request.SDKFunc.CreateMetaData().Denormalize()(normalized.Request)
	if reqInsErr != nil {
		return nil, reqInsErr
	}

	if req, ok := reqIns.(core_request.Request); ok {
		return createRequest(&id, req, normalized.ConcensusNeeded)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Request instance", reqIns.ID().String())
	return nil, errors.New(str)

}

func createRequestFromStorable(rep entity.Repository, storable *storableRequest) (Request, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	requestID, requestIDErr := uuid.FromString(storable.RequestID)
	if requestIDErr != nil {
		return nil, requestIDErr
	}

	reqIns, reqInsErr := rep.RetrieveByID(core_request.SDKFunc.CreateMetaData(), &requestID)
	if reqInsErr != nil {
		return nil, reqInsErr
	}

	if req, ok := reqIns.(core_request.Request); ok {
		return createRequest(&id, req, storable.ConcensusNeeded)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Request instance", reqIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *request) ID() *uuid.UUID {
	return obj.UUID
}

// Request returns the request
func (obj *request) Request() core_request.Request {
	return obj.Req
}

// ConcensusNeeded returns the concensus needed
func (obj *request) ConcensusNeeded() int {
	return obj.CNeeded
}
