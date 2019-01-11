package feature

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/proposal/project"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/datastore"
)

func retrieveAllFeatureKeyname() string {
	return "features"
}

func retrieveFeatureByProjectKeyname(proj project.Project) string {
	base := retrieveAllFeatureKeyname()
	return fmt.Sprintf("%s:by_project_id:%s", base, proj.ID().String())
}

func retrieveFeatureByCreatedByUserKeyname(createdBy user.User) string {
	base := retrieveAllFeatureKeyname()
	return fmt.Sprintf("%s:by_created_by_user_id:%s", base, createdBy.ID().String())
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Feature",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableFeature); ok {
				return createFeatureFromStorable(storable, rep)
			}

			if dataAsBytes, ok := data.([]byte); ok {
				ptr := new(normalizedFeature)
				jsErr := cdc.UnmarshalJSON(dataAsBytes, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				return createFeatureFromNormalized(ptr)
			}

			str := fmt.Sprintf("the given data does not represent a Feature instance: %s", data)
			return nil, errors.New(str)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if feat, ok := ins.(Feature); ok {
				out, outErr := createNormalizedFeature(feat)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Feature instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedFeature); ok {
				return createFeatureFromNormalized(normalized)
			}

			return nil, errors.New("the given normalized instance cannot be converted to a Feature instance")
		},
		EmptyStorable:   new(storableFeature),
		EmptyNormalized: new(normalizedFeature),
	})
}

func representation() entity.Representation {
	return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
		Met: createMetaData(),
		ToStorable: func(ins entity.Entity) (interface{}, error) {
			if feat, ok := ins.(Feature); ok {
				out := createStorableFeature(feat)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Feature instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Keynames: func(ins entity.Entity) ([]string, error) {
			if feat, ok := ins.(Feature); ok {
				return []string{
					retrieveAllFeatureKeyname(),
					retrieveFeatureByProjectKeyname(feat.Project()),
					retrieveFeatureByCreatedByUserKeyname(feat.CreatedBy()),
				}, nil
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Feature instance", ins.ID().String())
			return nil, errors.New(str)
		},
		OnSave: func(ds datastore.DataStore, ins entity.Entity) error {
			if feat, ok := ins.(Feature); ok {
				// crate metadata and representation:
				userMetaData := user.SDKFunc.CreateMetaData()
				projectMetaData := project.SDKFunc.CreateMetaData()

				// create the repository and service:
				entityRepository := entity.SDKFunc.CreateRepository(ds)

				// make sure the project exists:
				_, retProjErr := entityRepository.RetrieveByID(projectMetaData, feat.Project().ID())
				if retProjErr != nil {
					str := fmt.Sprintf("the given feature (ID: %s) contains a project (ID: %s) that does not exists", feat.ID().String(), feat.Project().ID().String())
					return errors.New(str)
				}

				// make sure the user exists:
				_, retUserErr := entityRepository.RetrieveByID(userMetaData, feat.CreatedBy().ID())
				if retUserErr != nil {
					str := fmt.Sprintf("the given feature (ID: %s) contains a createdBy user (ID: %s) that does not exists", feat.ID().String(), feat.CreatedBy().ID().String())
					return errors.New(str)
				}

				// everything is alright:
				return nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Feature instance", ins.ID().String())
			return errors.New(str)
		},
	})
}
