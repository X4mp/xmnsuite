package sell

import (
	"encoding/gob"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	gob.Register(&wish{})
	gob.Register(&sell{})
}
