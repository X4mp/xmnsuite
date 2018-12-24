package meta

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
)

type entityRequest struct {
	ent entity.Representation
	mp  map[string]entity.Representation
}

func createEntityRequest(ent entity.Representation, mp map[string]entity.Representation) EntityRequest {
	out := entityRequest{
		ent: ent,
		mp:  mp,
	}

	return &out
}

// RequestedBy returns the entity representation
func (obj *entityRequest) RequestedBy() entity.Representation {
	return obj.ent
}

// Map returns the map representation
func (obj *entityRequest) Map() map[string]entity.Representation {
	return obj.mp
}

// Add adds an entity representation
func (obj *entityRequest) Add(rep entity.Representation) EntityRequest {
	keyname := rep.MetaData().Keyname()
	if _, ok := obj.mp[keyname]; !ok {
		obj.mp[keyname] = rep
	}

	return obj
}
