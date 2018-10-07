package datastore

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/keys"
	"github.com/xmnservices/xmnsuite/lists"
	"github.com/xmnservices/xmnsuite/objects"
	"github.com/xmnservices/xmnsuite/roles"
	"github.com/xmnservices/xmnsuite/users"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	keys.RegisterGob()
	lists.RegisterGob()
	objects.RegisterGob()
	users.RegisterGob()
	roles.RegisterGob()
	gob.Register(&concreteDataStore{})
}
