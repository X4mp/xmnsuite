package external

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/link"
)

// External represents an external resource
type External interface {
	ID() *uuid.UUID
	Link() link.Link
	ResourceID() *uuid.UUID
}

// SDKFunc represents the External SDK func
var SDKFunc = struct {
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
}{
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if ewal, ok := ins.(External); ok {
					out := createStorableEWallet(ewal)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid External instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllExternalsKeyname(),
				}, nil
			},
		})
	},
}
