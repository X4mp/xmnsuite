package account

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/work"
)

type normalizedAccount struct {
	User user.Normalized `json:"user"`
	Work work.Normalized `json:"Work"`
}

func createNormalizedAccount(ins Account) (*normalizedAccount, error) {
	usr, usrErr := user.SDKFunc.CreateMetaData().Normalize()(ins.User())
	if usrErr != nil {
		return nil, usrErr
	}

	wrk := work.SDKFunc.Normalize(ins.Work())
	out := normalizedAccount{
		User: usr,
		Work: wrk,
	}

	return &out, nil
}
