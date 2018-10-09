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
	Ver string `json:"version"`
	St  State  `json:"state"`
}

func createInfoResponse(version string, st State) InfoResponse {
	out := infoResponse{
		Ver: version,
		St:  st,
	}

	return &out
}

// Version returns the blockchain version
func (obj *infoResponse) Version() string {
	return obj.Ver
}

// State returns the state
func (obj *infoResponse) State() State {
	return obj.St
}
