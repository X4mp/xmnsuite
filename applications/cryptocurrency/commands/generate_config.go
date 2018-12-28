package commands

import (
	"github.com/xmnservices/xmnsuite/configs"
)

func generateConfigs(pass string, retypedPass string, filename string) (configs.Configs, error) {
	// create the configs:
	conf := configs.SDKFunc.Generate()

	// create the service:
	service := configs.SDKFunc.CreateService()

	// save the configs:
	saveErr := service.Save(conf, filename, pass, retypedPass)
	if saveErr != nil {
		return nil, saveErr
	}

	// return the configs:
	return conf, nil
}
