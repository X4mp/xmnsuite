package deposit

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

// Deposit represents a deposit on an offer
type Deposit interface {
	ID() *uuid.UUID
	Offer() offer.Offer
	From() address.Address
	Amount() int
}

// Normalized represents a normalized offer
type Normalized interface {
}

// Data represents human-readable data
type Data struct {
	ID     string
	Offer  *offer.Data
	From   *address.Data
	Amount int
}

// DataSet represents the human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Deposits    []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	Offer  offer.Offer
	From   address.Address
	Amount int
}

// SDKFunc represents the Deposit SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Deposit
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	ToData               func(dep Deposit) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
}{
	Create: func(params CreateParams) Deposit {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createDeposit(params.ID, params.Offer, params.From, params.Amount)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return createRepresentation()
	},
	ToData: func(dep Deposit) *Data {
		return toData(dep)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
