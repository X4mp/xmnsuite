package web

type homeRequestList struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Requests    []*homeRequest
}

type homeRequest struct {
	ID         string
	FromUserID string
	NewName    string
}

type homeRequestSingle struct {
	ID         string
	FromUserID string
	NewName    string
	NewJS      string
	Votes      *homeVoteList
}

type homeVoteList struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Votes       []*homeVote
}

type homeVote struct {
	ID               string
	UserVoterID      string
	UserAmountShares int
	IsApproved       bool
}

type homeCategory struct {
	ID          string
	ParentID    string
	Name        string
	Description string
}

type homeCategoryList struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Categories  []*homeCategory
}

type homeCategoryNew struct {
	Users *homeUserList
}

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

type singleWallet struct {
	ID              string
	ConcensusNeeded int
	TokenAmount     int
	Users           *homeUserList
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
	ID                     string
	GazPricePerKb          int
	GazPriceInMatrixWorkKb int
	ConcensusNeeded        int
	MaxAmountOfValidators  int
	UserID                 string
	DepositID              string
}

type home struct {
	Genesis     *homeGenesis
	WalletPS    *homeWalletList
	AllWalletPS *homeWalletList
	UserPS      *homeUserList
}
