package uuid

import (
	lua "github.com/yuin/gopher-lua"
)

// UUID represents the uuid module
type UUID interface {
}

// CreateParams represents the context params:
type CreateParams struct {
	Context *lua.LState
}

// SDKFunc represents the uuid SDK func
var SDKFunc = struct {
	Create func(params CreateParams) UUID
}{
	Create: func(params CreateParams) UUID {
		mod := createModule(params.Context)
		return mod
	},
}
