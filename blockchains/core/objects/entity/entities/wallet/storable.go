package wallet

type storableWallet struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Creator string `json:"creator_pubkey"`
	CNeeded int    `json:"concensus_needed"`
}

func createStoredWallet(wallet Wallet) *storableWallet {
	out := storableWallet{
		ID:      wallet.ID().String(),
		Name:    wallet.Name(),
		Creator: wallet.Creator().String(),
		CNeeded: wallet.ConcensusNeeded(),
	}

	return &out
}
