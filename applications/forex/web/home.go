package web

type homeWalletList struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Wallets     []*homeWallet
}

type homeWallet struct {
	ID              string
	Creator         string
	ConcensusNeeded int
	TokenAmount     int
}

type homeGenesis struct {
	ID                    string
	GazPricePerKb         int
	ConcensusNeeded       int
	MaxAmountOfValidators int
	UserID                string
	DepositID             string
}

type home struct {
	Genesis  *homeGenesis
	WalletPS *homeWalletList
}
