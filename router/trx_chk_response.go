package router

type trxChkResponse struct {
	canBeExecuted   bool
	canBeAuthorized bool
	gazWanted       int64
	log             string
}

func createTrxChkResponse(canBeExecuted bool, canBeAuthorized bool, gazWanted int64, log string) TrxChkResponse {
	out := trxChkResponse{
		canBeExecuted:   canBeExecuted,
		canBeAuthorized: canBeAuthorized,
		gazWanted:       gazWanted,
		log:             log,
	}

	return &out
}

// CanBeExecuted returns true if the transaction can be executed, false otherwise
func (obj *trxChkResponse) CanBeExecuted() bool {
	return obj.canBeExecuted
}

// CanBeAuthorized returns true if the transaction can be executed, false otherwise
func (obj *trxChkResponse) CanBeAuthorized() bool {
	return obj.canBeAuthorized
}

// GazWanted returns the amount of gaz wanted to execute the transaction
func (obj *trxChkResponse) GazWanted() int64 {
	return obj.gazWanted
}

// Log returns the log
func (obj *trxChkResponse) Log() string {
	return obj.log
}
