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

type homeUserList struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Users       []*homeUser
}

type homeUser struct {
	ID       string
	Shares   int
	WalletID string
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
	Genesis     *homeGenesis
	WalletPS    *homeWalletList
	AllWalletPS *homeWalletList
	UserPS      *homeUserList
}
