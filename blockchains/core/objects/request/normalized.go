package request

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
)

type normalizedRequest struct {
	ID               string             `json:"id"`
	From             user.Normalized    `json:"from"`
	SaveEntityJSON   []byte             `json:"save_entity_json"`
	DeleteEntityJSON []byte             `json:"delete_entity_json"`
	Reason           string             `json:"reason"`
	Keyname          keyname.Normalized `json:"keyname"`
}

func createNormalizedRequest(req Request) (*normalizedRequest, error) {

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

	from, fromErr := user.SDKFunc.CreateMetaData().Normalize()(req.From())
	if fromErr != nil {
		return nil, fromErr
	}

	kname, knameErr := keyname.SDKFunc.CreateMetaData().Normalize()(req.Keyname())
	if knameErr != nil {
		return nil, knameErr
	}

	out := normalizedRequest{
		ID:               req.ID().String(),
		From:             from,
		SaveEntityJSON:   toSaveEntity,
		DeleteEntityJSON: toDeleteEntity,
		Reason:           req.Reason(),
		Keyname:          kname,
	}

	return &out, nil
}
