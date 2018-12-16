package seed

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
)

type normalizedSeed struct {
	ID   string          `json:"id"`
	Link link.Normalized `json:"link"`
	IP   string          `json:"ip"`
	Port int             `json:"port"`
}

func createNormalizedSeed(ins Seed) (*normalizedSeed, error) {
	lnkIns, lnkInsErr := link.SDKFunc.CreateMetaData().Normalize()(ins.Link())
	if lnkInsErr != nil {
		return nil, lnkInsErr
	}

	out := normalizedSeed{
		ID:   ins.ID().String(),
		Link: lnkIns,
		IP:   ins.IP().String(),
		Port: ins.Port(),
	}

	return &out, nil
}
