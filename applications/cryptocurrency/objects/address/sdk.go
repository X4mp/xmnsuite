package address

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
)

// Address represents an address bind between a wallet and a cryptocurrency address
type Address interface {
	ID() *uuid.UUID
	Wallet() wallet.Wallet
	Address() []byte
}

// Normalized represents a normalized offer
type Normalized interface {
}

// Repository represents the address repository
type Repository interface {
	RetrieveByID(id *uuid.UUID) (Address, error)
	RetrieveByAddress(addr []byte) (Address, error)
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
	Address []byte
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

// SDKFunc represents the Address SDK func
var SDKFunc = struct {
	Create               func(params CreateParams) Address
	CreateMetaData       func() entity.MetaData
	CreateRepresentation func() entity.Representation
	ToData               func(addr Address) *Data
	ToDataSet            func(ps entity.PartialSet) *DataSet
	Route                func(params RouteParams) func(w http.ResponseWriter, r *http.Request)
	RouteSet             func(params RouteSetParams) func(w http.ResponseWriter, r *http.Request)
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
}
