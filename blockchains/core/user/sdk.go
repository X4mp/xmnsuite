package user

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
)

// User represents a user
type User interface {
	ID() *uuid.UUID
	PubKey() crypto.PublicKey
	Shares() int
	Wallet() wallet.Wallet
}

// SDKFunc represents the User SDK func
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
				if usr, ok := ins.(User); ok {
					out := createStorableUser(usr)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllUserKeyname(),
				}, nil
			},
		})
	},
}
