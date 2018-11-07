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
	"github.com/xmnservices/xmnsuite/blockchains/core/user"
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
	metaDatas         map[string]entity.MetaData
	representations   map[string]entity.Representation
	entityService     entity.Service
	entityRepository  entity.Repository
	genesisRepository genesis.Repository
	userRepository    user.Repository
	voteService       vote.Service
}

func createCore20181108(store datastore.DataStore) *core20181108 {
	out := core20181108{
		metaDatas: map[string]entity.MetaData{
			"wallet": wallet.SDKFunc.CreateMetaData(),
		},
		representations: map[string]entity.Representation{
			"wallet": wallet.SDKFunc.CreateRepresentation(),
		},
		entityService: entity.SDKFunc.CreateService(entity.CreateServiceParams{
			Store: store,
		}),
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
	routerDS datastore.DataStore,
	routerRoleKey string,
	ds datastore.DataStore,
) applications.Application {

	// create core:
	core := createCore20181108(ds)

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
		RouterParams: routers.CreateRouterParams{
			DataStore: routerDS,
			RoleKey:   routerRoleKey,
			RtesParams: []routers.CreateRouteParams{
				core.saveEntity(),
				core.retrieveEntityByID(),
				core.saveRequest(),
			},
		},
	})

	return app
}

func (app *core20181108) saveEntity() routers.CreateRouteParams {
	return routers.CreateRouteParams{
		Pattern: "/<name|[a-z-]+>",
		SaveTrx: func(store datastore.DataStore, from crypto.PublicKey, path string, params map[string]string, data []byte, sig crypto.Signature) (routers.TransactionResponse, error) {

			// retrieve the genesis:
			gen, genErr := app.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Gensis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the name:
			if name, ok := params["name"]; ok {
				// retrieve the entity representation:
				if representation, ok := app.representations[name]; ok {
					// unmarshal the data:
					ptr := representation.MetaData().CopyStorable()
					jsErr := cdc.UnmarshalJSON(data, ptr)
					if jsErr != nil {
						return nil, jsErr
					}

					// save the entity:
					saveErr := app.entityService.Save(ptr.(entity.Entity), representation)
					if saveErr != nil {
						return nil, saveErr
					}

					// convert to json:
					jsData, jsDataErr := cdc.MarshalJSON(ptr)
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
			// retrieve the name:
			if name, ok := params["name"]; ok {
				// retrieve the entity metadata:
				if metadata, ok := app.metaDatas[name]; ok {
					// parse the id:
					id, idErr := uuid.FromString(params["id"])
					if idErr != nil {
						return nil, idErr
					}

					// retrieve the entity instance:
					retIns, retInsErr := app.entityRepository.RetrieveByID(metadata, &id)
					if retInsErr != nil {
						return nil, retInsErr
					}

					// convert the entity to json:
					js, jsErr := cdc.MarshalJSON(retIns)
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

			// retrieve the genesis:
			gen, genErr := app.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Gensis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the from user:
			fromUser, fromUserErr := app.userRepository.RetrieveByPubKey(from)
			if fromUserErr != nil {
				str := fmt.Sprintf("the from user (pubKey: %s) could not be found", from.String())
				return nil, errors.New(str)
			}

			// retrieve the name:
			if name, ok := params["name"]; ok {
				// retrieve the entity representation:
				if representation, ok := app.representations[name]; ok {
					if requestIDAsString, ok := params["id"]; ok {
						requestID, requestIDErr := uuid.FromString(requestIDAsString)
						if requestIDErr != nil {
							return nil, requestIDErr
						}

						// unmarshal the data:
						ptr := representation.MetaData().CopyStorable()
						jsErr := cdc.UnmarshalJSON(data, ptr)
						if jsErr != nil {
							return nil, jsErr
						}

						// build the request:
						req := request.SDKFunc.Create(request.CreateParams{
							ID:        &requestID,
							FromUser:  fromUser,
							NewEntity: ptr.(entity.Entity),
						})

						// save the request:
						saveErr := app.entityService.Save(req, request.SDKFunc.CreateRepresentation(request.CreateRepresentationParams{
							Met: representation.MetaData(),
						}))

						if saveErr != nil {
							return nil, saveErr
						}

						// convert to json:
						jsData, jsDataErr := cdc.MarshalJSON(req)
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

			// retrieve the genesis:
			gen, genErr := app.genesisRepository.Retrieve()
			if genErr != nil {
				str := fmt.Sprintf("there was an error while retrieving the Gensis instance: %s", genErr.Error())
				return nil, errors.New(str)
			}

			// retrieve the from user:
			fromUser, fromUserErr := app.userRepository.RetrieveByPubKey(from)
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
				if representation, ok := app.representations[name]; ok {
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
							req, reqErr := app.entityRepository.RetrieveByID(request.SDKFunc.CreateMetaData(request.CreateMetaDataParams{
								Met: representation.MetaData(),
							}), &requestID)

							if reqErr != nil {
								return nil, reqErr
							}

							// create the vote:
							vote := vote.SDKFunc.Create(vote.CreateParams{
								ID:         &voteID,
								Request:    req.(request.Request),
								Voter:      fromUser,
								IsApproved: ptr.IsApproved,
							})

							// save the vote:
							saveErr := app.voteService.Save(vote)
							if saveErr != nil {
								return nil, saveErr
							}

							// convert to json:
							jsData, jsDataErr := cdc.MarshalJSON(vote)
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
