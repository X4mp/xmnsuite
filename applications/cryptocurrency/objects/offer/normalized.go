package offer

import (
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

type normalizedOffer struct {
	ID            string             `json:"id"`
	Pledge        pledge.Normalized  `json:"pledge"`
	To            address.Normalized `json:"to"`
	Confirmations int                `json:"confirmations"`
	Amount        int                `json:"amount"`
	Price         int                `json:"price"`
	IP            string             `json:"ip_address"`
	Port          int                `json:"port"`
}

func createNormalizedOffer(ins Offer) (*normalizedOffer, error) {
	pldge, pldgeErr := pledge.SDKFunc.CreateMetaData().Normalize()(ins.Pledge())
	if pldgeErr != nil {
		return nil, pldgeErr
	}

	addr, addrErr := address.SDKFunc.CreateMetaData().Normalize()(ins.To())
	if addrErr != nil {
		return nil, addrErr
	}

	out := normalizedOffer{
		ID:            ins.ID().String(),
		Pledge:        pldge,
		To:            addr,
		Confirmations: ins.Confirmations(),
		Amount:        ins.Amount(),
		Price:         ins.Price(),
		IP:            ins.IP().String(),
		Port:          ins.Port(),
	}

	return &out, nil
}
