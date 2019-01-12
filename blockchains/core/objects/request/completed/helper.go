package completed

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	prev_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

func retrieveAllRequestsKeyname() string {
	return "completedrequests"
}

func retrieveRequestByRequestKeyname(req prev_request.Request) string {
	base := retrieveAllRequestsKeyname()
	return fmt.Sprintf("%s:by_request_id:%s", base, req.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "CompletedRequest",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableRequest); ok {
				return createRequestFromStorable(storable, rep)
			}

			ptr := new(normalizedRequest)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createRequestFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if req, ok := ins.(Request); ok {
				out := createStorableRequest(req)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Request instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedRequest); ok {
				return createRequestFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Request instance")
		},
		EmptyStorable:   new(storableRequest),
		EmptyNormalized: new(normalizedRequest),
	})
}
