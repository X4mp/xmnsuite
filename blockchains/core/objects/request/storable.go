package request

type storableRequest struct {
	ID          string `json:"id"`
	FromUserID  string `json:"from_user_id"`
	NewEntityJS []byte `json:"new_entity_js"`
	Reason      string `json:"reason"`
	KeynameID   string `json:"keyname"`
}

func createStorableRequest(req Request) (*storableRequest, error) {
	js, jsErr := reg.fromEntityToJSON(req.New(), req.Keyname().Name())
	if jsErr != nil {
		return nil, jsErr
	}

	out := storableRequest{
		ID:          req.ID().String(),
		FromUserID:  req.From().ID().String(),
		NewEntityJS: js,
		Reason:      req.Reason(),
		KeynameID:   req.Keyname().ID().String(),
	}

	return &out, nil
}
