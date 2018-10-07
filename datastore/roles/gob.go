package roles

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/datastore/lists"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	lists.RegisterGob()
	gob.Register(&concreteRoles{})
}
