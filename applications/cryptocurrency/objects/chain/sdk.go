package chain

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

// Chain represents a chain that hold bitcoin representations held by an offerer
type Chain interface {
	ID() *uuid.UUID
	Offer() offer.Offer
	TotalAmount() int
}

// Normalized represents a normalized offer
type Normalized interface {
}

// Data represents human-readable data
type Data struct {
	ID          string
	Offer       *offer.Data
	TotalAmount int
}

// DataSet represents the human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Chains      []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID          *uuid.UUID
	Offer       offer.Offer
	TotalAmount int
}

// SDKFunc represents the Chain SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Chain
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	ToData               func(chn Chain) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
}{
	Create: func(params CreateParams) Chain {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createChain(params.ID, params.Offer, params.TotalAmount)
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
	ToData: func(chn Chain) *Data {
		return toData(chn)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
