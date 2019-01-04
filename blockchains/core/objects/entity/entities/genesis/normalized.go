package genesis

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
)

type normalizedGenesis struct {
	ID      string                 `json:"id"`
	Info    information.Normalized `json:"information"`
	User    user.Normalized        `json:"user"`
	Deposit deposit.Normalized     `json:"deposit"`
}

func createNormalizedGenesis(ins Genesis) (*normalizedGenesis, error) {
	normalizedDeposit, normalizedDepositErr := deposit.SDKFunc.CreateMetaData().Normalize()(ins.Deposit())
	if normalizedDepositErr != nil {
		return nil, normalizedDepositErr
	}

	normalizedUser, normalizedUserErr := user.SDKFunc.CreateMetaData().Normalize()(ins.User())
	if normalizedUserErr != nil {
		return nil, normalizedUserErr
	}

	normalizedInfo, normalizedInfoErr := information.SDKFunc.CreateMetaData().Normalize()(ins.Info())
	if normalizedInfoErr != nil {
		return nil, normalizedInfoErr
	}

	out := normalizedGenesis{
		ID:      ins.ID().String(),
		Info:    normalizedInfo,
		User:    normalizedUser,
		Deposit: normalizedDeposit,
	}

	return &out, nil
}
