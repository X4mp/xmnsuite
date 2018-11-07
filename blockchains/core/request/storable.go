package request

type storableRequest struct {
	ID          string `json:"id"`
	FromUserID  string `json:"from_user_id"`
	NewEntityJS []byte `json:"new_entity_js"`
}

func createStorableRequest(req Request) (*storableRequest, error) {
	js, jsErr := cdc.MarshalJSON(req)
	if jsErr != nil {
		return nil, jsErr
	}

	out := storableRequest{
		ID:          req.ID().String(),
		FromUserID:  req.From().ID().String(),
		NewEntityJS: js,
	}

	return &out, nil
}
