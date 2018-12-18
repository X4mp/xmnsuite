package meta

import (
	"errors"
	"fmt"
	"log"

	"github.com/xmnservices/xmnsuite/applications/forex/objects/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
	"github.com/xmnservices/xmnsuite/datastore"
)

// CreateVoteServiceFn represents a func to create a VoteService
type CreateVoteServiceFn func(ds datastore.DataStore) vote.Service

// Meta represents the meta data for the core application
type Meta interface {
	Genesis() entity.Representation
	Wallet() entity.Representation
	Request() entity.Representation
	Vote() entity.Representation
	Retrieval() map[string]entity.MetaData
	Write() map[string]entity.Representation
	WriteOnAllEntityRequest() map[string]entity.Representation
	WriteOnEntityRequest() map[string]EntityRequest
	AddToWriteOnEntityRequest(requestedBy entity.MetaData, rep entity.Representation) error
}

// EntityRequest represents a save on entity request meta data
type EntityRequest interface {
	RequestedBy() entity.Representation
	Map() map[string]entity.Representation
	Add(rep entity.Representation) EntityRequest
	VoteService(store datastore.DataStore) vote.Service
}

// CreateParams represents the create params
type CreateParams struct {
	AdditionalRead                 map[string]entity.MetaData
	AdditionalWrite                map[string]entity.Representation
	AdditionalWriteOnEntityRequest []CreateEntityRequestParams
}

// CreateEntityRequestParams represents the create entity request params
type CreateEntityRequestParams struct {
	RequestedBy         entity.Representation
	Map                 map[string]entity.Representation
	CreateVoteServiceFn CreateVoteServiceFn
}

