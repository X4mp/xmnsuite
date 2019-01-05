package helpers

import (
	"encoding/json"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/configs"
)

// ProcessWalletRequestParams represents the proces wallet request params
type ProcessWalletRequestParams struct {
	CLIContext           *cliapp.Context
	EntityRepresentation entity.Representation
	Storable             interface{}
}

// RetrieveConfWithClientParams represents the retrieveConfWithClient params
type RetrieveConfWithClientParams struct {
	CLIContext *cliapp.Context
}

// SaveRequestParams represents the SaveRequest params
type SaveRequestParams struct {
	CLIContext           *cliapp.Context
	EntityRepresentation entity.Representation
	Ins                  entity.Entity
}

// PrintSuccessWithInstanceParams represents the print success new instance params
type PrintSuccessWithInstanceParams struct {
	Ins     interface{}
	Message string
}

// PrintErrorParams represents the print error params
type PrintErrorParams struct {
	Message string
}

// SDKFunc represents the helpers SDK func
var SDKFunc = struct {
	ProcessWalletRequest     func(params ProcessWalletRequestParams) request.Normalized
	RetrieveConfWithClient   func(params RetrieveConfWithClientParams) (configs.Configs, applications.Client)
	SaveRequest              func(params SaveRequestParams) request.Request
	PrintSuccessWithInstance func(params PrintSuccessWithInstanceParams)
	PrintError               func(params PrintErrorParams)
}{
	ProcessWalletRequest: func(params ProcessWalletRequestParams) request.Normalized {
		req, reqErr := processWalletRequest(params.CLIContext, params.EntityRepresentation, params.Storable)
		if reqErr != nil {
			panic(reqErr)
		}

		return req
	},
	RetrieveConfWithClient: func(params RetrieveConfWithClientParams) (configs.Configs, applications.Client) {
		conf, client, err := retrieveConfWithClient(params.CLIContext)
		if err != nil {
			panic(err)
		}

		return conf, client
	},
	SaveRequest: func(params SaveRequestParams) request.Request {
		req, reqErr := saveRequest(params.CLIContext, params.EntityRepresentation, params.Ins)
		if reqErr != nil {
			panic(reqErr)
		}

		return req
	},
	PrintSuccessWithInstance: func(params PrintSuccessWithInstanceParams) {
		js, jsErr := json.MarshalIndent(params.Ins, "", "    ")
		if jsErr != nil {
			panic(jsErr)
		}

		out := fmt.Sprintf("\n%s\n", js)
		fmt.Printf(out)
	},
	PrintError: func(params PrintErrorParams) {
		out := fmt.Sprintf("\n************ ERROR ************\n")
		out = fmt.Sprintf("%s%s", out, params.Message)
		out = fmt.Sprintf("%s\n********** END ERROR **********\n", out)
		fmt.Printf(out)
	},
}
