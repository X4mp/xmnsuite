package users

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/datastore/roles"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	roles.RegisterGob()
	gob.Register(&concreteUsers{})
}
