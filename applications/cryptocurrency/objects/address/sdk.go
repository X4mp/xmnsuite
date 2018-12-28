package address

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Address represents an address bind between a wallet and a cryptocurrency address
type Address interface {
	ID() *uuid.UUID
	Wallet() wallet.Wallet
	Address() string
}

// Normalized represents a normalized offer
type Normalized interface {
}

// Repository represents the address repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Address, error)
	RetrieveByAddress(addr string) (Address, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByWallet(wal wallet.Wallet, index int, amount int) (entity.PartialSet, error)
}

// Data represents human-readable data
type Data struct {
	ID      string
	Wallet  *wallet.Data
	Address string
}

// DataSet represents the human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Addresses   []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID      *uuid.UUID
	Wallet  wallet.Wallet
	Address string
}

// RouteParams represents the route params
type RouteParams struct {
	Tmpl             *template.Template
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

// SDKFunc represents the Address SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Address
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	ToData               func(addr Address) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	Route                func(params RouteParams) func(w http.ResponseWriter, r *http.Request)
	RouteSet             func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request)
	RouteNew             func(params RouteNewParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Address {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createAddress(params.ID, params.Wallet, params.Address)
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
	ToData: func(addr Address) *Data {
		return toData(addr)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	Route: func(params RouteParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create the metadata:
			metaData := createMetaData()

			// create the repositories:
			addressRepository := createRepository(metaData, params.EntityRepository)

			// get the id from the uri:
			vars := mux.Vars(r)
			if idAsString, ok := vars["id"]; ok {
				// convert the string to an id:
				id, idErr := uuid.FromString(idAsString)
				if idErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the Address ID (%s) is invalid", idAsString)
					w.Write([]byte(str))
					return
				}

				// retrieve the address by id:
				addr, addrErr := addressRepository.RetrieveByID(&id)
				if idErr != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(addrErr.Error()))
					return
				}

				w.WriteHeader(http.StatusOK)
				params.Tmpl.Execute(w, toData(addr))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the ID could not be found")
			w.Write([]byte(str))
		}
	},
	RouteSet: func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create the metadata:
			metaData := createMetaData()

			// create the repositories:
			addressRepository := createRepository(metaData, params.EntityRepository)

			// retrieve all the addresses:
			addrPS, addrPSErr := addressRepository.RetrieveSet(0, params.AmountOfElementsPerList)
			if addrPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving the address entity set: %s", addrPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// render:
			datSet, datSetErr := toDataSet(addrPS)
			if datSetErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while converting the address entity set to data: %s", datSetErr.Error())
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
			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
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

			// if the form has been submitted:
			fromWalletIDAsString := r.FormValue("fromwalletid")
			addressAsString := r.FormValue("address")
			reason := r.FormValue("reason")
			if fromWalletIDAsString != "" && addressAsString != "" {
				fromWalletID, fromWalletIDErr := uuid.FromString(fromWalletIDAsString)
				if fromWalletIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there given fromWalletID (%s) is invalid: %s", fromWalletIDAsString, fromWalletIDErr.Error())
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

				// create the address:
				id := uuid.NewV4()
				addr, addrErr := createAddress(&id, wal, addressAsString)
				if addrErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while creating an Address instance: %s", addrErr.Error())
					w.Write([]byte(str))
					return
				}

				// create the request:
				req := request.SDKFunc.Create(request.CreateParams{
					FromUser:  usr,
					NewEntity: addr,
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

			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, user.SDKFunc.ToDataSet(usrPS))
		}
	},
}
