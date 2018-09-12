package sdk

import (
	applications "github.com/xmnservices/xmnsuite/applications"
	lua "github.com/yuin/gopher-lua"
)

// SDK represents the SDK module
type SDK interface {
}

// CreateParams represents the create params
type CreateParams struct {
	Context *lua.LState
	Client  applications.Client
}

// SDKFunc represents the sdk module SDK func
var SDKFunc = struct {
	Create func(params CreateParams) SDK
}{
	Create: func(params CreateParams) SDK {
		out := createModule(params.Context, params.Client)
		return out
	},
}
