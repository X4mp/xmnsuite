package node

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/link"
)

func retrieveAllNodesKeyname() string {
	return "nodes"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Node",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			if storable, ok := data.(*storableNode); ok {

				lnkMetaData := link.SDKFunc.CreateMetaData()

				nodeID, nodeIDErr := uuid.FromString(storable.ID)
				if nodeIDErr != nil {
					return nil, nodeIDErr
				}

				lnkID, lnkIDErr := uuid.FromString(storable.LinkID)
				if lnkIDErr != nil {
					return nil, lnkIDErr
				}

				lnkIns, lnkInsErr := rep.RetrieveByID(lnkMetaData, &lnkID)
				if lnkInsErr != nil {
					return nil, lnkInsErr
				}

				if lnk, ok := lnkIns.(link.Link); ok {
					ip := net.ParseIP(storable.IP)
					out := createNode(&nodeID, lnk, storable.Pow, ip, storable.Port)
					return out, nil
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Link instance", lnkIns.ID().String())
				return nil, errors.New(str)
			}

			ptr := new(normalizedNode)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createNodeFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if node, ok := ins.(Node); ok {
				out, outErr := createNormalizedNode(node)
				if outErr != nil {
					return nil, outErr
				}

				return out, nil
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Node instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedNode); ok {
				return createNodeFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Node instance")
		},
		EmptyNormalized: new(normalizedNode),
		EmptyStorable:   new(storableNode),
	})
}
