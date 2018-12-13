package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/request/vote"
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
	routePrefix                   string
	routerRoleKey                 string
	meta                          meta.Meta
	maxAmountOfEntitiesToRetrieve int
}

func createCore20181108(met meta.Meta, routePrefix string, routerRoleKey string, maxAmountOfEntitiesToRetrieve int) *core20181108 {

	out := core20181108{
		routePrefix:   routePrefix,
		routerRoleKey: routerRoleKey,
		meta:          met,
		maxAmountOfEntitiesToRetrieve: maxAmountOfEntitiesToRetrieve,
	}

	return &out
}

func create20181106WithRootPubKey(
	namespace string,
	name string,
	id *uuid.UUID,
	fromBlockIndex int64,
	toBlockIndex int64,
	rootDir string,
	routePrefix string,
	routerRoleKey string,
	ds datastore.StoredDataStore,
	met meta.Meta,
	rootPubKey crypto.PublicKey,
	maxAmountOfEntitiesToRetrieve int,
) applications.Application {
	// enable the root user to have write access to the genesis route:
	store := ds.DataStore()
	store.Users().Insert(rootPubKey)
	store.Roles().Add(routerRoleKey, rootPubKey)
	store.Roles().EnableWriteAccess(routerRoleKey, fmt.Sprintf("%s/genesis", routePrefix))

	return create20181106(namespace, name, id, fromBlockIndex, toBlockIndex, rootDir, routePrefix, routerRoleKey, ds, met, maxAmountOfEntitiesToRetrieve)
}

func create20181106(
	namespace string,
	name string,
	id *uuid.UUID,
	fromBlockIndex int64,
	toBlockIndex int64,
	rootDir string,
	routePrefix string,
	routerRoleKey string,
	ds datastore.StoredDataStore,
	met meta.Meta,
	maxAmountOfEntitiesToRetrieve int,
) applications.Application {
	// create core:
	core := createCore20181108(met, routePrefix, routerRoleKey, maxAmountOfEntitiesToRetrieve)

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
			genRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
				Datastore: ds,
			})

			gen, genErr := genRepository.Retrieve()
			if genErr != nil {
				return nil, genErr
			}

			// retrieve the validators:
			validatorRepository := validator.SDKFunc.CreateRepository(ds)
			valPS, valPSErr := validatorRepository.RetrieveSet(0, gen.MaxAmountOfValidators())
			if valPSErr != nil {
				return nil, valPSErr
			}

			// create the application validators:
			valsIns := valPS.Instances()
			appVals := []applications.Validator{}
			for _, oneValIns := range valsIns {
				oneVal := oneValIns.(validator.Validator)
				appVals = append(appVals, applications.SDKFunc.CreateValidator(applications.CreateValidatorParams{
					IP:     oneVal.IP(),
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
				core.retrieveByIntersectKeynames(),
				core.retrieveSetByIntersectKeynames(),
				core.deleteEntityByID(),
				core.saveRequest(),
				core.saveEntityRequestVote(),
			},
		},
	})

	return app
}

