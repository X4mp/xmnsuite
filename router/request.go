package router

import (
	crypto "github.com/tendermint/tendermint/crypto"
)

type request struct {
	from crypto.PubKey
	path string
	data []byte
	sig  []byte
}

func createRequest(from crypto.PubKey, path string, data []byte, sig []byte) Request {
	out := request{
		from: from,
		path: path,
		data: data,
		sig:  sig,
	}

	return &out
}

// Path returns the path
func (obj *request) From() crypto.PubKey {
	return obj.from
}

// Path returns the path
func (obj *request) Path() string {
	return obj.path
}

// Data returns the data
func (obj *request) Data() []byte {
	return obj.data
}

// Signature returns the signature
func (obj *request) Signature() []byte {
	return obj.sig
}
