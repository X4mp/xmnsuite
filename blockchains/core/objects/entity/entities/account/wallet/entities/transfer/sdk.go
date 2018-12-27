package transfer

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
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

// Transfer represents a transfer of token that can be claimed
type Transfer interface {
	ID() *uuid.UUID
	Withdrawal() withdrawal.Withdrawal
	Deposit() deposit.Deposit
}

// Normalized represents the normalized transfer
type Normalized interface {
}

// Repository represents the transfer reposiotry
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Transfer, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByDeposit(dep deposit.Deposit, index int, amount int) (entity.PartialSet, error)
	RetrieveSetByWithdrawal(with withdrawal.Withdrawal, index int, amount int) (entity.PartialSet, error)
}

// Data represents human-redable data
type Data struct {
	ID         string
	Withdrawal *withdrawal.Data
	Deposit    *deposit.Data
}

// DataNew represents human-redable data for the new transfer page
type DataNew struct {
	From []*balance.Data
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Transfers   []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID         *uuid.UUID
	Withdrawal withdrawal.Withdrawal
	Deposit    deposit.Deposit
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// RouteSetParams represents the route transfer list params
type RouteSetParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// RouteParams represents the route params
type RouteParams struct {
	Tmpl             *template.Template
	EntityRepository entity.Repository
}

// RouteNewParams represents the route to create a new transfer params
type RouteNewParams struct {
	PK               crypto.PrivateKey
	Client           applications.Client
	Tmpl             *template.Template
	EntityRepository entity.Repository
}

// SDKFunc represents the Transfer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Transfer
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(trsf Transfer) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteSet             func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request)
	Route                func(params RouteParams) func(w http.ResponseWriter, r *http.Request)
	RouteNew             func(params RouteNewParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Transfer {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createTransfer(params.ID, params.Withdrawal, params.Deposit)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		out := createMetaData()
		return out
	},
	CreateRepresentation: func() entity.Representation {
		return createRepresentation()
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		out := createRepository(params.EntityRepository, metaData)
		return out
	},
	ToData: func(trsf Transfer) *Data {
		return toData(trsf)
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
			// create metadata:
			metaData := createMetaData()

			// create the repositories:
			transferRepository := createRepository(params.EntityRepository, metaData)

			// retrieve the transfers:
			trsfPS, trsfPSErr := transferRepository.RetrieveSet(0, params.AmountOfElementsPerList)
			if trsfPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving a transfer set: %s", trsfPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// convert to data:
			dataSet, dataSetErr := toDataSet(trsfPS)
			if dataSetErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while convert a Transfer partial set to data: %s", dataSetErr.Error())
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
			if transferIDAsString, ok := vars["id"]; ok {
				// parse the id:
				transferID, transferIDErr := uuid.FromString(transferIDAsString)
				if transferIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given ID (ID: %s) is invalid: %s", transferIDAsString, transferIDErr.Error())
					w.Write([]byte(str))
					return
				}

				// create the metadata:
				transferMetaData := createMetaData()

				// create the repositories:
				transferRepository := createRepository(params.EntityRepository, transferMetaData)

				// retrieve the transfer:
				retTrsf, retTrsfErr := transferRepository.RetrieveByID(&transferID)
				if retTrsfErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving a Transfer (ID: %s): %s", transferID.String(), retTrsfErr.Error())
					w.Write([]byte(str))
					return
				}

				// render:
				w.WriteHeader(http.StatusOK)
				params.Tmpl.Execute(w, toData(retTrsf))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the ID could not be found")
			w.Write([]byte(str))
		}
	},
	RouteNew: func(params RouteNewParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// parse the form:
			if parseFormErr := r.ParseForm(); parseFormErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while parsing form elements: %s", parseFormErr.Error())
				w.Write([]byte(str))
				return
			}

			// create the metadata and representation:
			representation := createRepresentation()

			// create the repositories:
			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
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

			requestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
				PK:          params.PK,
				Client:      params.Client,
				RoutePrefix: "",
			})

			// retrieve the genesis:
			gen, genErr := genesisRepository.Retrieve()
			if genErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
				w.Write([]byte(str))
				return
			}

			// if the form has been submitted:
			fromWalletIDAsString := r.FormValue("fromwalletid")
			toWalletIDAsString := r.FormValue("towalletid")
			amountAsString := r.FormValue("amount")
			reason := r.FormValue("reason")
			if fromWalletIDAsString != "" && toWalletIDAsString != "" {
				fromWalletID, fromWalletIDErr := uuid.FromString(fromWalletIDAsString)
				if fromWalletIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the fromWalletID (ID: %s) is invalid: %s", fromWalletIDAsString, fromWalletIDErr.Error())
					w.Write([]byte(str))
					return
				}

				toWalletID, toWalletIDErr := uuid.FromString(toWalletIDAsString)
				if toWalletIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the toWalletID (ID: %s) is invalid: %s", toWalletIDAsString, toWalletIDErr.Error())
					w.Write([]byte(str))
					return
				}

				amount, amountErr := strconv.Atoi(amountAsString)
				if amountErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the amount (%s) is invalid: %s", amountAsString, amountErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the from wallet:
				fromWallet, fromWalletErr := walletRepository.RetrieveByID(&fromWalletID)
				if fromWalletErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the from wallet (ID: %s): %s", fromWalletID.String(), fromWalletErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the to wallet:
				toWallet, toWalletErr := walletRepository.RetrieveByID(&toWalletID)
				if toWalletErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the to wallet (ID: %s): %s", toWalletID.String(), toWalletErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the from user:
				fromUser, fromUserErr := userRepository.RetrieveByPubKeyAndWallet(params.PK.PublicKey(), fromWallet)
				if fromUserErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the fromUser (PubKey: %s, FromWalletID: %s): %s", params.PK.PublicKey().String(), fromWallet.ID().String(), fromUserErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the keyname:
				keyname, keynameErr := keynameRepository.RetrieveByName(representation.MetaData().Keyname())
				if keynameErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the keyname (name: %s): %s", representation.MetaData().Keyname(), keynameErr.Error())
					w.Write([]byte(str))
					return
				}

				// create the request:
				id := uuid.NewV4()
				tok := gen.Deposit().Token()
				req := request.SDKFunc.Create(request.CreateParams{
					FromUser: fromUser,
					NewEntity: createTransfer(
						&id,
						withdrawal.SDKFunc.Create(withdrawal.CreateParams{
							From:   fromWallet,
							Token:  tok,
							Amount: amount,
						}),
						deposit.SDKFunc.Create(deposit.CreateParams{
							To:     toWallet,
							Token:  tok,
							Amount: amount,
						}),
					),
					Reason:  reason,
					Keyname: keyname,
				})

				// save the request:
				saveReqErr := requestService.Save(req, representation)
				if saveReqErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while saving a transfer request: %s", saveReqErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the active request:
				retActiveRequest, retActiveRequestErr := requestRepository.RetrieveByRequest(req)
				if retActiveRequestErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving an ActiveRequest from a Request (ID: %s): %s", req.ID().String(), retActiveRequestErr.Error())
					w.Write([]byte(str))
					return
				}

				// redirect:
				url := fmt.Sprintf("/requests/%s/%s/%s", req.Keyname().Group().Name(), req.Keyname().Name(), retActiveRequest.ID().String())
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				return
			}

			// retrieve the users:
			usrsPS, usrsPSErr := userRepository.RetrieveSetByPubKey(params.PK.PublicKey(), 0, -1)
			if usrsPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving a user set from pubkey (PubKey: %s): %s", params.PK.PublicKey().String(), usrsPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// for each user, retrieve the balance:
			usrs := usrsPS.Instances()
			balances := []*balance.Data{}
			for _, oneIns := range usrs {
				if usr, ok := oneIns.(user.User); ok {
					bal, balErr := balanceRepository.RetrieveByWalletAndToken(usr.Wallet(), gen.Deposit().Token())
					if balErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while retrieving the balance (WalletID: %s, TokenID: %s): %s", usr.Wallet().ID().String(), gen.Deposit().Token().ID(), balErr.Error())
						w.Write([]byte(str))
						return
					}

					// convert the balance to data:
					balances = append(balances, balance.SDKFunc.ToData(bal))
				}

			}

			// render:
			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, &DataNew{
				From: balances,
			})
			return
		}
	},
}
