package entity

import (
	"encoding/gob"
)

func init() {
	gob.Register(&testEntity{})
}
