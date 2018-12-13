package meta

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
)

// SDKFunc represents the forex SDK func
var SDKFunc = struct {
	CreateMetaData func() meta.Meta
}{
	CreateMetaData: func() meta.Meta {
		out, outErr := createMeta()
		if outErr != nil {
			panic(outErr)
		}

		return out
	},
}
