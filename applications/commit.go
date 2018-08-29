package applications

/*
 * CommitResponse
 */

type commitResponse struct {
	AppHsh []byte `json:"app_hash"`
}

func createCommitResponse(appHash []byte) CommitResponse {
	out := commitResponse{
		AppHsh: appHash,
	}

	return &out
}

// AppHash returns the app hash
func (obj *commitResponse) AppHash() []byte {
	return obj.AppHsh
}
