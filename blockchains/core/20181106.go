package core

import (
	"errors"
	"fmt"
	"unsafe"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/request/vote"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

type incomingVote struct {
	ID         string `json:"id"`
	IsApproved bool   `json:"is_approved"`
}

type incomingRequest struct {
	ID         string `json:"id"`
	WalletID   string `json:"wallet_id"`
	EntityJSON []byte `json:"entity"`
}

type core20181108 struct {
	genesisRepresentation  entity.Representation
	walletRepresentation   entity.Representation
	requestRepresentation  entity.Representation
	voteRepresentation     entity.Representation
	entityMetaDatas        map[string]entity.MetaData
	entityRepresentations  map[string]entity.Representation
	requestRepresentations map[string]entity.Representation
}

func createCore20181108() *core20181108 {

	// register the possible requests:
	request.SDKFunc.Register(pledge.SDKFunc.CreateMetaData())

	walletRepresentation := wallet.SDKFunc.CreateRepresentation()
	pledgeRepresentation := pledge.SDKFunc.CreateRepresentation()
	validatorRepresentation := validator.SDKFunc.CreateRepresentation()

	out := core20181108{
		genesisRepresentation: genesis.SDKFunc.CreateRepresentation(),
		walletRepresentation:  walletRepresentation,
		requestRepresentation: request.SDKFunc.CreateRepresentation(),
		voteRepresentation:    vote.SDKFunc.CreateRepresentation(),
		entityMetaDatas: map[string]entity.MetaData{
			"genesis":   genesis.SDKFunc.CreateMetaData(),
			"wallet":    walletRepresentation.MetaData(),
			"validator": validatorRepresentation.MetaData(),
		},
		entityRepresentations: map[string]entity.Representation{
			"wallet":    walletRepresentation,
			"validator": validatorRepresentation,
		},
		requestRepresentations: map[string]entity.Representation{
			"pledge": pledgeRepresentation,
		},
	}

	return &out
}

func create20181106(
	namespace string,
	name string,
	id *uuid.UUID,
	fromBlockIndex int64,
	toBlockIndex int64,
	rootDir string,
	routerRoleKey string,
	rootPubKey crypto.PublicKey,
	ds datastore.StoredDataStore,
) applications.Application {

	// enable the root user to have write access to the genesis route:
	store := ds.DataStore()
	store.Users().Insert(rootPubKey)
	store.Roles().Add(routerRoleKey, rootPubKey)
	store.Roles().EnableWriteAccess(routerRoleKey, "/genesis")
	store.Roles().EnableWriteAccess(routerRoleKey, "/[a-z-]+")
	store.Roles().EnableWriteAccess(routerRoleKey, "/[a-z-]+/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	store.Roles().EnableWriteAccess(routerRoleKey, "/[a-z-]+/requests")
	store.Roles().EnableWriteAccess(routerRoleKey, "/[a-z-]+/requests/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")

	// create core:
	core := createCore20181108()

	// create application:
	version := "2018.11.06"
	app := applications.SDKFunc.CreateApplication(applications.CreateApplicationParams{
		Namespace:      namespace,
		Name:           name,
		ID:             id,
		FromBlockIndex: fromBlockIndex,
		ToBlockIndex:   toBlockIndex,
		Version:        version,
		DirPath:        rootDir,
		Store:          ds,
		RetrieveValidators: func(ds datastore.DataStore) ([]applications.Validator, error) {
			// retrieve the genesis:
			genRepository := genesis.SDKFunc.CreateRepository(ds)
			gen, genErr := genRepository.Retrieve()
			if genErr != nil {
				return nil, genErr
			}

			// retrieve the validators:
			validatorRepository := validator.SDKFunc.CreateRepository(ds)
			valPS, valPSErr := validatorRepository.RetrieveSet(gen.MaxAmountOfValidators())
			if valPSErr != nil {
				return nil, valPSErr
			}

			// create the application validators:
			valsIns := valPS.Instances()
			appVals := []applications.Validator{}
			for _, oneValIns := range valsIns {
				oneVal := oneValIns.(validator.Validator)
				appVals = append(appVals, applications.SDKFunc.CreateValidator(applications.CreateValidatorParams{
					PubKey: oneVal.PubKey(),
					Power:  int64(oneVal.Pledge().From().Amount()),
				}))
			}

			return appVals, nil
		},
		RouterParams: routers.CreateRouterParams{
			DataStore: ds.DataStore(),
			RoleKey:   routerRoleKey,
			RtesParams: []routers.CreateRouteParams{
				core.saveGenesis(),
				core.saveEntity(),
				core.retrieveEntityByID(),
				core.deleteEntityByID(),
				core.saveRequest(),
				core.saveRequestVote(),
			},
		},
	})

	return app
}

func (app *core20181108) saveGenesis() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/genesis",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// converts the data to a genesis instance:
			ins, insErr := app.genesisRepresentation.MetaData().ToEntity()(dep.entityRepository, data)
			if insErr != nil {
				return nil, insErr
			}

			if gen, ok := ins.(genesis.Genesis); ok {
				// save the genesis instance:
				saveErr := dep.genesisService.Save(gen)
				if saveErr != nil {
					return nil, saveErr
				}

				// convert to json:
				normalized, normalizedErr := app.genesisRepresentation.MetaData().Normalize()(gen)
				if normalizedErr != nil {
					return nil, normalizedErr
				}

				jsData, jsDataErr := cdc.MarshalJSON(normalized)
				if jsDataErr != nil {
					return nil, jsDataErr
				}

				// there is no gaz cost for the genesis:
				gazUsed := 0

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
			}

			return nil, errors.New("the given data is not a normalized representation of aa Genesis instance")
		},
	}
}

