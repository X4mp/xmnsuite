package meta

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/affiliates"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/pledge"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/validator"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/link"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/node"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/withdrawal"
)

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
}

// CreateParams represents the create params
type CreateParams struct {
	AdditionalRead                 map[string]entity.MetaData
	AdditionalWrite                map[string]entity.Representation
	AdditionalWriteOnEntityRequest []CreateEntityRequestParams
}

// CreateEntityRequestParams represents the create entity request params
type CreateEntityRequestParams struct {
	RequestedBy entity.Representation
	Map         map[string]entity.Representation
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
				)
			}
		}

		// create the representations:
		genesisRepresentation := genesis.SDKFunc.CreateRepresentation()
		informationRepresentation := information.SDKFunc.CreateRepresentation()
		walletRepresentation := wallet.SDKFunc.CreateRepresentation()
		validatorRepresentation := validator.SDKFunc.CreateRepresentation()
		pledgeRepresentation := pledge.SDKFunc.CreateRepresentation()
		affiliatesRepresentation := affiliates.SDKFunc.CreateRepresentation()
		transferRepresentation := transfer.SDKFunc.CreateRepresentation()
		userRepresentation := user.SDKFunc.CreateRepresentation()
		linkRepresentation := link.SDKFunc.CreateRepresentation()
		nodeRepresentation := node.SDKFunc.CreateRepresentation()
		activeRequestRepresentation := active_request.SDKFunc.CreateRepresentation()
		keynameRepresentation := keyname.SDKFunc.CreateRepresentation()
		groupRepresentation := group.SDKFunc.CreateRepresentation()
		activeVoteRepresentation := active_vote.SDKFunc.CreateRepresentation()
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
			genesisRepresentation.MetaData().Keyname():       genesisRepresentation.MetaData(),
			informationRepresentation.MetaData().Keyname():   informationRepresentation.MetaData(),
			walletRepresentation.MetaData().Keyname():        walletRepresentation.MetaData(),
			validatorRepresentation.MetaData().Keyname():     validatorRepresentation.MetaData(),
			userRepresentation.MetaData().Keyname():          userRepresentation.MetaData(),
			activeRequestRepresentation.MetaData().Keyname(): activeRequestRepresentation.MetaData(),
			activeVoteRepresentation.MetaData().Keyname():    activeVoteRepresentation.MetaData(),
			pledgeRepresentation.MetaData().Keyname():        pledgeRepresentation.MetaData(),
			transferRepresentation.MetaData().Keyname():      transferRepresentation.MetaData(),
			linkRepresentation.MetaData().Keyname():          linkRepresentation.MetaData(),
			nodeRepresentation.MetaData().Keyname():          nodeRepresentation.MetaData(),
			withdrawalRepresentation.MetaData().Keyname():    withdrawalRepresentation.MetaData(),
			depositRepresentation.MetaData().Keyname():       depositRepresentation.MetaData(),
			keynameRepresentation.MetaData().Keyname():       keynameRepresentation.MetaData(),
			groupRepresentation.MetaData().Keyname():         groupRepresentation.MetaData(),
			affiliatesRepresentation.MetaData().Keyname():    affiliatesRepresentation.MetaData(),
		}

		// add the additional reads to the map:
		for keyname, oneAdditionalRead := range additionalReads {
			if _, ok := read[keyname]; ok {
				str := fmt.Sprintf("the keyname (%s) in the 'read' metaData is reserved for the core engine", keyname)
				panic(errors.New(str))
			}

			read[keyname] = oneAdditionalRead
		}

		// create the additional writes for wallets:
		additionalWriteForWallet := createEntityRequest(walletRepresentation, map[string]entity.Representation{
			pledgeRepresentation.MetaData().Keyname():     pledgeRepresentation,
			transferRepresentation.MetaData().Keyname():   transferRepresentation,
			userRepresentation.MetaData().Keyname():       userRepresentation,
			validatorRepresentation.MetaData().Keyname():  validatorRepresentation,
			walletRepresentation.MetaData().Keyname():     walletRepresentation, // for updates
			affiliatesRepresentation.MetaData().Keyname(): affiliatesRepresentation,
		})

		// create the additional writes for tokens:
		tokenRepresentation := token.SDKFunc.CreateRepresentation()
		additionalWriteForToken := createEntityRequest(tokenRepresentation, map[string]entity.Representation{
			informationRepresentation.MetaData().Keyname(): informationRepresentation,
			linkRepresentation.MetaData().Keyname():        linkRepresentation,
			nodeRepresentation.MetaData().Keyname():        nodeRepresentation,
		})

		// verify the additional writes for wallet:
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
			activeRequestRepresentation,
			activeVoteRepresentation,
			read,
			write,
			writeOnEntityRequest,
		)

		return out

	},
}
