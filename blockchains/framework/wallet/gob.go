package wallet

import (
	"encoding/gob"
)

func init() {
	RegisterGob()
}

// RegisterGob registers the hashtree for gob
func RegisterGob() {
	// wallet:
	gob.Register(&wallet{})
}
