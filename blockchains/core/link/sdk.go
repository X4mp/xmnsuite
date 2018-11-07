package link

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
)

// Link represents a blockchain link
type Link interface {
	ID() *uuid.UUID
	Keyname() string
	Title() string
	Description() string
	Nodes() []Node
}

// Node represents a node on a blockchain link
type Node interface {
	ID() *uuid.UUID
	PublicKey() tcrypto.PubKey
	Power() int
	IP() net.IP
	Port() int
}

// Daemon represents the daemon
type Daemon interface {
	Start() error
	Stop() error
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
