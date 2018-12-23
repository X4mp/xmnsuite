package meta

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/validator"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
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
		keynameRepresentation := keyname.SDKFunc.CreateRepresentation()
		groupRepresentation := group.SDKFunc.CreateRepresentation()
		voteRepresentation := vote.SDKFunc.CreateRepresentation()
		withdrawalRepresentation := withdrawal.SDKFunc.CreateRepresentation()
		depositRepresentation := deposit.SDKFunc.CreateRepresentation()

		// create the additional writes:
		additionalWrites := map[string]entity.Representation{
			keynameRepresentation.MetaData().Keyname(): keynameRepresentation,
			groupRepresentation.MetaData().Keyname():   groupRepresentation,
		}

		// add the additional writes to the map:
		for keyname, oneAdditionalWrite := range additionalWrites {
			if _, ok := write[keyname]; ok {
				str := fmt.Sprintf("the keyname (%s) in the 'write' metaData is reserved for the core engine", keyname)
				panic(errors.New(str))
			}

			write[keyname] = oneAdditionalWrite
		}

		// create the read:
		additionalReads := map[string]entity.MetaData{
			genesisRepresentation.MetaData().Keyname():    genesisRepresentation.MetaData(),
			walletRepresentation.MetaData().Keyname():     walletRepresentation.MetaData(),
			validatorRepresentation.MetaData().Keyname():  validatorRepresentation.MetaData(),
			userRepresentation.MetaData().Keyname():       userRepresentation.MetaData(),
			requestRepresentation.MetaData().Keyname():    requestRepresentation.MetaData(),
			voteRepresentation.MetaData().Keyname():       voteRepresentation.MetaData(),
			pledgeRepresentation.MetaData().Keyname():     pledgeRepresentation.MetaData(),
			transferRepresentation.MetaData().Keyname():   transferRepresentation.MetaData(),
			linkRepresentation.MetaData().Keyname():       linkRepresentation.MetaData(),
			nodeRepresentation.MetaData().Keyname():       nodeRepresentation.MetaData(),
			withdrawalRepresentation.MetaData().Keyname(): withdrawalRepresentation.MetaData(),
			depositRepresentation.MetaData().Keyname():    depositRepresentation.MetaData(),
			keynameRepresentation.MetaData().Keyname():    keynameRepresentation.MetaData(),
			groupRepresentation.MetaData().Keyname():      groupRepresentation.MetaData(),
		}

		// add the additional reads to the map:
		for keyname, oneAdditionalRead := range additionalReads {
			if _, ok := read[keyname]; ok {
				str := fmt.Sprintf("the keyname (%s) in the 'read' metaData is reserved for the core engine", keyname)
				panic(errors.New(str))
			}

			read[keyname] = oneAdditionalRead
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
								requesterWalletID := oneVote.Request().From().Wallet().ID()
								voterWalletID := oneVote.Voter().Wallet().ID()
								if bytes.Compare(requesterWalletID.Bytes(), voterWalletID.Bytes()) != 0 {
									str := fmt.Sprintf("the requester is binded to a wallet (ID: %s) that is different from the voter's wallet (ID: %s), therefore the wallet vote is not counted.  Skipping...", requesterWalletID.String(), voterWalletID.String())
									log.Printf(str)
									continue
								}

								if oneVote.IsApproved() {
									approved += oneVote.Voter().Shares()
									continue
								}

								disapproved += oneVote.Voter().Shares()
								continue
							}

							log.Printf("the entity (ID: %s) is not a valid Vote instance", oneVoteIns.ID().String())
						}

						// if there is enough approved or disapproved votes, the concensus passed:
						concensusPassed := (approved >= neededConcensus) || (disapproved >= neededConcensus)

						// vote is approved:
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
			pledgeRepresentation.MetaData().Keyname():    pledgeRepresentation,
			transferRepresentation.MetaData().Keyname():  transferRepresentation,
			userRepresentation.MetaData().Keyname():      userRepresentation,
			validatorRepresentation.MetaData().Keyname(): validatorRepresentation,
			walletRepresentation.MetaData().Keyname():    walletRepresentation, // for updates
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
			linkRepresentation.MetaData().Keyname(): linkRepresentation,
			nodeRepresentation.MetaData().Keyname(): nodeRepresentation,
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
