package active

type storableVote struct {
	ID     string `json:"id"`
	VoteID string `json:"vote"`
	Power  int    `json:"power"`
}

func createStorableVote(ins Vote) *storableVote {
	out := storableVote{
		ID:     ins.ID().String(),
		VoteID: ins.Vote().ID().String(),
		Power:  ins.Power(),
	}

	return &out
}
