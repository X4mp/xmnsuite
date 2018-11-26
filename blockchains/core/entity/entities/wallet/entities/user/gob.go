package user

import (
	"encoding/gob"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	// user:
	gob.Register(&user{})
}
