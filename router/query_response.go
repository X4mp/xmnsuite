package router

type queryResponse struct {
	isSuccess    bool
	isAuthorized bool
	hasNFS       bool
	gazUsed      int64
	data         []byte
	log          string
}

func createQueryResponse(isSuccess bool, isAuthorized bool, hasNFS bool, gazUsed int64, log string, data []byte) QueryResponse {
	out := queryResponse{
		isSuccess:    isSuccess,
		isAuthorized: isAuthorized,
		hasNFS:       hasNFS,
		gazUsed:      gazUsed,
		log:          log,
		data:         data,
	}

	return &out
}

// GazUsed returns the amount of gaz used to execute the query
func (obj *queryResponse) GazUsed() int64 {
	return obj.gazUsed
}

// IsSuccess returns true if the query was successful, false otherwise
func (obj *queryResponse) IsSuccess() bool {
	return obj.isSuccess
}

// IsAuthorized returns true if the request was authorized, false otherwise
func (obj *queryResponse) IsAuthorized() bool {
	return obj.isAuthorized
}

// HasInsufficientFunds returns true if the user had insufficient funds to conver the query costs, false otherwise
func (obj *queryResponse) HasInsufficientFunds() bool {
	return obj.hasNFS
}

// Log returns log
func (obj *queryResponse) Log() string {
	return obj.log
}

// Data returns the data
func (obj *queryResponse) Data() []byte {
	return obj.data
}

// Matshal marshals the data to the pointer
func (obj *queryResponse) Matshal(v interface{}) error {
	return nil
}
