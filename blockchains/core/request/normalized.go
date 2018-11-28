package request

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
)

type normalizedRequest struct {
	ID            string          `json:"id"`
	From          user.Normalized `json:"from"`
	NewEntityJS   []byte          `json:"new_entity_js"`
	NewEntityName string          `json:"new_entity_name"`
}

func createNormalizedRequest(req Request) (*normalizedRequest, error) {
	js, jsErr := reg.fromEntityToJSON(req.New(), req.NewName())
	if jsErr != nil {
		return nil, jsErr
	}

	from, fromErr := user.SDKFunc.CreateMetaData().Normalize()(req.From())
	if fromErr != nil {
		return nil, fromErr
	}

	out := normalizedRequest{
		ID:            req.ID().String(),
		From:          from,
		NewEntityJS:   js,
		NewEntityName: req.NewName(),
	}

	return &out, nil
}
