package address

import "encoding/hex"

type storableAddress struct {
	ID       string `json:"id"`
	WalletID string `json:"wallet_id"`
	Address  string `json:"address"`
}

func createStorableAddress(ins Address) *storableAddress {
	out := storableAddress{
		ID:       ins.ID().String(),
		WalletID: ins.Wallet().ID().String(),
		Address:  hex.EncodeToString(ins.Address()),
	}

	return &out
}