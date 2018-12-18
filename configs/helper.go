package configs

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	tcrypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/crypto"
)

func encrypt(conf Configs, pass string, retypedPass string) (string, error) {
	if len(pass) < 6 {
		return "", errors.New("The password must contain at least 6 characters")
	}

	if pass != retypedPass {
		return "", errors.New("The passwords do not match")
	}

	storable := createStorableConfigs(conf)
	js, jsErr := json.Marshal(storable)
	if jsErr != nil {
		return "", jsErr
	}

	encrypted := crypto.SDKFunc.Encrypt(crypto.EncryptParams{
		Pass: []byte(pass),
		Msg:  js,
	})

	return encrypted, nil
}

func decrypt(encryptedData string, pass string) (Configs, error) {
	decrypted := crypto.SDKFunc.Decrypt(crypto.DecryptParams{
		Pass:         []byte(pass),
		EncryptedMsg: encryptedData,
	})

	ptr := new(storableConfigs)
	jsErr := json.Unmarshal(decrypted, ptr)
	if jsErr != nil {
		return nil, jsErr
	}

	configs, configsErr := fromStorableToConfigs(ptr)
	if configsErr != nil {
		return nil, configsErr
	}

	return configs, nil
}

func fromEncodedStringToPrivKey(str string) (tcrypto.PrivKey, error) {
	privKeyAsBytes, privKeyAsBytesErr := hex.DecodeString(str)
	if privKeyAsBytesErr != nil {
		return nil, privKeyAsBytesErr
	}

	privKey := new(ed25519.PrivKeyEd25519)
	privKeyErr := cdc.UnmarshalBinaryBare(privKeyAsBytes, privKey)
	if privKeyErr != nil {
		str := fmt.Sprintf("the private key []byte is invalid: %s", privKeyErr.Error())
		return nil, errors.New(str)
	}

	return privKey, nil
}
