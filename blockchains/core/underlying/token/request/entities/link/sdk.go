package link

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/request/entities/node"
)

// Link represents a blockchain link
type Link interface {
	ID() *uuid.UUID
	Keyname() string
	Title() string
	Description() string
	Nodes() []node.Node
}

// Normalized represents the normalized link
type Normalized interface {
}

// SDKFunc represents the Link SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if lnk, ok := ins.(Link); ok {
					out := createStorableLink(lnk)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Link instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllLinksKeyname(),
				}, nil
			},
		})
	},
}
