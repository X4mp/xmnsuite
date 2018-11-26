package node

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

// Node represents a node on a blockchain link
type Node interface {
	ID() *uuid.UUID
	Power() int
	IP() net.IP
	Port() int
}

// Normalized represents a normalized node
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID    *uuid.UUID
	Power int
	IP    net.IP
	Port  int
}

// SDKFunc represents the Link SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Node
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Node {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createNode(params.ID, params.Power, params.IP, params.Port)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if node, ok := ins.(Node); ok {
					out := createStorableNode(node)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Node instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllNodesKeyname(),
				}, nil
			},
		})
	},
}
