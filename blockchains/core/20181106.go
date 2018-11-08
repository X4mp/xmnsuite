package core

import (
	"errors"
	"fmt"
	"unsafe"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/wallet"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

type incomingVote struct {
	IsApproved bool `json:"is_approved"`
}

type core20181108 struct {
	genesisRepresentation  entity.Representation
	entityMetaDatas        map[string]entity.MetaData
	entityRepresentations  map[string]entity.Representation
	requestRepresentations map[string]entity.Representation
}

func createCore20181108() *core20181108 {

	walletRepresentation := wallet.SDKFunc.CreateRepresentation()

	out := core20181108{
		genesisRepresentation: genesis.SDKFunc.CreateRepresentation(),
		entityMetaDatas: map[string]entity.MetaData{
			"genesis": genesis.SDKFunc.CreateMetaData(),
		},
		entityRepresentations: map[string]entity.Representation{
			"wallet": walletRepresentation,
		},
		requestRepresentations: map[string]entity.Representation{},
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
		RouterParams: routers.CreateRouterParams{
			DataStore: ds.DataStore(),
			RoleKey:   routerRoleKey,
			RtesParams: []routers.CreateRouteParams{
				core.saveGenesis(),
				core.saveEntity(),
				core.retrieveEntityByID(),
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
						return nil, idErr
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

func (app *core20181108) saveRequest() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<name|[a-z-]+>/request/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dep.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Gensis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the from user:
			fromUser, fromUserErr := dep.userRepository.RetrieveByPubKey(from)
			if fromUserErr != nil {
				str := fmt.Sprintf("the from user (pubKey: %s) could not be found", from.String())
				return nil, errors.New(str)
			}

			// retrieve the name:
			if name, ok := params["name"]; ok {
				// retrieve the entity representation:
				if representation, ok := app.requestRepresentations[name]; ok {
					if requestIDAsString, ok := params["id"]; ok {
						requestID, requestIDErr := uuid.FromString(requestIDAsString)
						if requestIDErr != nil {
							return nil, requestIDErr
						}

						// unmarshal the data:
						ins, insErr := representation.MetaData().ToEntity()(dep.entityRepository, data)
						if insErr != nil {
							return nil, insErr
						}

						// build the request:
						req := request.SDKFunc.Create(request.CreateParams{
							ID:        &requestID,
							FromUser:  fromUser,
							NewEntity: ins,
						})

						// save the request:
						saveErr := dep.entityService.Save(req, request.SDKFunc.CreateRepresentation(request.CreateRepresentationParams{
							EntityMetaData: representation.MetaData(),
						}))

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

					return nil, errors.New("the requestID is mandatory")
				}

				str := fmt.Sprintf("the given entity name (%s) is not supported", name)
				return nil, errors.New(str)
			}

			return nil, errors.New("an entity name must be provided")
		},
	}
}

func (app *core20181108) saveRequestVote() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<name|[a-z-]+>/request/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>/vote/<voteID|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {
			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dep.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Gensis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the from user:
			fromUser, fromUserErr := dep.userRepository.RetrieveByPubKey(from)
			if fromUserErr != nil {
				str := fmt.Sprintf("the from user (pubKey: %s) could not be found", from.String())
				return nil, errors.New(str)
			}

			// is approved:
			ptr := new(incomingVote)
			jsErr := cdc.UnmarshalJSON(data, ptr)
			if jsErr != nil {
				return nil, jsErr
			}

			// retrieve the name:
			if name, ok := params["name"]; ok {
				// retireve the representation:
				if representation, ok := app.requestRepresentations[name]; ok {
					// retrieve the requestID:
					if requestIDAsString, ok := params["id"]; ok {
						requestID, requestIDErr := uuid.FromString(requestIDAsString)
						if requestIDErr != nil {
							return nil, requestIDErr
						}

						if voteIDAsString, ok := params["voteID"]; ok {
							voteID, voteIDErr := uuid.FromString(voteIDAsString)
							if voteIDErr != nil {
								return nil, voteIDErr
							}

							// retrieve the request:
							req, reqErr := dep.entityRepository.RetrieveByID(request.SDKFunc.CreateMetaData(request.CreateMetaDataParams{
								EntityMetaData: representation.MetaData(),
							}), &requestID)

							if reqErr != nil {
								return nil, reqErr
							}

							// create the vote:
							voteIns := vote.SDKFunc.Create(vote.CreateParams{
								ID:         &voteID,
								Request:    req.(request.Request),
								Voter:      fromUser,
								IsApproved: ptr.IsApproved,
							})

							// save the vote:
							voteService := vote.SDKFunc.CreateService(vote.CreateServiceParams{
								EntityRepository: dep.entityRepository,
								EntityService:    dep.entityService,
								RequestRepresentation: request.SDKFunc.CreateRepresentation(request.CreateRepresentationParams{
									EntityMetaData: representation.MetaData(),
								}),
								NewEntityRepresentation: representation,
							})

							saveErr := voteService.Save(voteIns)
							if saveErr != nil {
								return nil, saveErr
							}

							// convert to json:
							storable, storableErr := representation.ToStorable()(voteIns)
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

						return nil, errors.New("the voteID is mandatory")
					}

					return nil, errors.New("the requestID is mandatory")
				}

				str := fmt.Sprintf("the given entity name (%s) is not supported", name)
				return nil, errors.New(str)
			}

			return nil, errors.New("an entity name must be provided")
		},
	}
}
