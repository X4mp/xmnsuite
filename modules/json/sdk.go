package json

import (
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type module struct {
}

// JSON represents the json module
type JSON interface {
}

// CreateParams represents the create params
type CreateParams struct {
	Context *lua.LState
}

// SDKFunc represents the json module SDK func
var SDKFunc = struct {
	Create func(params CreateParams) JSON
}{
	Create: func(params CreateParams) JSON {
		luajson.Preload(params.Context)
		return &module{}
	},
}
