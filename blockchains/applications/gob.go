package applications

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/datastore"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	datastore.RegisterGob()
	gob.Register(&state{})
}
