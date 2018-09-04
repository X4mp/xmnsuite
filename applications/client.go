package applications

type clientTransactionResponse struct {
	Chk TransactionResponse `json:"check_response"`
	Trx TransactionResponse `json:"transaction_response"`
	Ht  int64               `json:"height"`
	Hsh []byte              `json:"hash"`
}

func createClientTransactionResponse(chk TransactionResponse, trx TransactionResponse, height int64, hash []byte) ClientTransactionResponse {
	out := clientTransactionResponse{
		Chk: chk,
		Trx: trx,
		Ht:  height,
		Hsh: hash,
	}

	return &out
}

// Check returns the check response
func (obj *clientTransactionResponse) Check() TransactionResponse {
	return obj.Chk
}

// Transaction returns the transaction response
func (obj *clientTransactionResponse) Transaction() TransactionResponse {
	return obj.Trx
}

// Height returns the height
func (obj *clientTransactionResponse) Height() int64 {
	return obj.Ht
}

// Hash returns the hash
func (obj *clientTransactionResponse) Hash() []byte {
	return obj.Hsh
}
