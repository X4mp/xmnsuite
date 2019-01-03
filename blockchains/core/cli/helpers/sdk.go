package helpers

import (
	"encoding/json"
	"fmt"

	cliapp "github.com/urfave/cli"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

// ProcessWalletRequestParams represents the proces wallet request params
type ProcessWalletRequestParams struct {
	CLIContext           *cliapp.Context
	EntityRepresentation entity.Representation
	Storable             interface{}
}

// PrintSuccessNewInstanceParams represents the print success new instance params
type PrintSuccessNewInstanceParams struct {
	Ins     interface{}
	Message string
}

// SDKFunc represents the helpers SDK func
var SDKFunc = struct {
	ProcessWalletRequest    func(params ProcessWalletRequestParams) request.Normalized
	PrintSuccessNewInstance func(params PrintSuccessNewInstanceParams)
}{
	ProcessWalletRequest: func(params ProcessWalletRequestParams) request.Normalized {
		req, reqErr := processWalletRequest(params.CLIContext, params.EntityRepresentation, params.Storable)
		if reqErr != nil {
			panic(reqErr)
		}

		return req
	},
	PrintSuccessNewInstance: func(params PrintSuccessNewInstanceParams) {
		js, jsErr := json.MarshalIndent(params.Ins, "", "    ")
		if jsErr != nil {
			panic(jsErr)
		}

		out := fmt.Sprintf("\n************ XMN - SUCCESS ************\n")
		out = fmt.Sprintf("%s%s", out, params.Message)
		out = fmt.Sprintf("%s\n--------------New instance:-------------\n%s\n", out, string(js))
		out = fmt.Sprintf("%s\n********** END XMN - SUCCESS **********\n", out)
		fmt.Printf(out)
	},
}