func (app *core20181108) saveEntity() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<name|[a-z-]+>",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dep.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Gensis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the name:
			if name, ok := params["name"]; ok {
				// retrieve the entity representation:
				if representation, ok := app.entityRepresentations[name]; ok {
					// converts the data to an entity:
					ins, insErr := representation.MetaData().ToEntity()(dep.entityRepository, data)
					if insErr != nil {
						return nil, insErr
					}

					// save the entity:
					saveErr := dep.entityService.Save(ins, representation)
					if saveErr != nil {
						return nil, saveErr
					}

					// convert to json:
					storable, storableErr := representation.ToStorable()(ins)
					if storableErr != nil {
						return nil, storableErr
					}

					jsData, jsDataErr := cdc.MarshalJSON(storable)
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
				}

				str := fmt.Sprintf("the given entity name (%s) is not supported", name)
				return nil, errors.New(str)
			}

			return nil, errors.New("an entity name must be provided")
		},
	}
}

func (app *core20181108) retrieveEntityByID() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<name|[a-z-]+>/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
		QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the name:
			if name, ok := params["name"]; ok {
				// retrieve the entity metadata:
				if metaData, ok := app.entityMetaDatas[name]; ok {
					// parse the id:
					id, idErr := uuid.FromString(params["id"])
					if idErr != nil {
						str := fmt.Sprintf("the given ID (%s) is invalid: %s", params["id"], idErr.Error())
						return nil, errors.New(str)
					}

					// retrieve the entity instance:
					retIns, retInsErr := dep.entityRepository.RetrieveByID(metaData, &id)
					if retInsErr != nil {
						return nil, retInsErr
					}

					// normalize:
					normalized, normalizedErr := metaData.Normalize()(retIns)
					if normalizedErr != nil {
						return nil, normalizedErr
					}

					// convert the normalized entity to json:
					js, jsErr := cdc.MarshalJSON(normalized)
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
				}

				str := fmt.Sprintf("the given entity name (%s) is not supported", name)
				return nil, errors.New(str)
			}

			return nil, errors.New("an entity name must be provided")
		},
	}
}

func (app *core20181108) deleteEntityByID() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<name|[a-z-]+>/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
		DelTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.TransactionResponse, error) {
			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dep.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Genesis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the name:
			if name, ok := params["name"]; ok {
				// retrieve the entity representation:
				if representation, ok := app.entityRepresentations[name]; ok {
					// get the metadata:
					metaData := representation.MetaData()

					// parse the id:
					id, idErr := uuid.FromString(params["id"])
					if idErr != nil {
						str := fmt.Sprintf("the given ID (%s) is invalid: %s", params["id"], idErr.Error())
						return nil, errors.New(str)
					}

					// retrieve the entity instance:
					retIns, retInsErr := dep.entityRepository.RetrieveByID(metaData, &id)
					if retInsErr != nil {
						return nil, retInsErr
					}

					// delete the entity instance:
					delErr := dep.entityService.Delete(retIns, representation)
					if delErr != nil {
						return nil, delErr
					}

					// normalize:
					normalized, normalizedErr := metaData.Normalize()(retIns)
					if normalizedErr != nil {
						return nil, normalizedErr
					}

					// convert the normalized entity to json:
					js, jsErr := cdc.MarshalJSON(normalized)
					if jsErr != nil {
						return nil, jsErr
					}

					// calculate the gaz used:
					gazUsed := int(unsafe.Sizeof(js)) * gen.GazPricePerKb()

					// return the response:
					resp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
						Code:    routers.IsSuccessful,
						Log:     "success",
						GazUsed: int64(gazUsed),
						Tags: map[string][]byte{
							path: []byte("deleted"),
						},
					})

					return resp, nil
				}

				str := fmt.Sprintf("the given entity name (%s) is not supported", name)
				return nil, errors.New(str)
			}

			return nil, errors.New("an entity name must be provided")
		},
	}
}

