package validator

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/core/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
)

func retrieveAllValidatorsKeyname() string {
	return "validators"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Transfer",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableValidator) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				pledgeID, pledgeIDErr := uuid.FromString(storable.PledgeID)
				if pledgeIDErr != nil {
					return nil, pledgeIDErr
				}

				pubkey := new(ed25519.PubKeyEd25519)
				pubKeyErr := cdc.UnmarshalJSON([]byte(storable.PubKey), pubkey)
				if pubKeyErr != nil {
					return nil, pubKeyErr
				}

				// retrieve the pledge:
				pledgeMetaData := pledge.SDKFunc.CreateMetaData()
				pledgeIns, pledgeInsErr := rep.RetrieveByID(pledgeMetaData, &pledgeID)
				if pledgeInsErr != nil {
					return nil, pledgeInsErr
				}

				if pldge, ok := pledgeIns.(pledge.Pledge); ok {
					out := createValidator(&id, pubkey, pldge)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pledgeID.String())
				return nil, errors.New(str)
			}

			if storable, ok := data.(*storableValidator); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableValidator)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableValidator),
	})
}
