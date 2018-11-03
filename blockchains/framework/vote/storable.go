package vote

type storableVote struct {
	ID      string `json:"id"`
	ReqID   string `json:"request_id"`
	VoterID string `json:"voter_id"`
	IsAppr  bool   `json:"is_approved"`
}

func createStorableVote(vote Vote) *storableVote {
	out := storableVote{
		ID:      vote.ID().String(),
		ReqID:   vote.Request().ID().String(),
		VoterID: vote.Voter().ID().String(),
		IsAppr:  vote.IsApproved(),
	}

	return &out
}
