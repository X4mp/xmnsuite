package sell

import (
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/external"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
)

type sell struct {
	UUID        *uuid.UUID        `json:"id"`
	Frm         pledge.Pledge     `json:"from"`
	Wsh         Wish              `json:"wish"`
	DepToWallet external.External `json:"deposit_to"`
}

func createSell(id *uuid.UUID, from pledge.Pledge, wish Wish, depositToWallet external.External) Sell {
	out := sell{
		UUID:        id,
		Frm:         from,
		Wsh:         wish,
		DepToWallet: depositToWallet,
	}

	return &out
}

// ID returns the ID
func (obj *sell) ID() *uuid.UUID {
	return obj.UUID
}

// From returns the from pledge
func (obj *sell) From() pledge.Pledge {
	return obj.Frm
}

// Wish returns the wish
func (obj *sell) Wish() Wish {
	return obj.Wsh
}

// DepositTo returns the deposit to external wallet
func (obj *sell) DepositToWallet() external.External {
	return obj.DepToWallet
}
