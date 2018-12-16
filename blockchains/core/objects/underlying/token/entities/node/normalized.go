package node

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
)

type normalizedNode struct {
	ID   string          `json:"id"`
	Link link.Normalized `json:"link"`
	Pow  int             `json:"power"`
	IP   string          `json:"ip"`
	Port int             `json:"port"`
}

func createNormalizedNode(ins Node) (*normalizedNode, error) {
	lnkIns, lnkInsErr := link.SDKFunc.CreateMetaData().Normalize()(ins.Link())
	if lnkInsErr != nil {
		return nil, lnkInsErr
	}

	out := normalizedNode{
		ID:   ins.ID().String(),
		Link: lnkIns,
		Pow:  ins.Power(),
		IP:   ins.IP().String(),
		Port: ins.Port(),
	}

	return &out, nil
}
