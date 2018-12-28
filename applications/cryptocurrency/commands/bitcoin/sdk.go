package bitcoin

import (
	"math"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
)

const (
	namespace                    = "cryptocurrency"
	name                         = "bitcoin"
	id                           = "ac6da9be-c817-4527-9223-44cb6feabe10"
	databaseFilePath             = "db/blockchain/blockchain.db"
	blockchainRootDirectory      = "db/blockchain/files"
	routePrefix                  = ""
	tokenSymbol                  = "XBR"
	tokenName                    = "XMN Bitcoin Representation"
	tokenDescription             = "The XMN Bitcoin Representation is a decentralized marketplace of users that accept bitcoin deposits and spawn a PoB blockchain where representations of these bitcoins are created.  Upon withdrawal, the PoB tokens are destroyed and the real bitcoins are transfered to its owner."
	totalTokenAmount             = math.MaxInt64 - 1
	initialWalletConcensus       = 50
	initialGazPricePerKB         = 1
	initialTokenConcensusNeeded  = 50
	initialMaxAmountOfValidators = 200
	initialUserAmountOfShares    = 100
)

var peers = []string{}

// SpawnParams represents the spawn params
type SpawnParams struct {
	Pass     string
	Filename string
	Dir      string
	Port     int
}

// SDKFunc represents the commands SDK func
var SDKFunc = struct {
	Spawn func(params SpawnParams) applications.Node
}{
	Spawn: func(params SpawnParams) applications.Node {
		out, outErr := spawn(params.Pass, params.Filename, params.Dir, params.Port)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
