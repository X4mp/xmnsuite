package buy

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/sell"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/external"
)

// Buy represents a buy order
type Buy interface {
	ID() *uuid.UUID
	Sell() sell.Sell
	Transfer() external.External
}

// Repository represents the buy repository
type Repository interface {
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// Service represents the buy service
type Service interface {
	Save(ins Buy) error
}

// CreateParams represents the create params
type CreateParams struct {
	ID       *uuid.UUID
	Sell     sell.Sell
	Transfer external.External
}

// SDKFunc represents the Buy SDK func
var SDKFunc = struct {
	Create func(params CreateParams) Buy
}{
	Create: func(params CreateParams) Buy {
		return nil
	},
}
