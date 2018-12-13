package commands

import (
	"math"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/configs"
)

const (
	namespace                    = "xmn"
	name                         = "forex"
	id                           = "63166ddf-9cdf-440f-b357-e502d696f1ff"
	databaseFilePath             = "db/blockchain/blockchain.db"
	blockchainRootDirectory      = "db/blockchain/files"
	routePrefix                  = ""
	tokenSymbol                  = "XMN"
	tokenName                    = "XMN Foreign Exchange"
	tokenDescription             = "The XMN foreign exchange is a blockchain that enables anyone to create a bank and pledge XMN tokens in order to manage physical currency deposits.  The currency types can range from physical gold, US dollars or a plumber's time."
	totalTokenAmount             = math.MaxInt64 - 1
	initialWalletConcensus       = 50
	initialGazPricePerKB         = 1
	initialTokenConcensusNeeded  = 50
	initialMaxAmountOfValidators = 200
	initialUserAmountOfShares    = 100
)

var peers = []string{}

// GenerateConfigsParams represents the generate configs params
type GenerateConfigsParams struct {
	Pass        string
	RetypedPass string
	Filename    string
}

// SpawnMainParams represents the spawn main params
type SpawnMainParams struct {
	Pass     string
	Filename string
	Dir      string
	Port     int
}

// SDKFunc represents the commands SDK func
var SDKFunc = struct {
	GenerateConfigs func(params GenerateConfigsParams) configs.Configs
	SpawnMain       func(params SpawnMainParams) applications.Node
}{
	GenerateConfigs: func(params GenerateConfigsParams) configs.Configs {
		out, outErr := generateConfigs(params.Pass, params.RetypedPass, params.Filename)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	SpawnMain: func(params SpawnMainParams) applications.Node {
		out, outErr := spawnMain(params.Pass, params.Filename, params.Dir, params.Port)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
