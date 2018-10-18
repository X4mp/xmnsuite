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
	genService    GenesisService
	walletService WalletService
}

func createXMN(
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
				saveWallet(),
				retrieveWallets(),
				retrieveWalletByID(),
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
	tokService := createTokenService(store)
	intiialDepositService := createInitialDepositService(store, walletService)
	genesisService := createGenesisService(store, walletService, tokService, intiialDepositService)
	out := dependencies{
		genService:    genesisService,
		walletService: walletService,
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

/*
 * Save Wallet
 */

func saveWallet() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/wallets",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create dependencies:
			dependencies := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dependencies.genService.Retrieve()
			if genErr != nil {
				return nil, genErr
			}

			// unmarshal data:
			wal := new(wallet)
			jsErr := cdc.UnmarshalJSON(data, wal)
			if jsErr != nil {
				return nil, jsErr
			}

			// save the wallet instance:
			saveWalletErr := dependencies.walletService.Save(wal)
			if saveWalletErr != nil {
				return nil, saveWalletErr
			}

			// convert the wallet to json:
			jsData, jsDataErr := cdc.MarshalJSON(wal)
			if jsDataErr != nil {
				return nil, jsDataErr
			}

			// create the gaz price:
			gazUsed := int(unsafe.Sizeof(jsData)) * gen.GazPricePerKb()

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
 * Retrieve Wallets
 */
func retrieveWallets() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/wallets",
		QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
			// create dependencies:
			dependencies := createDependencies(store)

			// retrieve the wallets:
			index := fetchIndex(params)
			amount := fetchAmount(params)
			wals, walsErr := dependencies.walletService.Retrieve(index, amount)
			if walsErr != nil {
				return nil, walsErr
			}

			// convert the wallets to json:
			js, jsErr := cdc.MarshalJSON(wals)
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

/*
 * Retrieve Wallet By ID
 */
func retrieveWalletByID() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/wallets/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
		QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
			// create dependencies:
			dependencies := createDependencies(store)

			// create the ID:
			id, idErr := uuid.FromString(fetchParam(params, "id"))
			if idErr != nil {
				return nil, idErr
			}

			// retrieve the wallet by ID:
			wal, walsErr := dependencies.walletService.RetrieveByID(&id)
			if walsErr != nil {
				resp := routers.SDKFunc.CreateQueryResponse(routers.CreateQueryResponseParams{
					Code:  routers.NotFound,
					Log:   "not found",
					Key:   path,
					Value: []byte(""),
				})

				return resp, nil
			}

			// convert the wallets to json:
			js, jsErr := cdc.MarshalJSON(wal)
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
