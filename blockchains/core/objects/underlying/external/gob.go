package external

import (
	"encoding/gob"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the external resource
func RegisterGob() {
	gob.Register(&external{})
}
