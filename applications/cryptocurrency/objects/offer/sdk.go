package offer

import (
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

// Offer represents an offer to manage coins
type Offer interface {
	ID() *uuid.UUID
	Pledge() pledge.Pledge
	To() address.Address
	Amount() int
	Price() int
	IP() net.IP
	Port() int
}

// Normalized represents a normalized offer
type Normalized interface {
}

// Data represents human-readable data
type Data struct {
	ID        string
	Pledge    *pledge.Data
	ToAddress *address.Data
	Amount    int
	Price     int
}

// DataSet represents the human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Offers      []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	Pledge pledge.Pledge
	To     address.Address
	Amount int
	Price  int
	IP     net.IP
	Port   int
}

// SDKFunc represents the Offer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Offer
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	ToData               func(off Offer) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
}{
	Create: func(params CreateParams) Offer {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createOffer(params.ID, params.Pledge, params.To, params.Amount, params.Price, params.IP, params.Port)
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
	ToData: func(off Offer) *Data {
		return toData(off)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
