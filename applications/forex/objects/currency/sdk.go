package currency

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/crypto"
)

const (
	amountOfCharactersForSymbol         = 6
	maxAountOfCharactersForName         = 50
	maxAmountOfCharactersForDescription = 500
)

// Currency represents a currency
type Currency interface {
	ID() *uuid.UUID
	Category() category.Category
	Symbol() string
	Name() string
	Description() string
}

// Repository represents the currency repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Currency, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
}

// Normalized represents a normalized currency
type Normalized interface {
}

// Data represents the currency data
type Data struct {
	ID          string
	Category    *category.Data
	Symbol      string
	Name        string
	Description string
}

// DataSet represents the human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Currencies  []*Data
}

// DataNew represents the human-readable data for new section
type DataNew struct {
	Categories *category.DataSet
	Balances   []*balance.Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID          *uuid.UUID
	Category    category.Category
	Symbol      string
	Name        string
	Description string
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// RouteSetParams represents the RouteSet params
type RouteSetParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// RouteParams represents the Route params
type RouteParams struct {
	Tmpl             *template.Template
	EntityRepository entity.Repository
}

// RouteNewParams represents the RouteNew params
type RouteNewParams struct {
	PK                      crypto.PrivateKey
	Client                  applications.Client
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// SDKFunc represents the Currency SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Currency
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(curr Currency) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteSet             func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request)
	Route                func(params RouteParams) func(w http.ResponseWriter, r *http.Request)
	RouteNew             func(params RouteNewParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Currency {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createCurrency(params.ID, params.Category, params.Symbol, params.Name, params.Description)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return createRepresentation()
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(metaData, params.EntityRepository)
		return out
	},
	ToData: func(curr Currency) *Data {
		return toData(curr)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	RouteSet: func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create the metadata:
			metaData := createMetaData()

			// create the repositories:
			currencyRepository := createRepository(metaData, params.EntityRepository)

			// retrieve the currencies:
			currPS, currPSErr := currencyRepository.RetrieveSet(0, params.AmountOfElementsPerList)
			if currPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving the currencies entity set: %s", currPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// convert to data:
			dataSet, dataSetErr := toDataSet(currPS)
			if dataSetErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while convert a category partial set to data: %s", dataSetErr.Error())
				w.Write([]byte(str))
				return
			}

			// render:
			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, dataSet)
			return
		}
	},
	Route: func(params RouteParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			if currencyIDAsString, ok := vars["id"]; ok {
				currencyID, currencyIDErr := uuid.FromString(currencyIDAsString)
				if currencyIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the currency ID (%s) is invalid", currencyIDAsString)
					w.Write([]byte(str))
					return
				}

				// create the metadata:
				metaData := createMetaData()

				// create the repositories:
				currencyRepository := createRepository(metaData, params.EntityRepository)

				// retrieve the currency:
				retCurrency, retCurrencyErr := currencyRepository.RetrieveByID(&currencyID)
				if retCurrencyErr != nil {
					w.WriteHeader(http.StatusNotFound)
					str := fmt.Sprintf("there was an error while retrieving the Currency (ID: %s): %s", currencyID.String(), retCurrencyErr.Error())
					w.Write([]byte(str))
					return
				}

				// render:
				w.WriteHeader(http.StatusOK)
				params.Tmpl.Execute(w, toData(retCurrency))
				return

			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the ID could not be found")
			w.Write([]byte(str))
		}
	},
	RouteNew: func(params RouteNewParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create the metadata:
			representation := createRepresentation()

			// create the repositories:
			genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			requestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			balanceRepository := balance.SDKFunc.CreateRepository(balance.CreateRepositoryParams{
				DepositRepository: deposit.SDKFunc.CreateRepository(deposit.CreateRepositoryParams{
					EntityRepository: params.EntityRepository,
				}),
				WithdrawalRepository: withdrawal.SDKFunc.CreateRepository(withdrawal.CreateRepositoryParams{
					EntityRepository: params.EntityRepository,
				}),
			})

			categoryRepository := category.SDKFunc.CreateRepository(category.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			// create the services:
			requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
				PK:          params.PK,
				Client:      params.Client,
				RoutePrefix: "",
			})

			if parseFormErr := r.ParseForm(); parseFormErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while parsing form elements: %s", parseFormErr.Error())
				w.Write([]byte(str))
				return
			}

			fromWalletIDAsString := r.FormValue("fromwalletid")
			categoryIDAsString := r.FormValue("categoryid")
			symbol := r.FormValue("symbol")
			name := r.FormValue("name")
			description := r.FormValue("description")
			reason := r.FormValue("reason")
			if fromWalletIDAsString != "" && categoryIDAsString != "" && symbol != "" && name != "" && description != "" {
				// parse the walletID:
				frmWalletID, frmWalletIDErr := uuid.FromString(fromWalletIDAsString)
				if frmWalletIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given WalletID (ID: %s) is invalid: %s", fromWalletIDAsString, frmWalletIDErr.Error())
					w.Write([]byte(str))
					return
				}

				// parse the categoryID:
				catID, catIDErr := uuid.FromString(categoryIDAsString)
				if catIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given CategoryID (ID: %s) is invalid: %s", categoryIDAsString, catIDErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the wallet:
				wal, walErr := walletRepository.RetrieveByID(&frmWalletID)
				if walErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given Wallet (ID: %s) could not be retrieved: %s", frmWalletID.String(), walErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the category:
				cat, catErr := categoryRepository.RetrieveByID(&catID)
				if catErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given Category (ID: %s) could not be retrieved: %s", catID.String(), catErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the user:
				usr, usrErr := userRepository.RetrieveByPubKeyAndWallet(params.PK.PublicKey(), wal)
				if usrErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the user (Pubkey: %s, WalletID: %s): %s", params.PK.PublicKey().String(), wal.ID().String(), usrErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the keyname:
				kname, knameErr := keynameRepository.RetrieveByName(representation.MetaData().Keyname())
				if knameErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving a keyname: %s", knameErr.Error())
					w.Write([]byte(str))
					return
				}

				// create the new currency instance:
				id := uuid.NewV4()
				curr, currErr := createCurrency(&id, cat, symbol, name, description)
				if currErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while creating the Currency instance: %s", currErr.Error())
					w.Write([]byte(str))
					return
				}

				// create the request:
				req := request.SDKFunc.Create(request.CreateParams{
					FromUser:  usr,
					NewEntity: curr,
					Reason:    reason,
					Keyname:   kname,
				})

				// save the request:
				saveErr := requestService.Save(req, representation)
				if saveErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while saving a Request (Currency instance): %s", saveErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the active request:
				activeRequest, activeRequestErr := requestRepository.RetrieveByRequest(req)
				if activeRequestErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving an ActiveRequest from Request (ID: %s): %s", req.ID().String(), activeRequestErr.Error())
					w.Write([]byte(str))
					return
				}

				// redirect:
				uri := fmt.Sprintf("/requests/%s/%s/%s", kname.Group().Name(), kname.Name(), activeRequest.ID().String())
				http.Redirect(w, r, uri, http.StatusPermanentRedirect)
				return
			}

			// retrieve the genesis:
			gen, genErr := genesisRepository.Retrieve()
			if genErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
				w.Write([]byte(str))
				return
			}

			// retrieve the categories:
			catPS, catPSErr := categoryRepository.RetrieveSet(0, -1)
			if catPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving category partial set: %s", catPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// retrieve the users associated with our conf PK:
			usrPS, usrPSErr := userRepository.RetrieveSetByPubKey(params.PK.PublicKey(), 0, -1)
			if usrPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving the users entity set from creator's public key (PubKey: %s): %s", params.PK.PublicKey().String(), usrPSErr.Error())
				w.Write([]byte(str))
				return
			}

			usrsIns := usrPS.Instances()
			balances := []*balance.Data{}
			for _, oneUserIns := range usrsIns {
				if usr, ok := oneUserIns.(user.User); ok {
					bal, balErr := balanceRepository.RetrieveByWalletAndToken(usr.Wallet(), gen.Deposit().Token())
					if balErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while retrieving the balance (WalletID: %s, TokenID: %s): %s", usr.Wallet().ID().String(), gen.Deposit().Token().ID().String(), balErr.Error())
						w.Write([]byte(str))
						return
					}

					balances = append(balances, balance.SDKFunc.ToData(bal))
				}
			}

			// render:
			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, &DataNew{
				Categories: category.SDKFunc.ToDataSet(catPS),
				Balances:   balances,
			})

			return

		}
	},
}
