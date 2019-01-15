package balance

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
)

type balance struct {
	Wal wallet.Wallet `json:"wallet"`
	Am  int           `json:"amount"`
}

func createBalance(wal wallet.Wallet, amount int) Balance {
	out := balance{
		Wal: wal,
		Am:  amount,
	}

	return &out
}

// Wallet returns the wallet
func (obj *balance) Wallet() wallet.Wallet {
	return obj.Wal
}

// Amount returns the amount
func (obj *balance) Amount() int {
	return obj.Am
}
