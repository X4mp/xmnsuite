package web

import (
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	walletpkg "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
)

// Web represents a web server
type Web interface {
	Start() error
	Stop() error
}

// CreateParams represents the create params
type CreateParams struct {
	Port               int
	EntityService      entity.Service
	BalanceRepository  balance.Repository
	GenesisRepository  genesis.Repository
	WalletRepository   walletpkg.Repository
	CategoryRepository category.Repository
	CurrencyRepository currency.Repository
}

// SDKFunc represents the web server
var SDKFunc = struct {
	Create func(params CreateParams) Web
}{
	Create: func(params CreateParams) Web {
		out := createWeb(
			params.Port,
			params.EntityService,
			params.BalanceRepository,
			params.GenesisRepository,
			params.WalletRepository,
			params.CategoryRepository,
			params.CurrencyRepository,
		)

		return out
	},
}
