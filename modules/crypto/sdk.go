package crypto

import (
	lua "github.com/yuin/gopher-lua"
)

// Crypto represents the crypto module
type Crypto interface {
}

// CreateParams represents the context params:
type CreateParams struct {
	Context *lua.LState
}

// SDKFunc represents the crypto SDK func
var SDKFunc = struct {
	Create func(params CreateParams) Crypto
}{
	Create: func(params CreateParams) Crypto {
		mod := createModule(params.Context)
		return mod
	},
}
