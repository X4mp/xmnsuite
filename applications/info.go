package applications

/*
 * InfoRequest
 */

type infoRequest struct {
	Ver string `json:"version"`
}

func createInfoRequest(version string) InfoRequest {
	out := infoRequest{
		Ver: version,
	}

	return &out
}

// Version returns the version
func (obj *infoRequest) Version() string {
	return obj.Ver
}

/*
 * InfoResponse
 */

type infoResponse struct {
	Siz           int64  `json:"size"`
	Ver           string `json:"version"`
	LstBlkHeight  int64  `json:"last_block_height"`
	LstBlkAppHash []byte `json:"last_bock_app_hash"`
}

func createInfoResponse(size int64, version string, lastBlkHeight int64, lastBlkAppHash []byte) InfoResponse {
	out := infoResponse{
		Siz:           size,
		Ver:           version,
		LstBlkHeight:  lastBlkHeight,
		LstBlkAppHash: lastBlkAppHash,
	}

	return &out
}

// Size returns the blockchain size
func (obj *infoResponse) Size() int64 {
	return obj.Siz
}

// Version returns the blockchain version
func (obj *infoResponse) Version() string {
	return obj.Ver
}

// LastBlockHeight returns the amount of transactions the last block had
func (obj *infoResponse) LastBlockHeight() int64 {
	return obj.LstBlkHeight
}

// LastBlockAppHash returns the application hash the last block had
func (obj *infoResponse) LastBlockAppHash() []byte {
	return obj.LstBlkAppHash
}
