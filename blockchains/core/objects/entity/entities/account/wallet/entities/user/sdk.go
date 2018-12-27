package user

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// User represents a user
type User interface {
	ID() *uuid.UUID
	PubKey() crypto.PublicKey
	Shares() int
	Wallet() wallet.Wallet
}

// Normalized represents a normalized user
type Normalized interface {
}

// Repository represents the user repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (User, error)
	RetrieveByPubKeyAndWallet(pubKey crypto.PublicKey, wal wallet.Wallet) (User, error)
	RetrieveSetByPubKey(pubKey crypto.PublicKey, index int, amount int) (entity.PartialSet, error)
	RetrieveSetByWallet(wal wallet.Wallet, index int, amount int) (entity.PartialSet, error)
}

// Data represents human-redable data
type Data struct {
	ID     string
	PubKey string
	Shares int
	Wallet *wallet.Data
}

// WalletWithData represents human0readable wallet data with its users
type WalletWithData struct {
	Wallet *wallet.Data
	Users  *DataSet
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Users       []*Data
}

// CreateParams represents the Create params
type CreateParams struct {
	ID     *uuid.UUID
	PubKey crypto.PublicKey
	Shares int
	Wallet wallet.Wallet
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	Store            datastore.DataStore
	EntityRepository entity.Repository
}

// RouteWalletListParams represents the route wallet list params
type RouteWalletListParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// RouteUserSetInWalletParams represents the route user set in wallet params
type RouteUserSetInWalletParams struct {
	AmountOfElementsPerList int
	Tmpl                    *template.Template
	EntityRepository        entity.Repository
}

// RouteUserInWalletParams represents the route user in wallet params
type RouteUserInWalletParams struct {
	Tmpl             *template.Template
	EntityRepository entity.Repository
}

