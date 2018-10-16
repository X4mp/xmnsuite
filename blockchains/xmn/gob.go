package xmn

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	// dependencies:
	applications.RegisterGob()

	// genesis:
	gob.Register(&genesis{})
	gob.Register(&token{})
	gob.Register(&initialDeposit{})

	// wallet:
	gob.Register(&wallet{})
	gob.Register(&user{})
}
