package active

type storableRequest struct {
	ID              string `json:"id"`
	RequestID       string `json:"request_id"`
	ConcensusNeeded int    `json:"concensus_needed"`
}

func createStorable(ins Request) *storableRequest {
	out := storableRequest{
		ID:              ins.ID().String(),
		RequestID:       ins.Request().ID().String(),
		ConcensusNeeded: ins.ConcensusNeeded(),
	}

	return &out
}
