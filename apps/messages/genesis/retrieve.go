package genesis

import (
	"encoding/json"

	"github.com/xmnservices/xmnsuite/apps/messages"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
)

type retrievePayload struct {
	Host    string `json:"host"`
	Configs string `json:"configs"`
	Pass    string `json:"pass"`
}

func retrieve(input json.RawMessage) (interface{}, error) {
	ptr := new(retrievePayload)
	jsErr := json.Unmarshal(input, ptr)
	if jsErr != nil {
		return nil, jsErr
	}

	// create the services:
	_, _, entityRepository, _ := messages.SDKFunc.Create(messages.CreateParams{
		EncryptedConf: ptr.Configs,
		Pass:          ptr.Pass,
		Host:          ptr.Host,
	})

	// create the genesis repository:
	genRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
		EntityRepository: entityRepository,
	})

	// retrieve the genesis:
	gen, genErr := genRepository.Retrieve()
	if genErr != nil {
		return nil, genErr
	}

	// render:
	return messages.SDKFunc.RenderEntity(messages.RenderEntityParams{
		Meta: genesis.SDKFunc.CreateMetaData(),
		Ins:  gen,
	}), nil

}
