package claim

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/framework/entity"
)

// Claim represents a transfer claim
type Claim interface {
	ID() *uuid.UUID
	Transfer() transfer.Transfer
	Deposit() deposit.Deposit
}

// Service represents the claim service
type Service interface {
	Save(ins Claim) error
}

// SDKFunc represents the Transfer SDK func
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
				if claim, ok := ins.(Claim); ok {
					out := createStorableClaim(claim)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Claim instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				return []string{
					retrieveAllClaimsKeyname(),
				}, nil
			},
		})
	},
}
