package validator

import (
	"encoding/hex"

	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/pledge"
)

type normalizedValidator struct {
	ID     string            `json:"id"`
	IP     string            `json:"ip"`
	Port   int               `json:"port"`
	PubKey string            `json:"pubkey"`
	Pledge pledge.Normalized `json:"pledge"`
}

func createNormalizedValidator(ins Validator) (*normalizedValidator, error) {
	pldge, pldgeErr := pledge.SDKFunc.CreateMetaData().Normalize()(ins.Pledge())
	if pldgeErr != nil {
		return nil, pldgeErr
	}

	out := normalizedValidator{
		ID:     ins.ID().String(),
		IP:     ins.IP().String(),
		PubKey: hex.EncodeToString(ins.PubKey().Bytes()),
		Pledge: pldge,
	}

	return &out, nil
}
