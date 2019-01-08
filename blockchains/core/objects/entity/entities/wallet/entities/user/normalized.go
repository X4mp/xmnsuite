package user

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type normalizedUser struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	PubKey   string            `json:"pubkey"`
	Shares   int               `json:"shares"`
	Wallet   wallet.Normalized `json:"wallet"`
	Referral wallet.Normalized `json:"referral"`
}

func createNormalizedUser(ins User) (*normalizedUser, error) {
	walletMetaData := wallet.SDKFunc.CreateMetaData()
	wallet, walletErr := walletMetaData.Normalize()(ins.Wallet())
	if walletErr != nil {
		return nil, walletErr
	}

	var refUsr Normalized
	if ins.HasBeenReferred() {
		ref, refErr := walletMetaData.Normalize()(ins.Referral())
		if refErr != nil {
			return nil, refErr
		}

		refUsr = ref
	}

	out := normalizedUser{
		ID:       ins.ID().String(),
		Name:     ins.Name(),
		PubKey:   ins.PubKey().String(),
		Shares:   ins.Shares(),
		Wallet:   wallet,
		Referral: refUsr,
	}

	return &out, nil
}
