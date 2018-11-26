package link

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/request/entities/node"
)

func retrieveAllLinksKeyname() string {
	return "links"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Link",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableLink) (entity.Entity, error) {
				id, idErr := uuid.FromString(storable.ID)
				if idErr != nil {
					return nil, idErr
				}

				nodes := []node.Node{}
				nodeMetaData := node.SDKFunc.CreateMetaData()
				for _, oneNodeID := range storable.NodeIDs {

					nodeID, nodeIDErr := uuid.FromString(oneNodeID)
					if nodeIDErr != nil {
						return nil, nodeIDErr
					}

					oneNode, oneNodeErr := rep.RetrieveByID(nodeMetaData, &nodeID)
					if oneNodeErr != nil {
						return nil, oneNodeErr
					}

					if nod, ok := oneNode.(node.Node); ok {
						nodes = append(nodes, nod)
						continue
					}

					str := fmt.Sprintf("there is at least one entity (ID: %s) that was expected to be a node in the link (ID: %s), but is not", nodeID.String(), id.String())
					return nil, errors.New(str)
				}

				return createLink(&id, storable.Keyname, storable.Title, storable.Description, nodes)
			}

			if storable, ok := data.(*storableLink); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(normalizedLink)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return createLinkFromNormalized(ptr)

		},
		Normalize: func(ins entity.Entity) (interface{}, error) {
			if lnk, ok := ins.(Link); ok {
				return createNormalizedLink(lnk)
			}

			str := fmt.Sprintf("the given entity (ID: %s) is not a valid Link instance", ins.ID().String())
			return nil, errors.New(str)
		},
		Denormalize: func(ins interface{}) (entity.Entity, error) {
			if normalized, ok := ins.(*normalizedLink); ok {
				return createLinkFromNormalized(normalized)
			}

			return nil, errors.New("the given instance is not a valid normalized Link instance")
		},
		EmptyNormalized: new(normalizedLink),
		EmptyStorable:   new(storableLink),
	})
}
