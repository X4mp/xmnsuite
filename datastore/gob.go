package datastore

import (
	"encoding/gob"

	"github.com/xmnservices/xmnsuite/datastore/keys"
	"github.com/xmnservices/xmnsuite/datastore/lists"
	"github.com/xmnservices/xmnsuite/datastore/objects"
	"github.com/xmnservices/xmnsuite/datastore/roles"
	"github.com/xmnservices/xmnsuite/datastore/users"
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
