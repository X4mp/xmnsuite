package configs

import (
	tcrypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	crypto "github.com/xmnservices/xmnsuite/crypto"
)

// Configs represents the configs
type Configs interface {
	NodePK() tcrypto.PrivKey
	WalletPK() crypto.PrivateKey
}

// Normalized represents the normalized configs
type Normalized interface {
}

// Repository represents the configs repository
type Repository interface {
	Retrieve(filePath string, password string) (Configs, error)
}

// Service represents the configs service
type Service interface {
	Save(ins Configs, filePath string, password string, retypedPassword string) error
}

// SDKFunc represents the confgis sdk func
var SDKFunc = struct {
	Generate         func() Configs
	Normalize        func(ins Configs) Normalized
	CreateRepository func() Repository
	CreateService    func() Service
}{
	Generate: func() Configs {
		nodePK := ed25519.GenPrivKey()
		walletPK := crypto.SDKFunc.GenPK()
		out, outErr := createConfigs(nodePK, walletPK)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
	Normalize: func(ins Configs) Normalized {
		out := createStorableConfigs(ins)
		return out
	},
	CreateRepository: func() Repository {
		out := createRepository()
		return out
	},
	CreateService: func() Service {
		out := createService()
		return out
	},
}
