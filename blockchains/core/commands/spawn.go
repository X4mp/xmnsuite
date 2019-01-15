package commands

import (
	"path/filepath"

	"github.com/xmnservices/xmnsuite/blockchains"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/deposit"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/information"
	"github.com/xmnservices/xmnsuite/configs"
)

func spawn(pass string, filename string, rootDir string, port int) (applications.Node, error) {
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
		Meta:                    meta.SDKFunc.Create(meta.CreateParams{}),
		GenesisTransaction: genesis.SDKFunc.Create(genesis.CreateParams{
			Info: information.SDKFunc.Create(information.CreateParams{
				GazPricePerKb:         initialGazPricePerKB,
				ConcensusNeeded:       initialTokenConcensusNeeded,
				MaxAmountOfValidators: initialMaxAmountOfValidators,
				NetworkShare:          initialNetworkShare,
				ValidatorsShare:       initialValidatorShare,
				AffiliateShare:        initialReferralShare,
			}),
			User: user.SDKFunc.Create(user.CreateParams{
				PubKey: retConf.WalletPK().PublicKey(),
				Shares: initialUserAmountOfShares,
				Wallet: wal,
			}),
			Deposit: deposit.SDKFunc.Create(deposit.CreateParams{
				To:     wal,
				Amount: totalTokenAmount,
			}),
			Token: token.SDKFunc.Create(token.CreateParams{
				Symbol:      tokenSymbol,
				Name:        tokenName,
				Description: tokenDescription,
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
