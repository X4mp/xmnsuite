package messages

import (
	"encoding/json"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/configs"
)

func create(encryptedConf string, pass string, host string) (configs.Configs, applications.Client, entity.Repository, entity.Service) {
	// decrypt the configs:
	conf := configs.SDKFunc.Decrypt(configs.DecryptParams{
		Data: encryptedConf,
		Pass: pass,
	})

	// create the client:
	client := tendermint.SDKFunc.CreateClient(tendermint.CreateClientParams{
		IPAsString: host,
	})

	// create the entity repository:
	entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
		PK:     conf.WalletPK(),
		Client: client,
	})

	// create the entity service:
	entityService := entity.SDKFunc.CreateSDKService(entity.CreateSDKServiceParams{
		PK:     conf.WalletPK(),
		Client: client,
	})

	return conf, client, entityRepository, entityService
}

func renderEntity(met entity.MetaData, ins entity.Entity) (json.RawMessage, error) {
	// normalize:
	normalize, normalizeErr := met.Normalize()(ins)
	if normalizeErr != nil {
		return nil, normalizeErr
	}

	// convert to json:
	js, jsErr := json.Marshal(normalize)
	if jsErr != nil {
		return nil, jsErr
	}

	return json.RawMessage(js), nil
}
