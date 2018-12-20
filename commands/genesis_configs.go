package commands

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/crypto"
)

type genesisConfigs struct {
	conf               Configs
	rootPrivKey        crypto.PrivateKey
	genesisTransaction genesis.Genesis
}

func createGenesisConfigs(conf Configs, rootPrivKey crypto.PrivateKey, genesisTransaction genesis.Genesis) (GenesisConfigs, error) {
	out := genesisConfigs{
		conf:               conf,
		rootPrivKey:        rootPrivKey,
		genesisTransaction: genesisTransaction,
	}

	return &out, nil
}

// Configs returns the configs
func (obj *genesisConfigs) Configs() Configs {
	return obj.conf
}

// GenesisTransaction returns the genesis transaction
func (obj *genesisConfigs) GenesisTransaction() genesis.Genesis {
	return obj.genesisTransaction
}

// RootPrivateKey returns the root private key, if any
func (obj *genesisConfigs) RootPrivateKey() crypto.PrivateKey {
	return obj.rootPrivKey
}
