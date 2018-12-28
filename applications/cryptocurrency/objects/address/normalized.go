package address

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
)

type normalizedAddress struct {
	ID      string            `json:"id"`
	Wallet  wallet.Normalized `json:"wallet"`
	Address string            `json:"address"`
}

func createNormalizedAddress(ins Address) (*normalizedAddress, error) {
	wal, walErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.Wallet())
	if walErr != nil {
		return nil, walErr
	}

	out := normalizedAddress{
		ID:      ins.ID().String(),
		Wallet:  wal,
		Address: ins.Address(),
	}

	return &out, nil
}