// SDKFunc represents the meta SDK func
var SDKFunc = struct {
	Create func(params CreateParams) Meta
}{
	Create: func(params CreateParams) Meta {
		read := map[string]entity.MetaData{}
		write := map[string]entity.Representation{}
		writeOnEntityRequest := map[string]EntityRequest{}
		if params.AdditionalWrite != nil {
			write = params.AdditionalWrite
		}

		if params.AdditionalRead != nil {
			read = params.AdditionalRead
		}

		if params.AdditionalWriteOnEntityRequest != nil {
			for _, oneParams := range params.AdditionalWriteOnEntityRequest {
				writeOnEntityRequest[oneParams.RequestedBy.MetaData().Keyname()] = createEntityRequest(
					oneParams.RequestedBy,
					oneParams.Map,
					oneParams.CreateVoteServiceFn,
				)
			}
		}

		// create the representations:
		genesisRepresentation := genesis.SDKFunc.CreateRepresentation()
		walletRepresentation := wallet.SDKFunc.CreateRepresentation()
		validatorRepresentation := validator.SDKFunc.CreateRepresentation()
		pledgeRepresentation := pledge.SDKFunc.CreateRepresentation()
		transferRepresentation := transfer.SDKFunc.CreateRepresentation()
		userRepresentation := user.SDKFunc.CreateRepresentation()
		linkRepresentation := link.SDKFunc.CreateRepresentation()
		nodeRepresentation := node.SDKFunc.CreateRepresentation()
		requestRepresentation := request.SDKFunc.CreateRepresentation()
		voteRepresentation := vote.SDKFunc.CreateRepresentation()
		withdrawalRepresentation := withdrawal.SDKFunc.CreateRepresentation()
		depositRepresentation := deposit.SDKFunc.CreateRepresentation()

		// create the read:
		additionalReads := map[string]entity.MetaData{
			"genesis":    genesisRepresentation.MetaData(),
			"wallet":     walletRepresentation.MetaData(),
			"validator":  validatorRepresentation.MetaData(),
			"user":       userRepresentation.MetaData(),
			"request":    requestRepresentation.MetaData(),
			"vote":       voteRepresentation.MetaData(),
			"pledge":     pledgeRepresentation.MetaData(),
			"transfer":   transferRepresentation.MetaData(),
			"link":       linkRepresentation.MetaData(),
			"node":       nodeRepresentation.MetaData(),
			"withdrawal": withdrawalRepresentation.MetaData(),
			"deposit":    depositRepresentation.MetaData(),
		}

		// add the additional reads to the map:
		for keyname, oneAdditionalRead := range additionalReads {
			if _, ok := read[keyname]; ok {
				str := fmt.Sprintf("the keyname (%s) in the 'read' metaData is reserved for the core engine", keyname)
				panic(errors.New(str))
			}

			read[keyname] = oneAdditionalRead
		}

		// create the write:
		additionalWrites := map[string]entity.Representation{
			"wallet": walletRepresentation,
		}

		// add the additional writes to the map:
		for keyname, oneAdditionalWrite := range additionalWrites {
			if _, ok := write[keyname]; ok {
				str := fmt.Sprintf("the keyname (%s) in the 'write' representations is reserved for the core engine", keyname)
				panic(errors.New(str))
			}

			write[keyname] = oneAdditionalWrite
		}

		// create the wallet vote service:
		walletVoteService := func(store datastore.DataStore) vote.Service {
			return vote.SDKFunc.CreateService(vote.CreateServiceParams{
				CalculateVoteFn: func(votes entity.PartialSet) (bool, bool, error) {
					if votes.Amount() <= 0 {
						return false, false, errors.New("the votes cannot be empty")
					}

					// retrieve the needed concensus from the requester wallet:
					votesIns := votes.Instances()
					if firstVote, ok := votesIns[0].(vote.Vote); ok {
						// check the amount of concensus needed:
						neededConcensus := firstVote.Request().From().Wallet().ConcensusNeeded()

						// compile the vote's concensus:
						approved := 0
						disapproved := 0
						for _, oneVoteIns := range votesIns {
							if oneVote, ok := oneVoteIns.(vote.Vote); ok {
								if oneVote.IsApproved() {
									approved += oneVote.Voter().Shares()
									continue
								}

								disapproved += oneVote.Voter().Shares()
								continue
							}

							log.Printf("the entity (ID: %s) is not a valid Vote instance", oneVoteIns.ID().String())
						}

						// if there is enugh approved or disapproved votes, the concensus passed:
						concensusPassed := (approved >= neededConcensus) || (disapproved >= neededConcensus)

						// vote is approved, insert the new entity:
						if approved >= neededConcensus {
							return true, concensusPassed, nil
						}

						return false, concensusPassed, nil
					}

					return false, false, errors.New("the given entityPartialSet does not contain valid Vote instances")
				},
				DS: store,
			})
		}

		// create the additional writes for wallets:
		additionalWriteForWallet := createEntityRequest(walletRepresentation, map[string]entity.Representation{
			"pledge":    pledgeRepresentation,
			"transfer":  transferRepresentation,
			"user":      userRepresentation,
			"validator": validatorRepresentation,
			"wallet":    walletRepresentation, // for updates
		}, walletVoteService)

		// create the token vote service:
		tokenVoteService := func(store datastore.DataStore) vote.Service {
			return vote.SDKFunc.CreateService(vote.CreateServiceParams{
				CalculateVoteFn: func(votes entity.PartialSet) (bool, bool, error) {
					// retrieve the genesis:
					genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
						Datastore: store,
					})

					balanceRepository := balance.SDKFunc.CreateRepository(balance.CreateRepositoryParams{
						Datastore: store,
					})

					// retrieve genesis:
					gen, genErr := genesisRepository.Retrieve()
					if genErr != nil {
						str := fmt.Sprintf("there was an error while retrieving the Genesis instance: %s", genErr.Error())
						return false, false, errors.New(str)
					}

					// retrieve the token:
					tok := gen.Deposit().Token()

					// retrieve the needed concensus:
					neededConcensus := gen.ConcensusNeeded()

					// compile the vote's concensus:
					approved := 0
					disapproved := 0
					votesIns := votes.Instances()
					for _, oneVoteIns := range votesIns {
						if oneVote, ok := oneVoteIns.(vote.Vote); ok {
							// retrieve the balance:
							wal := oneVote.Voter().Wallet()
							balance, balanceErr := balanceRepository.RetrieveByWalletAndToken(wal, tok)
							if balanceErr != nil {
								str := fmt.Sprintf("there was an error while retrieving the balance on wallet (ID: %s) for token (ID: %s): %s", wal.ID().String(), tok.ID().String(), balanceErr.Error())
								return false, false, errors.New(str)
							}

							if oneVote.IsApproved() {
								approved += balance.Amount()
								continue
							}

							disapproved += balance.Amount()
							continue
						}

						log.Printf("the entity (ID: %s) is not a valid Vote instance", oneVoteIns.ID().String())
					}

					// if there is enugh approved or disapproved votes, the concensus passed:
					concensusPassed := (approved >= neededConcensus) || (disapproved >= neededConcensus)

					// vote is approved, insert the new entity:
					if approved >= neededConcensus {
						return true, concensusPassed, nil
					}

					return false, concensusPassed, nil
				},
				DS: store,
			})
		}

		// create the additional writes for tokens:
		tokenRepresentation := token.SDKFunc.CreateRepresentation()
		additionalWriteForToken := createEntityRequest(tokenRepresentation, map[string]entity.Representation{
			"link": linkRepresentation,
			"node": nodeRepresentation,
		}, tokenVoteService)

		// add the additional writes for wallet:
		walKeyname := additionalWriteForWallet.RequestedBy().MetaData().Keyname()
		tokKeyname := additionalWriteForToken.RequestedBy().MetaData().Keyname()
		for _, oneWriteOnEntityReq := range writeOnEntityRequest {
			keyname := oneWriteOnEntityReq.RequestedBy().MetaData().Keyname()
			if keyname == walKeyname || keyname == tokKeyname {
				str := fmt.Sprintf("the keyname (%s) in the 'write on entity request' representations is reserved for the core engine", keyname)
				panic(errors.New(str))
			}
		}

		// add the additional write on entity requests:
		writeOnEntityRequest[walKeyname] = additionalWriteForWallet
		writeOnEntityRequest[tokKeyname] = additionalWriteForToken

		// create the meta instance:
		out := createMeta(
			genesisRepresentation,
			walletRepresentation,
			requestRepresentation,
			voteRepresentation,
			read,
			write,
			writeOnEntityRequest,
		)

		return out

	},
}
