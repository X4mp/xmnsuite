package lists

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/datastore/keys"
	"github.com/xmnservices/xmnsuite/datastore/objects"
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
