package seed

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/link"
)

// Seed represents a seed on a blockchain link
type Seed interface {
	ID() *uuid.UUID
	Link() link.Link
	IP() net.IP
	Port() int
}

// Normalized represents a normalized seed
type Normalized interface {
}

// Repository represents the seed repository
type Repository interface {
	RetrieveByLink(lnk link.Link) ([]Seed, error)
}

// CreateParams represents the Create params
type CreateParams struct {
	ID   *uuid.UUID
	Link link.Link
	IP   net.IP
	Port int
}

// SDKFunc represents the Link SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Seed
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	Create: func(params CreateParams) Seed {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createSeed(params.ID, params.Link, params.IP, params.Port)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if seed, ok := ins.(Seed); ok {
					out := createStorableSeed(seed)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Seed instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllSeedsKeyname(),
				}, nil
			},
		})
	},
}
