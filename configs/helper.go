package configs

import (
	"encoding/hex"
	"errors"
	"fmt"

	tcrypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
)

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
