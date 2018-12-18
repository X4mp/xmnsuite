package banks

import (
	"github.com/gorilla/mux"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
)

// ShowBankParams represents the show bank params
type ShowBankParams struct {
	Router      *mux.Router
	TemplateDir string
}

// NewBankFormParams represents the new bank form params
type NewBankFormParams struct {
	Router             *mux.Router
	TemplateDir        string
	CurrencyRepository currency.Repository
}

// SDKFunc represents the banks controllers SDK func
var SDKFunc = struct {
	ShowBanks   func(params ShowBankParams) *mux.Route
	NewBankForm func(params NewBankFormParams) *mux.Route
}{
	ShowBanks: func(params ShowBankParams) *mux.Route {
		return showBanks(params.Router, params.TemplateDir)
	},
	NewBankForm: func(params NewBankFormParams) *mux.Route {
		return newBankForm(params.Router, params.TemplateDir, params.CurrencyRepository)
	},
}