func (app *core20181108) saveRequest() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<name|[a-z-]+>/requests",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dep.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Gensis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// convert the data to the incoming request:
			ptr := new(incomingRequest)
			jsErr := cdc.UnmarshalJSON(data, ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			// parse the walletID:
			walletID, walletIDErr := uuid.FromString(ptr.WalletID)
			if walletIDErr != nil {
				str := fmt.Sprintf("the given walletID (%s) is invalid: %s", ptr.WalletID, walletIDErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the wallet:
			walIns, walInsErr := dep.entityRepository.RetrieveByID(app.walletRepresentation.MetaData(), &walletID)
			if walInsErr != nil {
				return nil, walInsErr
			}

			// parse the requestID:
			reqID, reqIDErr := uuid.FromString(ptr.ID)
			if reqIDErr != nil {
				str := fmt.Sprintf("the given requestID (%s) is invalid: %s", ptr.ID, reqIDErr.Error())
				return nil, errors.New(str)
			}

			if wal, ok := walIns.(wallet.Wallet); ok {
				// retrieve the user:
				usr, usrErr := dep.userRepository.RetrieveByPubKeyAndWallet(from, wal)
				if usrErr != nil {
					return nil, usrErr
				}

				// retrieve the name:
				if name, ok := params["name"]; ok {
					// retrieve the entity representation:
					if representation, ok := app.requestRepresentations[name]; ok {
						// converts the data to an entity:
						ins, insErr := representation.MetaData().ToEntity()(dep.entityRepository, ptr.EntityJSON)
						if insErr != nil {
							return nil, insErr
						}

						// create the request:
						req := request.SDKFunc.Create(request.CreateParams{
							ID:        &reqID,
							FromUser:  usr,
							NewEntity: ins,
						})

						// save the request:
						representation := request.SDKFunc.CreateRepresentation()
						saveErr := dep.entityService.Save(req, representation)
						if saveErr != nil {
							return nil, saveErr
						}

						// convert to json:
						storable, storableErr := representation.ToStorable()(req)
						if storableErr != nil {
							return nil, storableErr
						}

						jsData, jsDataErr := cdc.MarshalJSON(storable)
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
					}

					str := fmt.Sprintf("the given entity name (%s) is not supported for requests", name)
					return nil, errors.New(str)
				}

				return nil, errors.New("an entity name must be provided")
			}

			str := fmt.Sprintf("the entity (ID: %s) was expected to be a wallet instance", walIns.ID().String())
			return nil, errors.New(str)
		},
	}
}

func (app *core20181108) saveRequestVote() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<name|[a-z-]+>/requests/<requestid|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {
			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dep.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Gensis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// parse the requestID:
			requestID, requestIDErr := uuid.FromString(params["requestid"])
			if requestIDErr != nil {
				str := fmt.Sprintf("the given requestID (%s) is invalid: %s", params["requestid"], requestIDErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the request:
			reqIns, reqInsErr := dep.entityRepository.RetrieveByID(app.requestRepresentation.MetaData(), &requestID)
			if reqInsErr != nil {
				return nil, reqInsErr
			}

			if req, ok := reqIns.(request.Request); ok {
				// convert the data to the incoming vote:
				ptr := new(incomingVote)
				jsErr := cdc.UnmarshalJSON(data, ptr)
				if jsErr != nil {
					return nil, jsErr
				}

				// retrieve the name:
				if name, ok := params["name"]; ok {
					// retrieve the representation:
					if representation, ok := app.requestRepresentations[name]; ok {

						// retrieve the voter:
						voter, voterErr := dep.userRepository.RetrieveByPubKeyAndWallet(from, req.From().Wallet())
						if voterErr != nil {
							return nil, voterErr
						}

						voteID, voteIDErr := uuid.FromString(ptr.ID)
						if voteIDErr != nil {
							str := fmt.Sprintf("the given voteID (%s) is invalid: %s", ptr.ID, voteIDErr.Error())
							return nil, errors.New(str)
						}

						// create the vote:
						voteIns := vote.SDKFunc.Create(vote.CreateParams{
							ID:         &voteID,
							Request:    req,
							Voter:      voter,
							IsApproved: ptr.IsApproved,
						})

						// save the vote:
						voteService := vote.SDKFunc.CreateService(vote.CreateServiceParams{
							EntityRepository: dep.entityRepository,
							EntityService:    dep.entityService,
						})

						saveErr := voteService.Save(voteIns, representation)
						if saveErr != nil {
							return nil, saveErr
						}

						// convert to json:
						storable, storableErr := app.voteRepresentation.ToStorable()(voteIns)
						if storableErr != nil {
							return nil, storableErr
						}

						jsData, jsDataErr := cdc.MarshalJSON(storable)
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

					}

					str := fmt.Sprintf("the given entity name (%s) is not supported", name)
					return nil, errors.New(str)
				}

				return nil, errors.New("an entity name must be provided")
			}

			str := fmt.Sprintf("the entity (ID: %s) was expected to be a request instance", reqIns.ID().String())
			return nil, errors.New(str)

		},
	}
}
