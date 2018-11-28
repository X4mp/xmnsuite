package request

type storableRequest struct {
	ID            string `json:"id"`
	FromUserID    string `json:"from_user_id"`
	NewEntityJS   []byte `json:"new_entity_js"`
	NewEntityName string `json:"new_entity_name"`
}

func createStorableRequest(req Request) (*storableRequest, error) {
	js, jsErr := reg.fromEntityToJSON(req.New(), req.NewName())
	if jsErr != nil {
		return nil, jsErr
	}

	out := storableRequest{
		ID:            req.ID().String(),
		FromUserID:    req.From().ID().String(),
		NewEntityJS:   js,
		NewEntityName: req.NewName(),
	}

	return &out, nil
}
