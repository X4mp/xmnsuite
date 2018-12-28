package web

import (
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/crypto"
)

// Web represents a web server
type Web interface {
	Start() error
	Stop() error
}

// CreateParams represents the create params
type CreateParams struct {
	Port   int
	Client applications.Client
	Meta   meta.Meta
	PK     crypto.PrivateKey
}

// SDKFunc represents the web server
var SDKFunc = struct {
	Create func(params CreateParams) Web
}{
	Create: func(params CreateParams) Web {
		out := createWeb(
			params.Port,
			params.Meta,
			params.Client,
			params.PK,
		)

		return out
	},
}
