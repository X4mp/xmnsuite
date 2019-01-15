package completed

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	prev_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

type request struct {
	UUID    *uuid.UUID           `json:"id"`
	Req     prev_request.Request `json:"request"`
	CNeeded int                  `json:"concensus_needed"`
	Appr    int                  `json:"approved"`
	DisAppr int                  `json:"disapproved"`
	Neutrl  int                  `json:"neutral"`
}

func createRequest(id *uuid.UUID, req prev_request.Request, concensusNeeded int, approved int, disapproved int, neutral int) (Request, error) {
	out := request{
		UUID:    id,
		Req:     req,
		CNeeded: concensusNeeded,
		Appr:    approved,
		DisAppr: disapproved,
		Neutrl:  neutral,
	}

	return &out, nil
}

func createRequestFromNormalized(normalized *normalizedRequest) (Request, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	reqIns, reqInsErr := prev_request.SDKFunc.CreateMetaData().Denormalize()(normalized.Request)
	if reqInsErr != nil {
		return nil, reqInsErr
	}

	if req, ok := reqIns.(prev_request.Request); ok {
		return createRequest(&id, req, normalized.ConcensusNeeded, normalized.Approved, normalized.Disapproved, normalized.Neutral)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Request instance", reqIns.ID().String())
	return nil, errors.New(str)
}

func createRequestFromStorable(storable *storableRequest, rep entity.Repository) (Request, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	reqID, reqIDErr := uuid.FromString(storable.RequestID)
	if reqIDErr != nil {
		return nil, reqIDErr
	}

	reqIns, reqInsErr := rep.RetrieveByID(prev_request.SDKFunc.CreateMetaData(), &reqID)
	if reqInsErr != nil {
		return nil, reqInsErr
	}

	if req, ok := reqIns.(prev_request.Request); ok {
		return createRequest(&id, req, storable.ConcensusNeeded, storable.Approved, storable.Disapproved, storable.Neutral)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Request instance", reqIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *request) ID() *uuid.UUID {
	return obj.UUID
}

// Request returns the request
func (obj *request) Request() prev_request.Request {
	return obj.Req
}

// ConcensusNeeded returns the concensus needed
func (obj *request) ConcensusNeeded() int {
	return obj.CNeeded
}

// Approved returns the amount of approved votes
func (obj *request) Approved() int {
	return obj.Appr
}

// Disapproved returns the amount of disapproved votes
func (obj *request) Disapproved() int {
	return obj.DisAppr
}

// Neutral returns the amount of neutral votes
func (obj *request) Neutral() int {
	return obj.Neutrl
}
