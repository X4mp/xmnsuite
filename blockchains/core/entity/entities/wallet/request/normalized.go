package request

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/user"
)

type normalizedRequest struct {
	ID          string          `json:"id"`
	From        user.Normalized `json:"from"`
	NewEntityJS []byte          `json:"new_entity_js"`
}

func createNormalizedRequest(req Request) (*normalizedRequest, error) {
	js, jsErr := reg.FromEntityToJSON(req.New())
	if jsErr != nil {
		return nil, jsErr
	}

	from, fromErr := user.SDKFunc.CreateMetaData().Normalize()(req.From())
	if fromErr != nil {
		return nil, fromErr
	}

	out := normalizedRequest{
		ID:          req.ID().String(),
		From:        from,
		NewEntityJS: js,
	}

	return &out, nil
}
