package balance

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Balance represents a wallet balance
type Balance interface {
	On() wallet.Wallet
	Of() token.Token
	Amount() int
}

// Repository represents a balance repository
type Repository interface {
	RetrieveByWalletAndToken(wal wallet.Wallet, tok token.Token) (Balance, error)
}

// Data represents human-readable data
type Data struct {
	On     *wallet.Data
	Of     *token.Data
	Amount int
}

// DataSet represents human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Balances    []*Data
}

// CreateRepositoryParams represents a CreateRepository params
type CreateRepositoryParams struct {
	Datastore            datastore.DataStore
	DepositRepository    deposit.Repository
	WithdrawalRepository withdrawal.Repository
}

// RouteListParams represents the route list params
type RouteListParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// RouteParams represents the route params
type RouteParams struct {
	Tmpl             *template.Template
	EntityRepository entity.Repository
}

// SDKFunc represents the balance SDK func
var SDKFunc = struct {
	CreateRepository func(params CreateRepositoryParams) Repository
	ToData           func(bal Balance) *Data
	RouteList        func(params RouteListParams) func(w http.ResponseWriter, r *http.Request)
	Route            func(params RouteParams) func(w http.ResponseWriter, r *http.Request)
}{
	CreateRepository: func(params CreateRepositoryParams) Repository {
		if params.Datastore != nil {
			if params.Datastore != nil {
				depositRepository := deposit.SDKFunc.CreateRepository(deposit.CreateRepositoryParams{
					Datastore: params.Datastore,
				})

				withdrawalRepository := withdrawal.SDKFunc.CreateRepository(withdrawal.CreateRepositoryParams{
					Datastore: params.Datastore,
				})

				out := createRepository(depositRepository, withdrawalRepository)
				return out
			}

			out := createRepository(params.DepositRepository, params.WithdrawalRepository)
			return out
		}

		out := createRepository(params.DepositRepository, params.WithdrawalRepository)
		return out
	},
	ToData: func(bal Balance) *Data {
		return toData(bal)
	},
	RouteList: func(params RouteListParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create the repositories:
			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			depositRepository := deposit.SDKFunc.CreateRepository(deposit.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			withdrawalRepository := withdrawal.SDKFunc.CreateRepository(withdrawal.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			balanceRepository := createRepository(depositRepository, withdrawalRepository)

			// retrieve the wallets:
			walsPS, walsPSErr := walletRepository.RetrieveSet(0, params.AmountOfElementsPerList)
			if walsPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving a wallet set: %s", walsPSErr.Error())
				w.Write([]byte(str))
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

			// convert to data:
			dataSet, dataSetErr := convertToDataSet(gen.Deposit().Token(), walsPS, balanceRepository)
			if dataSetErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while creating a Balance partial set data: %s", dataSetErr.Error())
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
			// create the repositories:
			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			depositRepository := deposit.SDKFunc.CreateRepository(deposit.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			withdrawalRepository := withdrawal.SDKFunc.CreateRepository(withdrawal.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			balanceRepository := createRepository(depositRepository, withdrawalRepository)

			vars := mux.Vars(r)
			if walletIDAsString, ok := vars["id"]; ok {
				walletID, walletIDErr := uuid.FromString(walletIDAsString)
				if walletIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given walletID (ID: %s) is invalid: %s", walletIDAsString, walletIDErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the walletID:
				retWal, retWalErr := walletRepository.RetrieveByID(&walletID)
				if retWalErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s): %s", walletID.String(), retWalErr.Error())
					w.Write([]byte(str))
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

				// retrieve the balance:
				bal, balErr := balanceRepository.RetrieveByWalletAndToken(retWal, gen.Deposit().Token())
				if balErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the balance (WalletID: %s, TokenID: %s): %s", retWal.ID().String(), gen.Deposit().Token().ID().String(), balErr.Error())
					w.Write([]byte(str))
					return
				}

				// render:
				w.WriteHeader(http.StatusOK)
				params.Tmpl.Execute(w, toData(bal))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the ID could not be found")
			w.Write([]byte(str))
		}
	},
}
