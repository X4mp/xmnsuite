package xmn

import (
	"unsafe"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

type dependencies struct {
	genService GenesisService
}

func createXMN(
	genService GenesisService,
	namespace string,
	name string,
	id *uuid.UUID,
	fromBlockIndex int64,
	toBlockIndex int64,
	version string,
	rootDir string,
	routerDS datastore.DataStore,
	routerRoleKey string,
) applications.Application {

	app := applications.SDKFunc.CreateApplication(applications.CreateApplicationParams{
		Namespace:      namespace,
		Name:           name,
		ID:             id,
		FromBlockIndex: fromBlockIndex,
		ToBlockIndex:   toBlockIndex,
		Version:        version,
		DirPath:        rootDir,
		RouterParams: routers.CreateRouterParams{
			DataStore: routerDS,
			RoleKey:   routerRoleKey,
			RtesParams: []routers.CreateRouteParams{
				saveGenesis(),
				retrieveGenesis(),
			},
		},
	})

	return app
}

/*
 * Create Dependencies
 */
func createDependencies(store datastore.DataStore) *dependencies {
	walletService := createWalletService(store)
	genesisService := createGenesisService(store, walletService)
	out := dependencies{
		genService: genesisService,
	}

	return &out
}

/*
 * Save Genesis
 */

func saveGenesis() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create dependencies:
			dependencies := createDependencies(store)

			// unmarshal data:
			gen := new(genesis)
			jsErr := cdc.UnmarshalJSON(data, gen)
			if jsErr != nil {
				return nil, jsErr
			}

			// save the genesis instance:
			saveGenErr := dependencies.genService.Save(gen)
			if saveGenErr != nil {
				return nil, saveGenErr
			}

			// convert the Genesis to json:
			jsData, jsDataErr := cdc.MarshalJSON(gen)
			if jsDataErr != nil {
				return nil, jsDataErr
			}

			// create the gaz price:
			gazUsed := int(unsafe.Sizeof(jsData)) * gen.GzPricePerKb

			// return the response:
			resp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
				Code:    routers.IsSuccessful,
				Log:     "success",
				GazUsed: int64(gazUsed),
				Tags: map[string][]byte{
					path: jsData,
				},
			})

			return resp, nil
		},
	}
}

/*
 * Retrieve Genesis
 */
func retrieveGenesis() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/",
		QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
			// create dependencies:
			dependencies := createDependencies(store)

			// retrieve the genesis instance:
			retGen, retGenErr := dependencies.genService.Retrieve()
			if retGenErr != nil {
				return nil, retGenErr
			}

			// convert the genesis to json:
			js, jsErr := cdc.MarshalJSON(retGen)
			if jsErr != nil {
				return nil, jsErr
			}

			// return the response:
			resp := routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
				Code:  routers.IsSuccessful,
				Log:   "success",
				Key:   path,
				Value: js,
			})

			return resp, nil
		},
	}
}
