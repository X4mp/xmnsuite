package router

type trxResponse struct {
	isSuccess    bool
	isAuthorized bool
	isNFS        bool
	tags         map[string][]byte
	gazUsed      int64
	log          string
}

func createTrxResponse(isSuccess bool, isAuthorized bool, isNFS bool, tags map[string][]byte, gazUsed int64, log string) TrxResponse {
	out := trxResponse{
		isSuccess:    isSuccess,
		isAuthorized: isAuthorized,
		isNFS:        isNFS,
		tags:         tags,
		gazUsed:      gazUsed,
		log:          log,
	}

	return &out
}

// IsSuccess returns true if the transaction is successful, false otherwise
func (obj *trxResponse) IsSuccess() bool {
	return obj.isSuccess
}

// IsAuthorized returns true if the transaction is authorized, false otherwise
func (obj *trxResponse) IsAuthorized() bool {
	return obj.isAuthorized
}

// HasInsufficientFunds returns true if the user had insufficient funds, false otherwise
func (obj *trxResponse) HasInsufficientFunds() bool {
	return obj.isNFS
}

// Tags returns the tags
func (obj *trxResponse) Tags() map[string][]byte {
	return obj.tags
}

// GazUsed returns the amount of gaz used
func (obj *trxResponse) GazUsed() int64 {
	return obj.gazUsed
}

// Log returns the logs
func (obj *trxResponse) Log() string {
	return obj.log
}