// SDKFunc represents the User SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) User
	CreateRepository     func(params CreateRepositoryParams) Repository
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	ToData               func(usr User) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	RouteWalletList      func(params RouteWalletListParams) func(w http.ResponseWriter, r *http.Request)
	RouteUserSetInWallet func(params RouteUserSetInWalletParams) func(w http.ResponseWriter, r *http.Request)
	RouteUserInWallet    func(params RouteUserInWalletParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) User {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createUser(params.ID, params.PubKey, params.Shares, params.Wallet)
		return out
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		if params.Store != nil {
			params.EntityRepository = entity.SDKFunc.CreateRepository(params.Store)
		}

		userMetaData := createMetaData()
		out := createRepository(userMetaData, params.EntityRepository)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if usr, ok := ins.(User); ok {
					out := createStorableUser(usr)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if usr, ok := ins.(User); ok {
					return []string{
						retrieveAllUserKeyname(),
						retrieveUserByPubKeyKeyname(usr.PubKey()),
						retrieveUserByWalletIDKeyname(usr.Wallet().ID()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
				return nil, errors.New(str)

			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {

				if usr, ok := ins.(User); ok {

					//create the metadata and representation:
					metaData := createMetaData()
					walRepresentation := wallet.SDKFunc.CreateRepresentation()

					// create the repositories and services:
					entityRepository := entity.SDKFunc.CreateRepository(ds)
					repository := createRepository(metaData, entityRepository)
					entityService := entity.SDKFunc.CreateService(ds)

					// make sure there is no other user with the given public key, on the same wallet:
					_, retUserErr := repository.RetrieveByPubKeyAndWallet(usr.PubKey(), usr.Wallet())
					if retUserErr == nil {
						str := fmt.Sprintf("the User instance (PubKey: %s, WalletID: %s) already exists", usr.PubKey().String(), usr.Wallet().ID().String())
						return errors.New(str)
					}

					// if the wallet does not exists, create it:
					wal := usr.Wallet()
					_, retWalErr := entityRepository.RetrieveByID(walRepresentation.MetaData(), wal.ID())
					if retWalErr != nil {
						saveWalErr := entityService.Save(wal, walRepresentation)
						if saveWalErr != nil {
							return saveWalErr
						}
					}

					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid User instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	ToData: func(usr User) *Data {
		return toData(usr)
	},
	ToDataSet: func(ps entity.PartialSet) *DataSet {
		out, outErr := toDataSet(ps)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	RouteWalletList: func(params RouteWalletListParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create the repositories:
			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			// retrieve the wallet set:
			walPS, walPSErr := walletRepository.RetrieveSet(0, params.AmountOfElementsPerList)
			if walPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieve wallet entity set: %s", walPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// render:
			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, wallet.SDKFunc.ToDataSet(walPS))
			return
		}
	},
	RouteUserSetInWallet: func(params RouteUserSetInWalletParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create metadata:
			userMetaData := createMetaData()

			// create the repositories:
			userRepository := createRepository(userMetaData, params.EntityRepository)
			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			vars := mux.Vars(r)
			if walletIDAsString, ok := vars["wallet_id"]; ok {
				// parse the walletID:
				walletID, walletIDErr := uuid.FromString(walletIDAsString)
				if walletIDErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the given walletID (ID: %s) is invalid: %s", walletIDAsString, walletIDErr.Error())
					w.Write([]byte(str))
					return
				}

				// retireve the wallet:
				wal, walErr := walletRepository.RetrieveByID(&walletID)
				if walErr != nil {
					w.WriteHeader(http.StatusNotFound)
					str := fmt.Sprintf("the given WalletID (ID: %s) does not exists: %s", walletID.String(), walErr.Error())
					w.Write([]byte(str))
					return
				}

				// retrieve the users:
				usrsPS, usrsPSErr := userRepository.RetrieveSetByWallet(wal, 0, params.AmountOfElementsPerList)
				if usrsPSErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving a user entity set: %s", usrsPSErr.Error())
					w.Write([]byte(str))
					return
				}

				// render:
				datSet, datSetErr := toDataSet(usrsPS)
				if datSetErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while converting the user entity set to data: %s", datSetErr.Error())
					w.Write([]byte(str))
					return
				}

				w.WriteHeader(http.StatusOK)
				params.Tmpl.Execute(w, &WalletWithData{
					Wallet: wallet.SDKFunc.ToData(wal),
					Users:  datSet,
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the walletID could not be found")
			w.Write([]byte(str))
		}
	},
	RouteUserInWallet: func(params RouteUserInWalletParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create metadata:
			userMetaData := createMetaData()

			// create the repositories:
			userRepository := createRepository(userMetaData, params.EntityRepository)
			walletRepository := wallet.SDKFunc.CreateRepository(wallet.CreateRepositoryParams{
				EntityRepository: params.EntityRepository,
			})

			vars := mux.Vars(r)
			if walletIDAsString, ok := vars["wallet_id"]; ok {
				if pubKeyAsString, ok := vars["pubkey"]; ok {
					// parse the walletID:
					walletID, walletIDErr := uuid.FromString(walletIDAsString)
					if walletIDErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("the given walletID (ID: %s) is invalid: %s", walletIDAsString, walletIDErr.Error())
						w.Write([]byte(str))
						return
					}

					// parse the pubkey:
					pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
						PubKeyAsString: pubKeyAsString,
					})

					// retrieve the wallet:
					wal, walErr := walletRepository.RetrieveByID(&walletID)
					if walErr != nil {
						w.WriteHeader(http.StatusNotFound)
						str := fmt.Sprintf("there was an error while retrieving the Wallet (ID: %s): %s", walletID.String(), walErr.Error())
						w.Write([]byte(str))
						return
					}

					// retrieve the user:
					usr, usrErr := userRepository.RetrieveByPubKeyAndWallet(pubKey, wal)
					if usrErr != nil {
						w.WriteHeader(http.StatusNotFound)
						str := fmt.Sprintf("there was an error while retrieving the User (PubKey: %s, WalletID: %s): %s", pubKey.String(), wal.ID().String(), usrErr.Error())
						w.Write([]byte(str))
						return
					}

					// render:
					w.WriteHeader(http.StatusOK)
					params.Tmpl.Execute(w, toData(usr))
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("the pubkey could not be found")
				w.Write([]byte(str))
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the walletID could not be found")
			w.Write([]byte(str))
		}
	},
}
