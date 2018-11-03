package user

type storableUser struct {
	ID       string `json:"id"`
	PubKey   string `json:"pubkey"`
	Shares   int    `json:"shares"`
	WalletID string `json:"wallet_id"`
}

func createStorableUser(usr User) *storableUser {
	out := storableUser{
		ID:       usr.ID().String(),
		PubKey:   usr.PubKey().String(),
		Shares:   usr.Shares(),
		WalletID: usr.Wallet().ID().String(),
	}

	return &out
}
