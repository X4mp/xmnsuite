package router

import (
	crypto "github.com/tendermint/tendermint/crypto"
)

type trxChkRequest struct {
	from            crypto.PubKey
	path            string
	dataSizeInBytes int64
	sig             []byte
}

func createTrxChkRequest(from crypto.PubKey, path string, dataSizeInBytes int64, sig []byte) TrxChkRequest {
	out := trxChkRequest{
		from:            from,
		path:            path,
		dataSizeInBytes: dataSizeInBytes,
		sig:             sig,
	}

	return &out
}

// Path returns the path
func (obj *trxChkRequest) From() crypto.PubKey {
	return obj.from
}

// Path returns the path
func (obj *trxChkRequest) Path() string {
	return obj.path
}

// Data returns the data
func (obj *trxChkRequest) DataSizeInBytes() int64 {
	return obj.dataSizeInBytes
}

// Signature returns the signature
func (obj *trxChkRequest) Signature() []byte {
	return obj.sig
}
