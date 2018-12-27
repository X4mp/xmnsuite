package wallet

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
)

// Wallet represents a wallet
type Wallet interface {
	ID() *uuid.UUID
	Creator() crypto.PublicKey
	ConcensusNeeded() int
}

// Data represents human-redable data
type Data struct {
	ID              string
	Creator         string
	ConcensusNeeded int
}

// DataSet represents human-redable data set
type DataSet struct {
	Index       int
	Amount      int
	TotalAmount int
	IsLast      bool
	Wallets     []*Data
}

// Repository represents the wallet repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Wallet, error)
	RetrieveSet(index int, amount int) (entity.PartialSet, error)
	RetrieveSetByCreatorPublicKey(pubKey crypto.PublicKey, index int, amount int) (entity.PartialSet, error)
}

// Normalized represents a normalized wallet
type Normalized interface {
}

// CreateParams represents the Create params
type CreateParams struct {
	ID              *uuid.UUID
	Creator         crypto.PublicKey
	ConcensusNeeded int
}

// CreateRepositoryParams represents the CreateRepository params
type CreateRepositoryParams struct {
	EntityRepository entity.Repository
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

// SDKFunc represents the Wallet SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Wallet
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	CreateRepository     func(params CreateRepositoryParams) Repository
	ToData               func(wal Wallet) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	Route                func(params RouteParams) func(w http.ResponseWriter, r *http.Request)
	RouteSet             func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request)
}{
	Create: func(params CreateParams) Wallet {
		if params.ID == nil {
			id := uuid.NewV4()
			params.ID = &id
		}

		out := createWallet(params.ID, params.Creator, params.ConcensusNeeded)
		return out
	},
	CreateMetaData: func() entity.MetaData {
		return createMetaData()
	},
	CreateRepresentation: func() entity.Representation {
		return entity.SDKFunc.CreateRepresentation(entity.CreateRepresentationParams{
			Met: createMetaData(),
			ToStorable: func(ins entity.Entity) (interface{}, error) {
				if wallet, ok := ins.(Wallet); ok {
					out := createStoredWallet(wallet)
					return out, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", ins.ID().String())
				return nil, errors.New(str)
			},
			Keynames: func(ins entity.Entity) ([]string, error) {
				if wal, ok := ins.(Wallet); ok {
					return []string{
						retrieveAllWalletKeyname(),
						retrieveByPublicKeyWalletKeyname(wal.Creator()),
					}, nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", ins.ID().String())
				return nil, errors.New(str)

			},
			Sync: func(ds datastore.DataStore, ins entity.Entity) error {
				// create the repository and service:
				repository := entity.SDKFunc.CreateRepository(ds)

				// create the metadata:
				metaData := createMetaData()

				if wal, ok := ins.(Wallet); ok {
					// if the wallet already exists:
					retWalletIns, retWalletInsErr := repository.RetrieveByID(metaData, wal.ID())
					if retWalletInsErr == nil {
						// cast the returned wallet:
						if retWal, ok := retWalletIns.(Wallet); ok {
							// make sure the creator is the same:
							if !retWal.Creator().Equals(wal.Creator()) {
								str := fmt.Sprintf("the Wallet (ID: %s) already existed but the creator pubKey does not match.  Expected: %s, Received: %s", retWal.ID().String(), wal.Creator().String(), retWal.Creator().String())
								return errors.New(str)
							}

							// everything is fine, it will update the wallet:
							return nil
						}
					}

					// the wallet doesnt exists, so everything is fine:
					return nil
				}

				str := fmt.Sprintf("the given entity (ID: %s) is not a valid Wallet instance", ins.ID().String())
				return errors.New(str)
			},
		})
	},
	CreateRepository: func(params CreateRepositoryParams) Repository {
		metaData := createMetaData()
		rep := createRepository(metaData, params.EntityRepository)
		return rep
	},
	ToData: func(wal Wallet) *Data {
		return toData(wal)
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
			walletRepository := createRepository(metaData, params.EntityRepository)

			// retrieve all the wallets:
			allWalPS, allWalPSErr := walletRepository.RetrieveSet(0, params.AmountOfElementsPerList)
			if allWalPSErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while retrieving the wallet entity set: %s", allWalPSErr.Error())
				w.Write([]byte(str))
				return
			}

			// render:
			datSet, datSetErr := toDataSet(allWalPS)
			if datSetErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while converting the wallet entity set to data: %s", datSetErr.Error())
				w.Write([]byte(str))
				return
			}

			w.WriteHeader(http.StatusOK)
			params.Tmpl.Execute(w, datSet)
		}
	},
	Route: func(params RouteParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			// create the metadata:
			metaData := createMetaData()

			// create the repositories:
			walletRepository := createRepository(metaData, params.EntityRepository)

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

				// retrieve the wallet by id:
				wal, walErr := walletRepository.RetrieveByID(&id)
				if idErr != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(walErr.Error()))
					return
				}

				w.WriteHeader(http.StatusOK)
				params.Tmpl.Execute(w, toData(wal))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the ID could not be found")
			w.Write([]byte(str))
		}
	},
}
