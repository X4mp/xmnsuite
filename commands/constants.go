package commands

import uuid "github.com/satori/go.uuid"

type constants struct {
	namespace     string
	name          string
	id            *uuid.UUID
	routePrefix   string
	routerRoleKey string
}

func createConstants(namespace string, name string, id *uuid.UUID, routePrefix string, routerRoleKey string) (Constants, error) {
	out := constants{
		namespace:     namespace,
		name:          name,
		id:            id,
		routePrefix:   routePrefix,
		routerRoleKey: routerRoleKey,
	}

	return &out, nil
}

// Namespace returns the namespace
func (obj *constants) Namespace() string {
	return obj.namespace
}

// Name returns the name
func (obj *constants) Name() string {
	return obj.name
}

// ID returns the ID
func (obj *constants) ID() *uuid.UUID {
	return obj.id
}

// RoutePrefix returns the route prefix
func (obj *constants) RoutePrefix() string {
	return obj.routePrefix
}

// RouterRoleKey returns the router role key
func (obj *constants) RouterRoleKey() string {
	return obj.routerRoleKey
}
