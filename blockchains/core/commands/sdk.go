package commands

import (
	"math"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/configs"
)

const (
	namespace                    = "xmn"
	name                         = "core"
	id                           = "5aac4a8d-7cb5-457c-b43d-ca0933c20dab"
	databaseFilePath             = "db/blockchain/blockchain.db"
	blockchainRootDirectory      = "db/blockchain/files"
	tokenSymbol                  = "XMN"
	tokenName                    = "XMN"
	tokenDescription             = "The XMN core is a decentralized core blockchain infrastructure."
	totalTokenAmount             = math.MaxInt64 - 1
	initialWalletConcensus       = 50
	initialGazPricePerKB         = 1
	initialTokenConcensusNeeded  = 50
	initialMaxAmountOfValidators = 200
	initialNetworkShare          = 5
	initialValidatorShare        = 80
	initialReferralShare         = 15
	initialUserAmountOfShares    = 100
)

var peers = []string{}

// GenerateConfigsParams represents the generate configs params
type GenerateConfigsParams struct {
	Pass        string
	RetypedPass string
	Filename    string
}

// SpawnParams represents the spawn params
type SpawnParams struct {
	Pass     string
	Filename string
	Dir      string
	Port     int
}

// SDKFunc represents the commands SDK func
var SDKFunc = struct {
	GenerateConfigs func(params GenerateConfigsParams) configs.Configs
	Spawn           func(params SpawnParams) applications.Node
}{
	GenerateConfigs: func(params GenerateConfigsParams) configs.Configs {
		out, outErr := generateConfigs(params.Pass, params.RetypedPass, params.Filename)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	Spawn: func(params SpawnParams) applications.Node {
		out, outErr := spawn(params.Pass, params.Filename, params.Dir, params.Port)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
