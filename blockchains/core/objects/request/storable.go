package request

type storableRequest struct {
	ID               string `json:"id"`
	FromUserID       string `json:"from_user_id"`
	SaveEntityJSON   []byte `json:"save_entity_json"`
	DeleteEntityJSON []byte `json:"delete_entity_json"`
	Reason           string `json:"reason"`
	KeynameID        string `json:"keyname"`
}

func createStorableRequest(req Request) (*storableRequest, error) {

	var toSaveEntity []byte
	var toDeleteEntity []byte

	if req.HasSave() {
		saveJSON, saveJSONErr := reg.fromEntityToJSON(req.Save(), req.Keyname().Name())
		if saveJSONErr != nil {
			return nil, saveJSONErr
		}

		toSaveEntity = saveJSON
	}

	if req.HasDelete() {
		deleteJSON, deleteJSONErr := reg.fromEntityToJSON(req.Delete(), req.Keyname().Name())
		if deleteJSONErr != nil {
			return nil, deleteJSONErr
		}

		toDeleteEntity = deleteJSON
	}

	out := storableRequest{
		ID:               req.ID().String(),
		FromUserID:       req.From().ID().String(),
		SaveEntityJSON:   toSaveEntity,
		DeleteEntityJSON: toDeleteEntity,
		Reason:           req.Reason(),
		KeynameID:        req.Keyname().ID().String(),
	}

	return &out, nil
}
