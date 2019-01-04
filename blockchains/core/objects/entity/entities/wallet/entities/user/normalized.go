package user

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type normalizedUser struct {
	ID     string            `json:"id"`
	PubKey string            `json:"pubkey"`
	Shares int               `json:"shares"`
	Wallet wallet.Normalized `json:"wallet"`
}

func createNormalizedUser(ins User) (*normalizedUser, error) {
	wallet, walletErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.Wallet())
	if walletErr != nil {
		return nil, walletErr
	}

	out := normalizedUser{
		ID:     ins.ID().String(),
		PubKey: ins.PubKey().String(),
		Shares: ins.Shares(),
		Wallet: wallet,
	}

	return &out, nil
}

func createNormalizedUserFromParams(id string, pubKey string, shares int, wal wallet.Normalized) *normalizedUser {
	out := normalizedUser{
		ID:     id,
		PubKey: pubKey,
		Shares: shares,
		Wallet: wal,
	}

	return &out
}
