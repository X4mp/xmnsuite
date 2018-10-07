package lists

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/keys"
	"github.com/xmnservices/xmnsuite/objects"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	keys.RegisterGob()
	objects.RegisterGob()
	gob.Register(&concreteLists{})
}
