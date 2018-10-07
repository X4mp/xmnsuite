package objects

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/keys"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	keys.RegisterGob()
	gob.Register(&concreteObjects{})
}
