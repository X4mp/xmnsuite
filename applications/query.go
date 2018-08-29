package applications

import (
	"errors"
	"fmt"
)

/*
 * QueryRequest
 */

type queryRequest struct {
	Ptr ResourcePointer `json:"pointer"`
	Sig []byte          `json:"sig"`
}

func createQueryRequest(ptr ResourcePointer, sig []byte) (QueryRequest, error) {

	if !ptr.From().VerifyBytes(ptr.Hash(), sig) {
		str := fmt.Sprintf("the signature and resource pointer's hash could not be validated by the resource pointer's public key")
		return nil, errors.New(str)
	}

	out := queryRequest{
		Ptr: ptr,
		Sig: sig,
	}

	return &out, nil
}

// Pointer returns the resource pointer
func (obj *queryRequest) Pointer() ResourcePointer {
	return obj.Ptr
}

// Signature returns the signature
func (obj *queryRequest) Signature() []byte {
	return obj.Sig
}

/*
 * QueryResponse
 */

type queryResponse struct {
	Cod int    `json:"code"`
	Lg  string `json:"log"`
	K   string `json:"key"`
	Val []byte `json:"value"`
}

func createEmptyQueryResponse(code int, log string) (QueryResponse, error) {

	if !isCodeValid(code) {
		str := fmt.Sprintf("the code (%d) is invalid", code)
		return nil, errors.New(str)
	}

	out := queryResponse{
		Cod: code,
		Lg:  log,
		K:   "",
		Val: []byte(""),
	}

	return &out, nil
}

func createQueryResponse(code int, log string, key string, value []byte) (QueryResponse, error) {

	if !isCodeValid(code) {
		str := fmt.Sprintf("the code (%d) is invalid", code)
		return nil, errors.New(str)
	}

	out := queryResponse{
		Cod: code,
		Lg:  log,
		K:   key,
		Val: value,
	}

	return &out, nil
}

// Code returns the status code
func (obj *queryResponse) Code() int {
	return obj.Cod
}

// Log returns the log
func (obj *queryResponse) Log() string {
	return obj.Lg
}

// Key returns the key
func (obj *queryResponse) Key() string {
	return obj.K
}

// Value returns the value
func (obj *queryResponse) Value() []byte {
	return obj.Val
}
