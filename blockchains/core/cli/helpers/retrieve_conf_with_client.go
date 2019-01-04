package helpers

import (
	"errors"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/tendermint"
	"github.com/xmnservices/xmnsuite/configs"
)

func retrieveConfWithClient(c *cliapp.Context) (configs.Configs, applications.Client, error) {
	// retrieve the configurations:
	fileAsString := c.String("file")
	confRepository := configs.SDKFunc.CreateRepository()
	conf, confErr := confRepository.Retrieve(fileAsString, c.String("pass"))
	if confErr != nil {
		str := fmt.Sprintf("the given file (%s) either does not exist or the given password is invalid", fileAsString)
		return nil, nil, errors.New(str)
	}

	// create the blockchain client:
	client := tendermint.SDKFunc.CreateClient(tendermint.CreateClientParams{
		IPAsString: c.String("host"),
	})

	return conf, client, nil
}
