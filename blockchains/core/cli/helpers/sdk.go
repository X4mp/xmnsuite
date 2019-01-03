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
	PrintSuccessWithInstance: func(params PrintSuccessWithInstanceParams) {
		js, jsErr := json.MarshalIndent(params.Ins, "", "    ")
		if jsErr != nil {
			panic(jsErr)
		}

		out := fmt.Sprintf("\n************ XMN - SUCCESS ************\n")
		out = fmt.Sprintf("%s%s", out, params.Message)
		out = fmt.Sprintf("%s\n-------------------------------------\n%s\n", out, string(js))
		out = fmt.Sprintf("%s\n********** END XMN - SUCCESS **********\n", out)
		fmt.Printf(out)
	},
	PrintError: func(params PrintErrorParams) {
		out := fmt.Sprintf("\n************ XMN - ERROR ************\n")
		out = fmt.Sprintf("%s%s", out, params.Message)
		out = fmt.Sprintf("%s\n********** END XMN - ERROR **********\n", out)
		fmt.Printf(out)
	},
}
