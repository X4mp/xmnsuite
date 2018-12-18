package configs

import (
	"encoding/base64"

	tcrypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	crypto "github.com/xmnservices/xmnsuite/crypto"
)

// Configs represents the configs
type Configs interface {
	NodePK() tcrypto.PrivKey
	WalletPK() crypto.PrivateKey
	String() string
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

// CreateParams represents the create params
type CreateParams struct {
	Encoded string
}

// EncryptParams represents the encrypt params
type EncryptParams struct {
	Conf        Configs
	Pass        string
	RetypedPass string
}

// DecryptParams represents the decrypt params
type DecryptParams struct {
	Data string
	Pass string
}

// SDKFunc represents the confgis sdk func
var SDKFunc = struct {
	Generate         func() Configs
	Create           func(params CreateParams) Configs
	Encrypt          func(params EncryptParams) string
	Decrypt          func(params DecryptParams) Configs
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
	Create: func(params CreateParams) Configs {
		decoded, decodedErr := base64.StdEncoding.DecodeString(params.Encoded)
		if decodedErr != nil {
			panic(decodedErr)
		}

		ptr := new(storableConfigs)
		jsErr := cdc.UnmarshalJSON(decoded, ptr)
		if jsErr != nil {
			panic(jsErr)
		}

		conf, confErr := fromStorableToConfigs(ptr)
		if confErr != nil {
			panic(confErr)
		}

		return conf
	},
	Encrypt: func(params EncryptParams) string {
		encrypted, encryptedErr := encrypt(params.Conf, params.Pass, params.RetypedPass)
		if encryptedErr != nil {
			panic(encryptedErr)
		}

		return encrypted
	},
	Decrypt: func(params DecryptParams) Configs {
		conf, confErr := decrypt(params.Data, params.Pass)
		if confErr != nil {
			panic(confErr)
		}

		return conf
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
