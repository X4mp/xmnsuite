package router

import (
	"errors"
	"fmt"

	crypto "github.com/tendermint/tendermint/crypto"
)

/*
 * Request
 */

type requestSignedStruct struct {
	Path string `json:"path"`
	Data []byte `json:"data"`
}

type request struct {
	Frm crypto.PubKey `json:"from"`
	Pth string        `json:"path"`
	Dat []byte        `json:"data"`
	Sig []byte        `json:"signature"`
}

func createRequest(from crypto.PubKey, path string, data []byte, sig []byte) (Request, error) {

	if from == nil {
		return nil, errors.New("the from pubKey is mandatory")
	}

	if path == "" {
		return nil, errors.New("the path is mandatory")
	}

	if data == nil {
		return nil, errors.New("the data is mandatory")
	}

	if sig == nil {
		return nil, errors.New("the signature is mandatory")
	}

	str := requestSignedStruct{
		Path: path,
		Data: data,
	}

	js, jsErr := cdc.MarshalJSON(str)
	if jsErr != nil {
		return nil, jsErr
	}

	if !from.VerifyBytes(js, sig) {
		str := fmt.Sprintf("the path, data and signature could not be verified by the given public key")
		return nil, errors.New(str)
	}

	out := request{
		Frm: from,
		Pth: path,
		Dat: data,
		Sig: sig,
	}

	return &out, nil
}

// Path returns the path
func (obj *request) From() crypto.PubKey {
	return obj.Frm
}

// Path returns the path
func (obj *request) Path() string {
	return obj.Pth
}

// Data returns the data
func (obj *request) Data() []byte {
	return obj.Dat
}

// Signature returns the signature
func (obj *request) Signature() []byte {
	return obj.Sig
}
