package node

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func retrieveAllNodesKeyname() string {
	return "nodes"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Node",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableNode); ok {
				return createNodeFromStorable(storable)
			}

			ptr := new(storableNode)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createNodeFromStorable(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if node, ok := ins.(Node); ok {
				return createStorableNode(node), nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Node instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*storableNode); ok {
				return createNodeFromStorable(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Node instance")
		},
		EmptyNormalized: new(storableNode),
		EmptyStorable:   new(storableNode),
	})
}
