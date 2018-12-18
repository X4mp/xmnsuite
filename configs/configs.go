package configs

import (
	"encoding/base64"
	"errors"

	tcrypto "github.com/tendermint/tendermint/crypto"
	crypto "github.com/xmnservices/xmnsuite/crypto"
)

type configs struct {
	nodePK   tcrypto.PrivKey
	walletPK crypto.PrivateKey
}

func createConfigs(nodePK tcrypto.PrivKey, walletPK crypto.PrivateKey) (Configs, error) {
	if nodePK == nil {
		return nil, errors.New("the nodePK is mandatory")
	}

	if walletPK == nil {
		return nil, errors.New("the walletPK is mandatory")
	}

	out := configs{
		nodePK:   nodePK,
		walletPK: walletPK,
	}

	return &out, nil
}

func fromStorableToConfigs(storable *storableConfigs) (Configs, error) {
	nodePK, nodePKErr := fromEncodedStringToPrivKey(storable.NodePK)
	if nodePKErr != nil {
		return nil, nodePKErr
	}

	walletPK := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
		PKAsString: storable.WalletPK,
	})

	return createConfigs(nodePK, walletPK)
}

// NodePK returns the nodePK
func (obj *configs) NodePK() tcrypto.PrivKey {
	return obj.nodePK
}

// WalletPK returns the walletPK
func (obj *configs) WalletPK() crypto.PrivateKey {
	return obj.walletPK
}

// String returns the configs as string
func (obj *configs) String() string {
	storable := createStorableConfigs(obj)
	js, jsErr := cdc.MarshalJSON(storable)
	if jsErr != nil {
		panic(jsErr)
	}

	return base64.StdEncoding.EncodeToString(js)
}
