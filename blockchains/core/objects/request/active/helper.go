package active

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	core_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

func retrieveAllRequestsKeyname() string {
	return "active_requests"
}

func retrieveRequestsByRequestKeyname(req core_request.Request) string {
	base := retrieveAllRequestsKeyname()
	return fmt.Sprintf("%s:by_request_id:%s", base, req.ID().String())
}

func retrieveRequestsFromUserKeyname(usr user.User) string {
	base := retrieveAllRequestsKeyname()
	return fmt.Sprintf("%s:by_from_id:%s", base, usr.ID().String())
}

func retrieveRequestsByKeynameKeyname(kname keyname.Keyname) string {
	base := retrieveAllRequestsKeyname()
	return fmt.Sprintf("%s:by_keyname_id:%s", base, kname.ID().String())
}

func retrieveRequestsByWalletKeyname(wal wallet.Wallet) string {
	base := retrieveAllRequestsKeyname()
	return fmt.Sprintf("%s:by_wallet_id:%s", base, wal.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "ActiveRequest",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableRequest); ok {
				return createRequestFromStorable(rep, storable)
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
				return createNormalizedRequest(req)
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
