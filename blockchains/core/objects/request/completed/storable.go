package completed

type storableRequest struct {
	ID              string `json:"id"`
	RequestID       string `json:"request_id"`
	ConcensusNeeded int    `json:"concensus_needed"`
	Approved        int    `json:"approved"`
	Disapproved     int    `json:"disapproved"`
	Neutral         int    `json:"neutral"`
}

func createStorableRequest(ins Request) *storableRequest {
	out := storableRequest{
		ID:              ins.ID().String(),
		RequestID:       ins.Request().ID().String(),
		ConcensusNeeded: ins.ConcensusNeeded(),
		Approved:        ins.Approved(),
		Disapproved:     ins.Disapproved(),
		Neutral:         ins.Neutral(),
	}

	return &out
}
