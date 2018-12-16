package project

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

func retrieveAllProjectsKeyname() string {
	return "projects"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Project",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableProject); ok {
				return createProjectFromStorable(storable)
			}

			ptr := new(storableProject)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createProjectFromStorable(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if proj, ok := ins.(Project); ok {
				out := createStorableProject(proj)
				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Project instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if storable, ok := ins.(*storableProject); ok {
				return createProjectFromStorable(storable)
			}

			return nil, errors.New("the given instance is not a valid normalized Project instance")
		},
		EmptyStorable:   new(storableProject),
		EmptyNormalized: new(storableProject),
	})
}
