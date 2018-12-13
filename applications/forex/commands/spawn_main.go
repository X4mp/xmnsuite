package commands

import (
	"path/filepath"

	"github.com/xmnservices/xmnsuite/applications/forex/meta"
	"github.com/xmnservices/xmnsuite/blockchains"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token"
	"github.com/xmnservices/xmnsuite/configs"
)

func spawnMain(pass string, filename string, rootDir string, port int) (applications.Node, error) {
	// create the repository:
	repository := configs.SDKFunc.CreateRepository()

	// retrieve the configs:
	retConf, retConfErr := repository.Retrieve(filename, pass)
	if retConfErr != nil {
		return nil, retConfErr
	}

	wal := wallet.SDKFunc.Create(wallet.CreateParams{
		Creator:         retConf.WalletPK().PublicKey(),
		ConcensusNeeded: initialWalletConcensus,
	})

	blkChain := blockchains.SDKFunc.Create(blockchains.CreateParams{
		Port:      port,
		Name:      name,
		Namespace: namespace,
		ID:        id,
		Conf:      retConf,
		BlockchainRootDirectory: filepath.Join(rootDir, blockchainRootDirectory),
		DatabaseFilePath:        filepath.Join(rootDir, databaseFilePath),
		Peers:                   peers,
		Meta:                    meta.SDKFunc.CreateMetaData(),
		GenesisTransaction: genesis.SDKFunc.Create(genesis.CreateParams{
			GazPricePerKb:         initialGazPricePerKB,
			ConcensusNeeded:       initialTokenConcensusNeeded,
			MaxAmountOfValidators: initialMaxAmountOfValidators,
			User: user.SDKFunc.Create(user.CreateParams{
				PubKey: retConf.WalletPK().PublicKey(),
				Shares: initialUserAmountOfShares,
				Wallet: wal,
			}),
			Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
				To: wal,
				Token: token.SDKFunc.Create(token.CreateParams{
					Symbol:      tokenSymbol,
					Name:        tokenName,
					Description: tokenDescription,
				}),
				Amount: totalTokenAmount,
			}),
		}),
	})

	// start the blockchain:
	node, nodeErr := blkChain.Start()
	if nodeErr != nil {
		return nil, nodeErr
	}

	return node, nil
}
