package link

import (
	"net"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
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

				nodes := []Node{}
				for _, oneStorableNode := range storable.Nodes {

					nodeID, nodeIDErr := uuid.FromString(oneStorableNode.ID)
					if nodeIDErr != nil {
						return nil, nodeIDErr
					}

					pubkey := new(ed25519.PubKeyEd25519)
					pubKeyErr := cdc.UnmarshalJSON([]byte(oneStorableNode.PubKey), pubkey)
					if pubKeyErr != nil {
						return nil, pubKeyErr
					}

					ip := net.ParseIP(oneStorableNode.IP)
					nodes = append(nodes, createNode(&nodeID, pubkey, oneStorableNode.Pow, ip, oneStorableNode.Port))
				}

				out := createLinkWithNodes(&id, storable.Keyname, storable.Title, storable.Description, nodes)
				return out, nil
			}

			if storable, ok := data.(*storableLink); ok {
				return fromStorableToEntity(storable)
			}

			ptr := new(storableLink)
			jsErr := cdc.UnmarshalJSON(data.([]byte), ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			return fromStorableToEntity(ptr)

		},
		EmptyStorable: new(storableLink),
	})
}
