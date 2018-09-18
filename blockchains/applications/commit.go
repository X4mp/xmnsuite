package applications

/*
 * CommitResponse
 */

type commitResponse struct {
	AppHsh    []byte `json:"app_hash"`
	BlkHeight int64  `json:"block_height"`
}

func createCommitResponse(appHash []byte, blkHeight int64) CommitResponse {
	out := commitResponse{
		AppHsh:    appHash,
		BlkHeight: blkHeight,
	}

	return &out
}

// AppHash returns the app hash
func (obj *commitResponse) AppHash() []byte {
	return obj.AppHsh
}

// BlockHeight returns the block height
func (obj *commitResponse) BlockHeight() int64 {
	return obj.BlkHeight
}
