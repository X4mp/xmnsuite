package deposit

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Deposit represents a deposit on an offer
type Deposit interface {
	ID() *uuid.UUID
	Offer() offer.Offer
	From() address.Address
	Amount() int
}

// Normalized represents a normalized offer
type Normalized interface {
}

// Repository represents the deposit repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Deposit, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByOffer(off offer.Offer, index int, amount int) (entity.PartialSet, error)
	RetrieveSetByFromAddress(frmAddress address.Address, index int, amount int) (entity.PartialSet, error)
}

// Data represents human-readable data
type Data struct {
	ID     string
	Offer  *offer.Data
	From   *address.Data
	Amount int
}

// DataSet represents the human-readable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Deposits    []*Data
}

// DataSetWithOffer represents the human-readable data set with offer
type DataSetWithOffer struct {
	Offer       *offer.Data
	Deposits    *DataSet
	MyAddresses []*address.Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	Offer  offer.Offer
	From   address.Address
	Amount int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
}

// RouteSetParams represents the route set params
type RouteSetParams struct {
	PK               crypto.PrivateKey
	Client           applications.Client
	Tmpl             *template.Template
	EntityRepository entity.Repository
}

// SDKFunc represents the Deposit SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Deposit
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(dep Deposit) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteSet             func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Deposit {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out, outErr := createDeposit(params.ID, params.Offer, params.From, params.Amount)
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
	ToData: func(dep Deposit) *Data {
		return toData(dep)
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
			representation := createRepresentation()

			// create the repositories:
			depositRepository := createRepository(representation.MetaData(), params.EntityRepository)
			offerRepository := offer.SDKFunc.CreateRepository(offer.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			addressRepository := address.SDKFunc.CreateRepository(address.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			userRepository := user.SDKFunc.CreateRepository(user.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			keynameRepository := keyname.SDKFunc.CreateRepository(keyname.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			requestService := request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
				PK:          params.PK,
				Client:      params.Client,
				RoutePrefix: "",
			})

			requestRepository := active_request.SDKFunc.CreateRepository(active_request.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			// get the id from the uri:
			vars := mux.Vars(r)
			if idAsString, ok := vars["id"]; ok {
				// convert the string to an id:
				id, idErr := uuid.FromString(idAsString)
				if idErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the wallet ID (%s) is invalid", idAsString)
					w.Write([]byte(str))
					return
				}

				// retrieve the offer by id:
				off, offErr := offerRepository.RetrieveByID(&id)
				if offErr != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(offErr.Error()))
					return
				}

				// if the form has been submitted:
				addressIDAsString := r.FormValue("addressid")
				amountAsString := r.FormValue("amount")
				reason := r.FormValue("reason")
				if addressIDAsString != "" && amountAsString != "" {
					addressID, addressIDErr := uuid.FromString(addressIDAsString)
					if addressIDErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there given addressID (%s) is invalid: %s", addressIDAsString, addressIDErr.Error())
						w.Write([]byte(str))
						return
					}

					amt, amtErr := strconv.Atoi(amountAsString)
					if amtErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("the given amount (%s) is not a valid number: %s", amountAsString, amtErr.Error())
						w.Write([]byte(str))
						return
					}

					// retrieve the from address:
					fromAddr, fromAddrErr := addressRepository.RetrieveByID(&addressID)
					if fromAddrErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while retrieving an Address (ID: %s): %s", &addressID, fromAddrErr.Error())
						w.Write([]byte(str))
						return
					}

					// retrieve the from user:
					fromUser, fromUSerErr := userRepository.RetrieveByPubKeyAndWallet(params.PK.PublicKey(), fromAddr.Wallet())
					if fromUSerErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while retrieving the from User (PubKey: %s, WalletID: %s): %s", params.PK.PublicKey().String(), fromAddr.Wallet().ID().String(), fromUSerErr.Error())
						w.Write([]byte(str))
						return
					}

					// retrieve the keyname:
					kname, knameErr := keynameRepository.RetrieveByName(representation.MetaData().Keyname())
					if knameErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while retrieving the Keyname (Name: %s): %s", representation.MetaData().Keyname(), knameErr.Error())
						w.Write([]byte(str))
						return
					}

					// create the new deposit:
					id := uuid.NewV4()
					dep, depErr := createDeposit(&id, off, fromAddr, amt)
					if depErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while creating a Deposit instance: %s", depErr.Error())
						w.Write([]byte(str))
						return
					}

					// create the request:
					req := request.SDKFunc.Create(request.CreateParams{
						FromUser:  fromUser,
						NewEntity: dep,
						Reason:    reason,
						Keyname:   kname,
					})

					// save the request:
					saveReqErr := requestService.Save(req, representation)
					if saveReqErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while saving the request: %s", saveReqErr.Error())
						w.Write([]byte(str))
						return
					}

					// retrieve the active request:
					activeReq, activeReqErr := requestRepository.RetrieveByRequest(req)
					if activeReqErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while retrieving an ActiveRequest by Request (ID: %s): %s", req.ID().String(), activeReqErr.Error())
						w.Write([]byte(str))
						return
					}

					// redirect:
					url := fmt.Sprintf("/requests/%s/%s/%s", activeReq.Request().Keyname().Group().Name(), activeReq.Request().Keyname().Name(), activeReq.ID().String())
					http.Redirect(w, r, url, http.StatusTemporaryRedirect)
					return
				}

				// retrieve the deposits related to the offer:
				depPS, depPSErr := depositRepository.RetrieveSetByOffer(off, 0, -1)
				if depPSErr != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(depPSErr.Error()))
					return
				}

				// retrieve the users related to our pubKey:
				myUsersPS, myUsersPSErr := userRepository.RetrieveSetByPubKey(params.PK.PublicKey(), 0, -1)
				if myUsersPSErr != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(myUsersPSErr.Error()))
					return
				}

				// for each user, retrieve the addresses:
				usersIns := myUsersPS.Instances()
				addresses := []*address.Data{}
				for _, oneUserIns := range usersIns {
					if usr, ok := oneUserIns.(user.User); ok {
						addrPS, addrPSErr := addressRepository.RetrieveSetByWallet(usr.Wallet(), 0, -1)
						if addrPSErr != nil {
							w.WriteHeader(http.StatusNotFound)
							w.Write([]byte(myUsersPSErr.Error()))
							return
						}

						addrsIns := addrPS.Instances()
						for _, oneAddrIns := range addrsIns {
							if addr, ok := oneAddrIns.(address.Address); ok {
								addresses = append(addresses, address.SDKFunc.ToData(addr))
							}
						}
					}
				}

				// render:
				datPS, datPSErr := toDataSet(depPS)
				if datPSErr != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(datPSErr.Error()))
					return
				}

				w.WriteHeader(http.StatusOK)
				params.Tmpl.Execute(w, &DataSetWithOffer{
					Offer:       offer.SDKFunc.ToData(off),
					Deposits:    datPS,
					MyAddresses: addresses,
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the ID could not be found")
			w.Write([]byte(str))
		}
	},
}
