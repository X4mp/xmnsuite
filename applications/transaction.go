package applications

import (
	"errors"
	"fmt"

	crypto "github.com/xmnservices/xmnsuite/crypto"
)

/*
 * TransactionRequest
 */

type transactionRequest struct {
	Res Resource         `json:"resource"`
	Ptr ResourcePointer  `json:"resource_pointer"`
	Sig crypto.Signature `json:"signature"`
}

func createTransactionRequestWithResource(res Resource, sig crypto.Signature) (TransactionRequest, error) {
	if !res.Pointer().From().Equals(sig.PublicKey(res.Hash())) {
		str := fmt.Sprintf("the signature and resource hash could not be validated by the resource pointer's public key")
		return nil, errors.New(str)
	}

	out := transactionRequest{
		Res: res,
		Ptr: nil,
		Sig: sig,
	}

	return &out, nil
}

func createTransactionRequestWithResourcePointer(ptr ResourcePointer, sig crypto.Signature) (TransactionRequest, error) {
	if !ptr.From().Equals(sig.PublicKey(ptr.Hash())) {
		str := fmt.Sprintf("the signature and resource pointer hash could not be validated by the resource pointer's public key")
		return nil, errors.New(str)
	}

	out := transactionRequest{
		Res: nil,
		Ptr: ptr,
		Sig: sig,
	}

	return &out, nil
}

// Resource returns the resource, if any
func (obj *transactionRequest) Resource() Resource {
	return obj.Res
}

// Pointer represents the resource pointer, if any
func (obj *transactionRequest) Pointer() ResourcePointer {
	return obj.Ptr
}

// Signature returns the signature
func (obj *transactionRequest) Signature() crypto.Signature {
	return obj.Sig
}

/*
 * TransactionResponse
 */

type transactionResponse struct {
	Cod    int               `json:"code"`
	Lg     string            `json:"log"`
	GzUsed int64             `json:"gaz_used"`
	Tgs    map[string][]byte `json:"tags"`
}

func createFreeTransactionResponse(code int, log string) (TransactionResponse, error) {

	if !isCodeValid(code) {
		str := fmt.Sprintf("the code (%d) is invalid", code)
		return nil, errors.New(str)
	}

	out := transactionResponse{
		Cod:    code,
		Lg:     log,
		GzUsed: 0,
		Tgs:    map[string][]byte{},
	}

	return &out, nil
}

func createTransactionResponse(code int, log string, gazUsed int64, tags map[string][]byte) (TransactionResponse, error) {

	if !isCodeValid(code) {
		str := fmt.Sprintf("the code (%d) is invalid", code)
		return nil, errors.New(str)
	}

	out := transactionResponse{
		Cod:    code,
		Lg:     log,
		GzUsed: gazUsed,
		Tgs:    tags,
	}

	return &out, nil
}

// Code returns the status code
func (obj *transactionResponse) Code() int {
	return obj.Cod
}

// Log returns the log
func (obj *transactionResponse) Log() string {
	return obj.Lg
}

// GazUsed returns the gaz used to execute the transaction
func (obj *transactionResponse) GazUsed() int64 {
	return obj.GzUsed
}

// Tags returns the transaction tags
func (obj *transactionResponse) Tags() map[string][]byte {
	return obj.Tgs
}
