package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/affiliates"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/fees"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/datastore"
	"github.com/xmnservices/xmnsuite/routers"
)

const maxAmountOfEntitiesToRetrieve = 500

type incomingVote struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Reason     string `json:"reason"`
	IsNeutral  bool   `json:"is_neutral"`
	IsApproved bool   `json:"is_approved"`
}

type incomingRequest struct {
	ID               string `json:"id"`
	Reason           string `json:"reason"`
	WalletID         string `json:"wallet_id"`
	SaveEntityJSON   []byte `json:"save_entity_json"`
	DeleteEntityJSON []byte `json:"delete_entity_json"`
}

type core20181108 struct {
	routerRoleKey string
	meta          meta.Meta
}

func createCore20181108(met meta.Meta, routerRoleKey string) *core20181108 {

	out := core20181108{
		routerRoleKey: routerRoleKey,
		meta:          met,
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
	routerRoleKey string,
	ds datastore.StoredDataStore,
	met meta.Meta,
	rootPubKey crypto.PublicKey,
) applications.Application {
	// enable the root user to have write access to the genesis route:
	store := ds.DataStore()
	store.Users().Insert(rootPubKey)
	store.Roles().Add(routerRoleKey, rootPubKey)
	store.Roles().EnableWriteAccess(routerRoleKey, "/genesis")

	return create20181106(namespace, name, id, fromBlockIndex, toBlockIndex, rootDir, routerRoleKey, ds, met)
}

func create20181106(
	namespace string,
	name string,
	id *uuid.UUID,
	fromBlockIndex int64,
	toBlockIndex int64,
	rootDir string,
	routerRoleKey string,
	ds datastore.StoredDataStore,
	met meta.Meta,
) applications.Application {
	// create core:
	core := createCore20181108(met, routerRoleKey)

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
			validatorRepository := validator.SDKFunc.CreateRepository(validator.CreateRepositoryParams{
				Store: ds,
			})

			vals, valsErr := validatorRepository.RetrieveSetOrderedByPledgeAmount(0, gen.Info().MaxAmountOfValidators())
			if valsErr != nil {
				return nil, valsErr
			}

			// create the application validators:
			appVals := []applications.Validator{}
			for _, oneValIns := range vals {
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
		Pattern: "/genesis",
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

				// enable the route to save accounts:
				store.Roles().EnableWriteAccess(app.routerRoleKey, "/account")

				// enable the route to save instances:
				store.Roles().EnableWriteAccess(app.routerRoleKey, "/[a-z-]+")

				// enable the route to save requests:
				store.Roles().EnableWriteAccess(app.routerRoleKey, "/[a-z-]+/requests")

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
		Pattern: "/<name|[a-z-]+>",
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
					store.Roles().EnableWriteAccess(app.routerRoleKey, fmt.Sprintf("/%s/%s", name, ins.ID().String()))

					// convert to json:
					storable, storableErr := representation.ToStorable()(ins)
					if storableErr != nil {
						return nil, storableErr
					}

					jsData, jsDataErr := cdc.MarshalJSON(storable)
					if jsDataErr != nil {
						return nil, jsDataErr
					}

					// retrieve the client:
					client, clientErr := dep.userRepository.RetrieveByPubKey(from)
					if clientErr != nil {
						return nil, clientErr
					}

					// retrieve the affiliate:
					var aff affiliates.Affiliate
					if client.HasBeenReferred() {
						aff, _ = dep.affiliateRepository.RetrieveByWallet(client.Referral())
					}

					vals, valsErr := dep.validatorRepository.RetrieveSetOrderedByPledgeAmount(0, gen.Info().MaxAmountOfValidators())
					if valsErr != nil {
						return nil, valsErr
					}

					// create the fees:
					fee := fees.SDKFunc.Create(fees.CreateParams{
						Gen:        gen,
						StoredData: jsData,
						Client:     client,
						Affiliate:  aff,
						Validators: vals,
					})

					// save the fees:
					saveFeesErr := dep.entityService.Save(fee, fees.SDKFunc.CreateRepresentation())
					if saveFeesErr != nil {
						return nil, saveFeesErr
					}

					// return the response:
					resp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
						Code:    routers.IsSuccessful,
						Log:     "success",
						GazUsed: int64(fee.Client().Amount()),
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
		Pattern: "/<name|[a-z-]+>/<keynames|[^/]+>/intersect",
		QueryTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, sig crypto.Signature) (routers.QueryResponse, error) {

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the name:
			entityMetaDatas := app.meta.Retrieval()
			if name, ok := params["name"]; ok {
				// retrieve the entity metadata:
				if metaData, ok := entityMetaDatas[name]; ok {
					// decode the keynames:
					keynamesList, keynamesListErr := base64.StdEncoding.DecodeString(params["keynames"])
					if keynamesListErr != nil {
						return nil, keynamesListErr
					}

					// create the slice:
					keynames := strings.Split(string(keynamesList), ",")

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
		Pattern: "/<name|[a-z-]+>/<keynames|[^/]+>/set/intersect",
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

			amount := maxAmountOfEntitiesToRetrieve
			if amountAsString, ok := params["amount"]; ok {
				am, amErr := strconv.Atoi(amountAsString)
				if amErr != nil {
					str := fmt.Sprintf("the given amount (%s) is invalid: %s", amountAsString, amErr.Error())
					return nil, errors.New(str)
				}

				amount = am
			}

			if amount > maxAmountOfEntitiesToRetrieve {
				amount = maxAmountOfEntitiesToRetrieve
			}

			// create the dependencies:
			dep := createDependencies(store)

			// retrieve the name:
			entityMetaDatas := app.meta.Retrieval()
			if name, ok := params["name"]; ok {
				// retrieve the entity metadata:
				if metaData, ok := entityMetaDatas[name]; ok {
					// decode the keynames:
					keynamesList, keynamesListErr := base64.StdEncoding.DecodeString(params["keynames"])
					if keynamesListErr != nil {
						return nil, keynamesListErr
					}

					// create the slice:
					keynames := strings.Split(string(keynamesList), ",")

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
					gazUsed := int(unsafe.Sizeof(js)) * gen.Info().GazPricePerKb()

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
		Pattern: "/<keyname|[a-z-]+>/requests",
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

			// retrieve the user:
			usr, usrErr := dep.userRepository.RetrieveByPubKey(from)
			if usrErr != nil {
				str := fmt.Sprintf("the requester PublicKey (%s) is not a valid user", from.String())
				return nil, errors.New(str)
			}

			if wal, ok := walIns.(wallet.Wallet); ok {
				if keynameName, ok := params["keyname"]; ok {
					// retrieve the keyname by name:
					kname, knameErr := dep.keynameRepository.RetrieveByName(keynameName)
					if knameErr != nil {
						str := fmt.Sprintf("the keyname (name: %s) is invalid: %s", keynameName, knameErr.Error())
						return nil, errors.New(str)
					}

					// retrieve the representation:
					wrOnEntityReq := app.meta.WriteOnEntityRequest()
					if wrReq, ok := wrOnEntityReq[kname.Group().Name()]; ok {
						mp := wrReq.Map()
						if representation, ok := mp[kname.Name()]; ok {

							// instances:
							var toSaveIns entity.Entity
							var toDelIns entity.Entity

							// if the request is a save entity:
							if ptr.SaveEntityJSON != nil {
								saveIns, saveInsErr := representation.MetaData().ToEntity()(dep.entityRepository, ptr.SaveEntityJSON)
								if saveInsErr != nil {
									return nil, saveInsErr
								}

								toSaveIns = saveIns
							}

							// if the request is a delete entity:
							if ptr.DeleteEntityJSON != nil {
								delIns, delInsErr := representation.MetaData().ToEntity()(dep.entityRepository, ptr.DeleteEntityJSON)
								if delInsErr != nil {
									return nil, delInsErr
								}

								toDelIns = delIns
							}

							// create the request:
							req := request.SDKFunc.Create(request.CreateParams{
								ID:           &reqID,
								FromUser:     usr,
								SaveEntity:   toSaveIns,
								DeleteEntity: toDelIns,
								Reason:       ptr.Reason,
								Keyname:      kname,
							})

							// create the active request:
							var activeReq active_request.Request
							keyname := wrReq.RequestedBy().MetaData().Keyname()
							if keyname == token.SDKFunc.CreateMetaData().Keyname() {
								activeReq = active_request.SDKFunc.Create(active_request.CreateParams{
									Request:         req,
									ConcensusNeeded: gen.Info().ConcensusNeeded(),
								})
							}

							if keyname == wallet.SDKFunc.CreateMetaData().Keyname() {
								activeReq = active_request.SDKFunc.Create(active_request.CreateParams{
									Request:         req,
									ConcensusNeeded: wal.ConcensusNeeded(),
								})
							}

							// save the active request:
							activeRepresentation := active_request.SDKFunc.CreateRepresentation()
							saveActiveReqErr := dep.entityService.Save(activeReq, activeRepresentation)
							if saveActiveReqErr != nil {
								return nil, saveActiveReqErr
							}

							// enable the voting on the request:
							store.Roles().EnableWriteAccess(app.routerRoleKey, fmt.Sprintf("/%s/requests/%s", kname.Name(), activeReq.ID().String()))

							// convert to json:
							storable, storableErr := activeRepresentation.ToStorable()(activeReq)
							if storableErr != nil {
								return nil, storableErr
							}

							jsData, jsDataErr := cdc.MarshalJSON(storable)
							if jsDataErr != nil {
								return nil, jsDataErr
							}

							// retrieve the client:
							client, clientErr := dep.userRepository.RetrieveByPubKey(from)
							if clientErr != nil {
								return nil, clientErr
							}

							// retrieve the affiliate:
							var aff affiliates.Affiliate
							if client.HasBeenReferred() {
								aff, _ = dep.affiliateRepository.RetrieveByWallet(client.Referral())
							}

							vals, valsErr := dep.validatorRepository.RetrieveSetOrderedByPledgeAmount(0, gen.Info().MaxAmountOfValidators())
							if valsErr != nil {
								return nil, valsErr
							}

							// create the fees:
							fee := fees.SDKFunc.Create(fees.CreateParams{
								Gen:        gen,
								StoredData: jsData,
								Client:     client,
								Affiliate:  aff,
								Validators: vals,
							})

							// save the fees:
							saveFeesErr := dep.entityService.Save(fee, fees.SDKFunc.CreateRepresentation())
							if saveFeesErr != nil {
								return nil, saveFeesErr
							}

							// return the response:
							resp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
								Code:    routers.IsSuccessful,
								Log:     "success",
								GazUsed: int64(fee.Client().Amount()),
								Tags: map[string][]byte{
									path: jsData,
								},
							})

							return resp, nil
						}

						str := fmt.Sprintf("the keyname (%s) is not supported on the group (%s) for requests", kname.Name(), kname.Group().Name())
						return nil, errors.New(str)
					}

					str := fmt.Sprintf("the group (%s) is not supported for requests", kname.Group().Name())
					return nil, errors.New(str)
				}

				return nil, errors.New("the keyname is mandatory")
			}

			str := fmt.Sprintf("the entity (ID: %s) was expected to be a wallet instance", walIns.ID().String())
			return nil, errors.New(str)
		},
	}
}

func (app *core20181108) saveEntityRequestVote() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<keyname|[a-z-]+>/requests/<requestid|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {
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

			// retrieve the active request:
			reqIns, reqInsErr := dep.entityRepository.RetrieveByID(app.meta.Request().MetaData(), &requestID)
			if reqInsErr != nil {
				return nil, reqInsErr
			}

			if req, ok := reqIns.(active_request.Request); ok {
				entityRequests := app.meta.WriteOnEntityRequest()
				if entityRequest, ok := entityRequests[req.Request().Keyname().Group().Name()]; ok {
					if keynameName, ok := params["keyname"]; ok {
						// convert the data to the incoming vote:
						ptr := new(incomingVote)
						jsErr := cdc.UnmarshalJSON(data, ptr)
						if jsErr != nil {
							return nil, jsErr
						}

						voterID, voterIDErr := uuid.FromString(ptr.UserID)
						if voterIDErr != nil {
							return nil, voterIDErr
						}

						// retrieve the voter:
						voter, voterErr := dep.userRepository.RetrieveByID(&voterID)
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
							Reason:     ptr.Reason,
							IsNeutral:  ptr.IsNeutral,
							IsApproved: ptr.IsApproved,
						})

						representations := entityRequest.Map()
						if representation, ok := representations[keynameName]; ok {
							// create the active vote:
							var activeVote active_vote.Vote
							keyname := entityRequest.RequestedBy().MetaData().Keyname()
							if keyname == token.SDKFunc.CreateMetaData().Keyname() {
								balance, balanceErr := dep.balanceRepository.RetrieveByWallet(voter.Wallet())
								if balanceErr != nil {
									return nil, balanceErr
								}

								activeVote = active_vote.SDKFunc.Create(active_vote.CreateParams{
									Vote:  voteIns,
									Power: balance.Amount(),
								})
							}

							if keyname == wallet.SDKFunc.CreateMetaData().Keyname() {
								activeVote = active_vote.SDKFunc.Create(active_vote.CreateParams{
									Vote:  voteIns,
									Power: voter.Shares(),
								})
							}

							saveErr := dep.voteService.Save(activeVote, representation)
							if saveErr != nil {
								return nil, saveErr
							}

							// convert to json:
							storable, storableErr := app.meta.Vote().ToStorable()(activeVote)
							if storableErr != nil {
								return nil, storableErr
							}

							jsData, jsDataErr := cdc.MarshalJSON(storable)
							if jsDataErr != nil {
								return nil, jsDataErr
							}

							// retrieve the client:
							client, clientErr := dep.userRepository.RetrieveByPubKey(from)
							if clientErr != nil {
								return nil, clientErr
							}

							// retrieve the affiliate:
							var aff affiliates.Affiliate
							if client.HasBeenReferred() {
								aff, _ = dep.affiliateRepository.RetrieveByWallet(client.Referral())
							}

							vals, valsErr := dep.validatorRepository.RetrieveSetOrderedByPledgeAmount(0, gen.Info().MaxAmountOfValidators())
							if valsErr != nil {
								return nil, valsErr
							}

							// create the fees:
							fee := fees.SDKFunc.Create(fees.CreateParams{
								Gen:        gen,
								StoredData: jsData,
								Client:     client,
								Affiliate:  aff,
								Validators: vals,
							})

							// save the fees:
							saveFeesErr := dep.entityService.Save(fee, fees.SDKFunc.CreateRepresentation())
							if saveFeesErr != nil {
								return nil, saveFeesErr
							}

							// return the response:
							resp := routers.SDKFunc.CreateTransactionResponse(routers.CreateTransactionResponseParams{
								Code:    routers.IsSuccessful,
								Log:     "success",
								GazUsed: int64(fee.Client().Amount()),
								Tags: map[string][]byte{
									path: jsData,
								},
							})

							return resp, nil
						}

						str := fmt.Sprintf("the keyname (name: %s) cannot be voted on by group (name: %s)", keynameName, req.Request().Keyname().Group().Name())
						return nil, errors.New(str)
					}

					return nil, errors.New("an keyname must be provided")
				}

				str := fmt.Sprintf("the group (name: %s) is not an entity that can be voted on", req.Request().Keyname().Group().Name())
				return nil, errors.New(str)
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Request instance", reqIns.ID().String())
			return nil, errors.New(str)
		},
	}
}
