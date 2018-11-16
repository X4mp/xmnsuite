package sell

type storableWish struct {
	ExternalTokenID string `json:"external_token_id"`
	Amount          int    `json:"amount"`
}

func createStorableWish(ins Wish) *storableWish {
	out := storableWish{
		ExternalTokenID: ins.Token().ID().String(),
		Amount:          ins.Amount(),
	}

	return &out
}

type storableSell struct {
	ID                string        `json:"id"`
	FromPledgeID      string        `json:"from_pledge_id"`
	Wish              *storableWish `json:"wish"`
	DepositToWalletID string        `json:"deposit_to_external_wallet_id"`
}

func createStorableSell(ins Sell) *storableSell {
	out := storableSell{
		ID:                ins.ID().String(),
		FromPledgeID:      ins.From().ID().String(),
		Wish:              createStorableWish(ins.Wish()),
		DepositToWalletID: ins.DepositToWallet().ID().String(),
	}

	return &out
}
