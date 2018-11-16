package wallet

type storableWallet struct {
	ID      string `json:"id"`
	CNeeded int    `json:"concensus_needed"`
}

func createStoredWallet(wallet Wallet) *storableWallet {
	out := storableWallet{
		ID:      wallet.ID().String(),
		CNeeded: wallet.ConcensusNeeded(),
	}

	return &out
}