func (app *core20181108) saveGenesis() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: fmt.Sprintf("%s/genesis", app.routePrefix),
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// converts the data to a genesis instance:
			ins, insErr := app.meta.Genesis().MetaData().ToEntity()(dep.entityRepository, data)
			if insErr != nil {
				return nil, insErr
			}

			if gen, ok := ins.(genesis.Genesis); ok {
				// save the genesis instance:
				saveErr := dep.genesisService.Save(gen)
				if saveErr != nil {
					return nil, saveErr
				}

				// enable the route to save instances:
				store.Roles().EnableWriteAccess(app.routerRoleKey, fmt.Sprintf("%s/[a-z-]+", app.routePrefix))

				// enable the route to save requests:
				store.Roles().EnableWriteAccess(app.routerRoleKey, fmt.Sprintf("%s/[a-z-]+/requests", app.routePrefix))

				// convert to json:
				normalized, normalizedErr := app.meta.Genesis().MetaData().Normalize()(gen)
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
		Pattern: fmt.Sprintf("%s/<name|[a-z-]+>", app.routePrefix),
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dep.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Genesis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the name:
			entityRepresentations := app.meta.Write()
			if name, ok := params["name"]; ok {
				// retrieve the entity representation:
				if representation, ok := entityRepresentations[name]; ok {
					// converts the data to an entity:
					ins, insErr := representation.MetaData().ToEntity()(dep.entityRepository, data)
					if insErr != nil {
						return nil, insErr
					}

					// make sure the entity does not already exists:
					_, alreadyExistsErr := dep.entityRepository.RetrieveByID(representation.MetaData(), ins.ID())
					if alreadyExistsErr == nil {
						str := fmt.Sprintf("the entity (Name: %s, ID: %s) already exists and therefore cannot be updated directly", representation.MetaData().Name(), ins.ID().String())
						return nil, errors.New(str)
					}

					// save the entity:
					saveErr := dep.entityService.Save(ins, representation)
					if saveErr != nil {
						return nil, saveErr
					}

					// enable the ability to update/delete the entity:
					store.Roles().EnableWriteAccess(app.routerRoleKey, fmt.Sprintf("%s/%s/%s", app.routePrefix, name, ins.ID().String()))

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
		Pattern: fmt.Sprintf("%s/<name|[a-z-]+>/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", app.routePrefix),
		QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the name:
			entityMetaDatas := app.meta.Retrieval()
			if name, ok := params["name"]; ok {
				// retrieve the entity metadata:
				if metaData, ok := entityMetaDatas[name]; ok {
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

func (app *core20181108) retrieveByIntersectKeynames() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: fmt.Sprintf("%s/<name|[a-z-]+>/<keynames|[a-z0-9-+]+>/intersect", app.routePrefix),
		QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the name:
			entityMetaDatas := app.meta.Retrieval()
			if name, ok := params["name"]; ok {
				// retrieve the entity metadata:
				if metaData, ok := entityMetaDatas[name]; ok {
					// parse the keynames:
					keynames := strings.Split(params["keynames"], "|")

					// retrieve the entity instance:
					retIns, retInsErr := dep.entityRepository.RetrieveByIntersectKeynames(metaData, keynames)
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

func (app *core20181108) retrieveSetByIntersectKeynames() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: fmt.Sprintf("%s/<name|[a-z-]+>/<keynames|[a-z0-9-+]+>/set/intersect", app.routePrefix),
		QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {
			index := 0
			if indexAsString, ok := params["index"]; ok {
				idx, idxErr := strconv.Atoi(indexAsString)
				if idxErr != nil {
					str := fmt.Sprintf("the given index (%s) is invalid: %s", indexAsString, idxErr.Error())
					return nil, errors.New(str)
				}

				index = idx
			}

			amount := app.maxAmountOfEntitiesToRetrieve
			if amountAsString, ok := params["amount"]; ok {
				am, amErr := strconv.Atoi(amountAsString)
				if amErr != nil {
					str := fmt.Sprintf("the given amount (%s) is invalid: %s", amountAsString, amErr.Error())
					return nil, errors.New(str)
				}

				amount = am
			}

			if amount > app.maxAmountOfEntitiesToRetrieve {
				amount = app.maxAmountOfEntitiesToRetrieve
			}

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the name:
			entityMetaDatas := app.meta.Retrieval()
			if name, ok := params["name"]; ok {
				// retrieve the entity metadata:
				if metaData, ok := entityMetaDatas[name]; ok {
					// parse the keynames:
					keynames := strings.Split(params["keynames"], "|")

					// retrieve the entity partial set:
					retPS, retPSErr := dep.entityRepository.RetrieveSetByIntersectKeynames(metaData, keynames, index, amount)
					if retPSErr != nil {
						return nil, retPSErr
					}

					// normalize:
					normalized := entity.SDKFunc.NormalizePartialSet(entity.NormalizePartialSetParams{
						PartialSet: retPS,
						MetaData:   metaData,
					})

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
		Pattern: fmt.Sprintf("%s/<name|[a-z-]+>/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", app.routePrefix),
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
			entityRepresentations := app.meta.Write()
			if name, ok := params["name"]; ok {
				// retrieve the entity representation:
				if representation, ok := entityRepresentations[name]; ok {
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
		Pattern: fmt.Sprintf("%s/<name|[a-z-]+>/requests", app.routePrefix),
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the genesis:
			gen, genErr := dep.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Genesis instance: %s", genErr.Error())
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
			walIns, walInsErr := dep.entityRepository.RetrieveByID(app.meta.Wallet().MetaData(), &walletID)
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
					str := fmt.Sprintf("the requester PublicKey (%s) is not a user on the given wallet (ID: %s)", from.String(), wal.ID().String())
					return nil, errors.New(str)
				}

				// retrieve the name:
				allRequestRepresentations := app.meta.WriteOnAllEntityRequest()
				if name, ok := params["name"]; ok {
					// retrieve the entity representation:
					if representation, ok := allRequestRepresentations[name]; ok {
						// converts the data to an entity:
						ins, insErr := representation.MetaData().ToEntity()(dep.entityRepository, ptr.EntityJSON)
						if insErr != nil {
							return nil, insErr
						}

						// create the request:
						req := request.SDKFunc.Create(request.CreateParams{
							ID:             &reqID,
							FromUser:       usr,
							NewEntity:      ins,
							EntityMetaData: representation.MetaData(),
						})

						// save the request:
						representation := request.SDKFunc.CreateRepresentation()
						saveErr := dep.entityService.Save(req, representation)
						if saveErr != nil {
							return nil, saveErr
						}

						// enable the voting on the request:
						store.Roles().EnableWriteAccess(app.routerRoleKey, fmt.Sprintf("%s/%s/requests/%s/[a-z-]+", app.routePrefix, name, reqID.String()))

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

func (app *core20181108) saveEntityRequestVote() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: fmt.Sprintf("%s/<entityname|[a-z-]+>/requests/<requestid|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>/<name|[a-z-]+>", app.routePrefix),
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {
			entityRequests := app.meta.WriteOnEntityRequest()
			if entityRequest, ok := entityRequests[params["name"]]; ok {
				// create the dependencies:
				dep := createDependencies(store)

				// retrieve the genesis:
				gen, genErr := dep.genesisRepository.Retrieve()
				if genErr != nil {
					str := fmt.Sprintf("there was an error while retrieving the Genesis instance: %s", genErr.Error())
					return nil, errors.New(str)
				}

				// parse the requestID:
				requestID, requestIDErr := uuid.FromString(params["requestid"])
				if requestIDErr != nil {
					str := fmt.Sprintf("the given requestID (%s) is invalid: %s", params["requestid"], requestIDErr.Error())
					return nil, errors.New(str)
				}

				// retrieve the request:
				reqIns, reqInsErr := dep.entityRepository.RetrieveByID(app.meta.Request().MetaData(), &requestID)
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
					if entityName, ok := params["entityname"]; ok {
						// retrieve the representation:
						requestRepresentations := entityRequest.Map()
						if representation, ok := requestRepresentations[entityName]; ok {

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

							saveErr := entityRequest.VoteService(store).Save(voteIns, representation)
							if saveErr != nil {
								return nil, saveErr
							}

							// convert to json:
							storable, storableErr := app.meta.Vote().ToStorable()(voteIns)
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

						str := fmt.Sprintf("the given entityName (%s) is not supported", entityName)
						return nil, errors.New(str)
					}

					return nil, errors.New("an entityName must be provided")
				}

				str := fmt.Sprintf("the entity (ID: %s) was expected to be a request instance", reqIns.ID().String())
				return nil, errors.New(str)
			}

			str := fmt.Sprintf("the name (%s) is not an entity that can be voted on", params["name"])
			return nil, errors.New(str)
		},
	}
}
