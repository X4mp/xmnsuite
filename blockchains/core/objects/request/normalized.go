package request

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

type normalizedRequest struct {
	ID          string             `json:"id"`
	From        user.Normalized    `json:"from"`
	NewEntityJS []byte             `json:"new_entity_js"`
	Reason      string             `json:"reason"`
	Keyname     keyname.Normalized `json:"keyname"`
}

func createNormalizedRequest(req Request) (*normalizedRequest, error) {
	js, jsErr := reg.fromEntityToJSON(req.New(), req.Keyname().Name())
	if jsErr != nil {
		return nil, jsErr
	}

	from, fromErr := user.SDKFunc.CreateMetaData().Normalize()(req.From())
	if fromErr != nil {
		return nil, fromErr
	}

	kname, knameErr := keyname.SDKFunc.CreateMetaData().Normalize()(req.Keyname())
	if knameErr != nil {
		return nil, knameErr
	}

	out := normalizedRequest{
		ID:          req.ID().String(),
		From:        from,
		NewEntityJS: js,
		Reason:      req.Reason(),
		Keyname:     kname,
	}

	return &out, nil
}
