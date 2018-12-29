package offer

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
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

// Offer represents an offer to manage coins
type Offer interface {
	ID() *uuid.UUID
	Pledge() pledge.Pledge
	To() address.Address
	Confirmations() int
	Amount() int
	Price() int
	IP() net.IP
	Port() int
}

// Normalized represents a normalized offer
type Normalized interface {
}

// Repository represents the offer repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Offer, error)
	RetrieveByPledge(pldge pledge.Pledge) (Offer, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByToAddress(toAddr address.Address, index int, amount int) (entity.PartialSet, error)
}

// Data represents human-readable data
type Data struct {
	ID            string
	Pledge        *pledge.Data
	ToAddress     *address.Data
	Confirmations int
	Amount        int
	Price         int
}

// DataSet represents the human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Offers      []*Data
}

// DataNew represents the human-readable data new
type DataNew struct {
	Balances  []*balance.Data
	Addresses []*address.Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID            *uuid.UUID
	Pledge        pledge.Pledge
	To            address.Address
	Confirmations int
	Amount        int
	Price         int
	IP            net.IP
	Port          int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// RouteSetParams represents the route set params
type RouteSetParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// RouteNewParams represents the route new params
type RouteNewParams struct {
	AmountOfElementsPerList int
	Client                  applications.Client
	PK                      crypto.PrivateKey
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// SDKFunc represents the Offer SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Offer
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(off Offer) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteSet             func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request)
	RouteNew             func(params RouteNewParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Offer {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createOffer(params.ID, params.Pledge, params.To, params.Confirmations, params.Amount, params.Price, params.IP, params.Port)
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
	ToData: func(off Offer) *Data {
		return toData(off)
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
			offerRepository := createRepository(metaData, params.EntityRepository)

			// retrieve all the offers:
			offPS, offPSErr := offerRepository.RetrieveSet(0, params.AmountOfElementsPerList)
			if offPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving the offer entity set: %s", offPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// render:
			datSet, datSetErr := toDataSet(offPS)
			if datSetErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while converting the offer entity set to data: %s", datSetErr.Error())
				w.Write([]byte(str))
				return
			}

			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, datSet)
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

			// create metadata:
			representation := createRepresentation()

			// create the repositories:
			genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
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

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			addressRepository := address.SDKFunc.CreateRepository(address.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
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
				str := fmt.Sprintf("there was an error while retrieving the Genesis instance: %s", genErr.Error())
				w.Write([]byte(str))
				return
			}

			// if the form has been submitted:
			fromWalletIDAsString := r.FormValue("fromwalletid")
			addressIDAsString := r.FormValue("addressid")
			pledgeAmountAsString := r.FormValue("pledgeamount")
			satoshiAmountAsString := r.FormValue("satoshiamount")
			servicePriceAsString := r.FormValue("serviceprice")
			ipAddressAsString := r.FormValue("ipaddress")
			portAsString := r.FormValue("port")
			confirmationsAsString := r.FormValue("confirmations")
			reason := r.FormValue("reason")
			if fromWalletIDAsString != "" && pledgeAmountAsString != "" && addressIDAsString != "" && satoshiAmountAsString != "" && servicePriceAsString != "" && ipAddressAsString != "" && portAsString != "" && confirmationsAsString != "" {
				fromWalletID, fromWalletIDErr := uuid.FromString(fromWalletIDAsString)
				if fromWalletIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there given fromWalletID (%s) is invalid: %s", fromWalletIDAsString, fromWalletIDErr.Error())
					w.Write([]byte(str))
					return
				}

				addressID, addressIDErr := uuid.FromString(addressIDAsString)
				if addressIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there given addressID (%s) is invalid: %s", addressIDAsString, addressIDErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the wallet:
				wal, walErr := walletRepository.RetrieveByID(&fromWalletID)
				if walErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the Wallet (ID: %s): %s", fromWalletID, walErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the address:
				addr, addrErr := addressRepository.RetrieveByID(&addressID)
				if addrErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the Address (ID: %s): %s", addressID.String(), addrErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the user:
				usr, usrErr := userRepository.RetrieveByPubKeyAndWallet(params.PK.PublicKey(), wal)
				if usrErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the User (Pubkey: %s, WalletID: %s): %s", params.PK.PublicKey().String(), wal.ID().String(), usrErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the keyname:
				kname, knameErr := keynameRepository.RetrieveByName(representation.MetaData().Keyname())
				if knameErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the Keyname (Name; %s): %s", representation.MetaData().Keyname(), knameErr.Error())
					w.Write([]byte(str))
					return
				}

				pldge, pldgeErr := strconv.Atoi(pledgeAmountAsString)
				if pldgeErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given pledge amount (%s) is not a valid number: %s", pledgeAmountAsString, pldgeErr.Error())
					w.Write([]byte(str))
					return
				}

				amount, amountErr := strconv.Atoi(satoshiAmountAsString)
				if amountErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given satoshi amount (%s) is not a valid number: %s", satoshiAmountAsString, amountErr.Error())
					w.Write([]byte(str))
					return
				}

				price, priceErr := strconv.Atoi(servicePriceAsString)
				if priceErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given service price (%s) is not a valid number: %s", servicePriceAsString, priceErr.Error())
					w.Write([]byte(str))
					return
				}

				port, portErr := strconv.Atoi(portAsString)
				if portErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given port (%s) is not a valid number: %s", portAsString, portErr.Error())
					w.Write([]byte(str))
					return
				}

				conf, confErr := strconv.Atoi(confirmationsAsString)
				if portErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given confirmations amount (%s) is not a valid number: %s", confirmationsAsString, confErr.Error())
					w.Write([]byte(str))
					return
				}

				// create the offer:
				id := uuid.NewV4()
				off, offErr := createOffer(&id, pledge.SDKFunc.Create(pledge.CreateParams{
					From: withdrawal.SDKFunc.Create(withdrawal.CreateParams{
						From:   usr.Wallet(),
						Token:  gen.Deposit().Token(),
						Amount: pldge,
					}),
					To: gen.Deposit().To(),
				}), addr, conf, amount, price, net.ParseIP(ipAddressAsString), port)

				if offErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while creating an offer instance: %s", offErr.Error())
					w.Write([]byte(str))
					return
				}

				// create the request:
				req := request.SDKFunc.Create(request.CreateParams{
					FromUser:  usr,
					NewEntity: off,
					Reason:    reason,
					Keyname:   kname,
				})

				// save the request:
				saveErr := requestService.Save(req, representation)
				if saveErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while saving a Request instance: %s", saveErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the active request:
				activeReq, activeReqErr := requestRepository.RetrieveByRequest(req)
				if activeReqErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while cretrieving an ActiveRequest: %s", activeReqErr.Error())
					w.Write([]byte(str))
					return
				}

				// redirect:
				url := fmt.Sprintf("/requests/%s/%s/%s", activeReq.Request().Keyname().Group().Name(), activeReq.Request().Keyname().Name(), activeReq.ID().String())
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				return
			}

			// retrieve all the users related to my pubkey:
			usrPS, usrPSErr := userRepository.RetrieveSetByPubKey(params.PK.PublicKey(), 0, params.AmountOfElementsPerList)
			if usrPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving the user entity set: %s", usrPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// for each user, retrieve the balances and addresses:
			usrs := usrPS.Instances()
			tok := gen.Deposit().Token()
			balances := []*balance.Data{}
			addresses := []*address.Data{}
			for _, oneUserIns := range usrs {
				if usr, ok := oneUserIns.(user.User); ok {
					// retrieve the balance:
					bal, balErr := balanceRepository.RetrieveByWalletAndToken(usr.Wallet(), tok)
					if balErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while retrieving the balance (WalletID: %s, TokenID: %s): %s", usr.Wallet().ID().String(), tok.ID().String(), balErr.Error())
						w.Write([]byte(str))
						return
					}

					// store the balance:
					balances = append(balances, balance.SDKFunc.ToData(bal))

					// retrieve the address:
					addrPS, addrPSErr := addressRepository.RetrieveSetByWallet(usr.Wallet(), 0, -1)
					if addrPSErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while retrieving the address entity set (WalletID: %s): %s", usr.Wallet().ID().String(), addrPSErr.Error())
						w.Write([]byte(str))
						return
					}

					addrIns := addrPS.Instances()
					for _, oneAddrIns := range addrIns {
						if oneAddr, ok := oneAddrIns.(address.Address); ok {
							addresses = append(addresses, address.SDKFunc.ToData(oneAddr))
						}
					}
				}

			}

			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, &DataNew{
				Balances:  balances,
				Addresses: addresses,
			})
		}
	},
}
