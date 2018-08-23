package router

import (
	crypto "github.com/tendermint/tendermint/crypto"
)

/*
 * Request
 */

type request struct {
	Frm crypto.PubKey `json:"from"`
	Pth string        `json:"path"`
	Dat []byte        `json:"data"`
	Sig []byte        `json:"signature"`
}

func createRequest(from crypto.PubKey, path string, data []byte, sig []byte) Request {
	out := request{
		Frm: from,
		Pth: path,
		Dat: data,
		Sig: sig,
	}

	return &out
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

/*
 * Trx Chk Request
 */

type trxChkRequest struct {
	Frm            crypto.PubKey `json:"from"`
	Pth            string        `json:"path"`
	DtaSizeInBytes int64         `json:"data_size_in_bytes"`
	Sig            []byte        `json:"signature"`
}

func createTrxChkRequest(from crypto.PubKey, path string, dataSizeInBytes int64, sig []byte) TrxChkRequest {
	out := trxChkRequest{
		Frm:            from,
		Pth:            path,
		DtaSizeInBytes: dataSizeInBytes,
		Sig:            sig,
	}

	return &out
}

// Path returns the path
func (obj *trxChkRequest) From() crypto.PubKey {
	return obj.Frm
}

// Path returns the path
func (obj *trxChkRequest) Path() string {
	return obj.Pth
}

// Data returns the data
func (obj *trxChkRequest) DataSizeInBytes() int64 {
	return obj.DtaSizeInBytes
}

// Signature returns the signature
func (obj *trxChkRequest) Signature() []byte {
	return obj.Sig
}
