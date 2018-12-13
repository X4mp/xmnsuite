package commands

import (
	tcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
)

type configs struct {
	cons              Constants
	port              int
	nodePK            tcrypto.PrivKey
	blockchainRootDir string
	databaseFilePath  string
	met               meta.Meta
}

func createConfigs(met meta.Meta, cons Constants, port int, nodePK tcrypto.PrivKey, blockchainRootDir string, databaseFilePath string) (Configs, error) {
	out := configs{
		met:               met,
		cons:              cons,
		port:              port,
		nodePK:            nodePK,
		blockchainRootDir: blockchainRootDir,
		databaseFilePath:  databaseFilePath,
	}

	return &out, nil
}

// Meta returns the meta
func (obj *configs) Meta() meta.Meta {
	return obj.met
}

// Constants returns the constants
func (obj *configs) Constants() Constants {
	return obj.cons
}

// Port returns the port
func (obj *configs) Port() int {
	return obj.port
}

// NodePrivateKey returns the node private key
func (obj *configs) NodePrivateKey() tcrypto.PrivKey {
	return obj.nodePK
}

// BlockchainRootDirectory returns the blockchain root directory
func (obj *configs) BlockchainRootDirectory() string {
	return obj.blockchainRootDir
}

// DatabaseFilePath returns the database file path
func (obj *configs) DatabaseFilePath() string {
	return obj.databaseFilePath
}
