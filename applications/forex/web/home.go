package web

type homeRequestGroupList struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Requests    []*homeRequestGroup
}

type homeRequestGroup struct {
	ID   string
	Name string
}

type homeRequestKeynamesOfGroup struct {
	Group    *homeRequestGroup
	Keynames *homeRequestKeynamesList
}

type homeRequestKeynamesList struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Keynames    []*homeRequestKeyname
}

type homeRequestKeyname struct {
	ID    string
	Name  string
	Group *homeRequestGroup
}

type homeRequests struct {
	Keyname  *homeRequestKeyname
	Requests *homeRequestList
}

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
	Reason     string
}

type homeRequestSingle struct {
	ID              string
	FromUserID      string
	Reason          string
	NewJS           string
	ConcensusNeeded int
	Keyname         *homeRequestKeyname
	MyUsers         *homeUserList
	Votes           *homeVoteList
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
	Reason           string
	IsNeutral        bool
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
