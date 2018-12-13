package commands

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/xmnservices/xmnsuite/crypto"
)

func createNodePK(nodePkAsString string) (tcrypto.PrivKey, error) {
	if nodePkAsString == "" {
		log.Printf("there was no nodePrivateKey provided. generating a new one...")
		nodePK := ed25519.GenPrivKey()
		return nodePK, nil
	}

	nodePKAsBytes, nodePKAsBytesErr := hex.DecodeString(nodePkAsString)
	if nodePKAsBytesErr != nil {
		// log:
		log.Printf("there was an error while decoding a string to hex: %s", nodePKAsBytesErr.Error())

		// output error:
		str := fmt.Sprintf("the given nodepk (%s) is not a valid private key", nodePkAsString)
		return nil, errors.New(str)
	}

	nodePK := new(ed25519.PrivKeyEd25519)
	nodePKErr := cdc.UnmarshalBinaryBare(nodePKAsBytes, nodePK)
	if nodePKErr != nil {
		// log:
		log.Printf("there was an error while Unmarshalling []byte to PrivateKey instance: %s", nodePKErr.Error())

		// output error:
		str := fmt.Sprintf("the given nodepk (%s) is not a valid private key", nodePkAsString)
		return nil, errors.New(str)
	}

	return nodePK, nil
}

func createRootPrivateKey(rootPrivateKeyAsString string) (crypto.PrivateKey, error) {
	if rootPrivateKeyAsString == "" {
		log.Printf("there was no rootPublicKey provided. generating a new root PrivateKey...")
		rootPrivKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{})
		return rootPrivKey, nil
	}

	privKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
		PKAsString: rootPrivateKeyAsString,
	})

	return privKey, nil
}

func createConstantFromParams(params CreateConstantsParams) Constants {
	id, idErr := uuid.FromString(params.ID)
	if idErr != nil {
		panic(idErr)
	}

	cons, consErr := createConstants(params.Namespace, params.Name, &id, params.RoutePrefix, params.RouterRoleKey)
	if consErr != nil {
		panic(consErr)
	}

	return cons
}

func createConfigsFromParams(params CreateParams) Configs {
	// create the constants:
	cons := createConstantFromParams(params.Constants)

	nodePK, nodePKErr := createNodePK(params.NodePrivateKey)
	if nodePKErr != nil {
		panic(nodePKErr)
	}

	conf, confErr := createConfigs(params.Meta, cons, params.Port, nodePK, params.BlockchainRootDirectory, params.DatabaseFilePath)
	if confErr != nil {
		panic(confErr)
	}

	return conf
}

func write(str string) string {
	out := fmt.Sprintf("\n************ xdns ************\n")
	out = fmt.Sprintf("%s%s", out, str)
	out = fmt.Sprintf("%s\n********** end xdns **********\n", out)
	return out
}

func print(str string) {
	fmt.Printf("%s", write(str))
}
