package messages

import (
	"encoding/json"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/configs"
)

// CreateParams represents the create params
type CreateParams struct {
	EncryptedConf string
	Pass          string
	Host          string
}

// RenderEntityParams represents the render entity params
type RenderEntityParams struct {
	Meta entity.MetaData
	Ins  entity.Entity
}

// SDKFunc represents the messages SDK func
var SDKFunc = struct {
	Create       func(params CreateParams) (configs.Configs, applications.Client, entity.Repository, entity.Service)
	RenderEntity func(params RenderEntityParams) json.RawMessage
}{
	Create: func(params CreateParams) (configs.Configs, applications.Client, entity.Repository, entity.Service) {
		return create(params.EncryptedConf, params.Pass, params.Host)
	},
	RenderEntity: func(params RenderEntityParams) json.RawMessage {
		out, outErr := renderEntity(params.Meta, params.Ins)
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
