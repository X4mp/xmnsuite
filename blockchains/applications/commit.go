package applications

/*
 * CommitResponse
 */

type commitResponse struct {
	AppHsh     []byte `json:"app_hash"`
	PrevAppHsh []byte `json:"prev_app_hash"`
	BlkHeight  int64  `json:"block_height"`
}

func createCommitResponse(prevAppHash []byte, appHash []byte, blkHeight int64) CommitResponse {
	out := commitResponse{
		AppHsh:     appHash,
		PrevAppHsh: prevAppHash,
		BlkHeight:  blkHeight,
	}

	return &out
}

// AppHash returns the app hash
func (obj *commitResponse) AppHash() []byte {
	return obj.AppHsh
}

// PrevAppHash returns the previous app hash
func (obj *commitResponse) PrevAppHash() []byte {
	return obj.PrevAppHsh
}

// BlockHeight returns the block height
func (obj *commitResponse) BlockHeight() int64 {
	return obj.BlkHeight
}
