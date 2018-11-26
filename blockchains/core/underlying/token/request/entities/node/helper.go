package node

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

func retrieveAllNodesKeyname() string {
	return "nodes"
}

func createMetaData() entity.MetaData {
	return entity.SDKFunc.CreateMetaData(entity.CreateMetaDataParams{
		Name: "Node",
		ToEntity: func(rep entity.Repository, data interface{}) (entity.Entity, error) {
			fromStorableToEntity := func(storable *storableNode) (entity.Entity, error) {
				nodeID, nodeIDErr := uuid.FromString(storable.ID)
				if nodeIDErr != nil {
					return nil, nodeIDErr
				}

				pubkey := new(ed25519.PubKeyEd25519)
				pubKeyErr := cdc.UnmarshalJSON([]byte(storable.PubKey), pubkey)
				if pubKeyErr != nil {
					return nil, pubKeyErr
				}

				ip := net.ParseIP(storable.IP)
				out := createNode(&nodeID, pubkey, storable.Pow, ip, storable.Port)
				return out, nil
			}

			if storable, ok := data.(*storableNode); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableNode)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

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
