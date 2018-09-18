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
	Siz int64  `json:"size"`
	Ver string `json:"version"`
}

func createInfoResponse(size int64, version string) InfoResponse {
	out := infoResponse{
		Siz: size,
		Ver: version,
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
