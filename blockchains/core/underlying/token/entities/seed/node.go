package seed

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/link"
)

type seed struct {
	UUID      *uuid.UUID `json:"id"`
	Lnk       link.Link  `json:"link"`
	IPAddress net.IP     `json:"ip"`
	Prt       int        `json:"port"`
}

func createSeed(id *uuid.UUID, lnk link.Link, ip net.IP, port int) Seed {
	out := seed{
		UUID:      id,
		Lnk:       lnk,
		IPAddress: ip,
		Prt:       port,
	}

	return &out
}

func createSeedFromNormalized(normalized *normalizedSeed) (Seed, error) {
	seedID, seedIDErr := uuid.FromString(normalized.ID)
	if seedIDErr != nil {
		return nil, seedIDErr
	}

	lnkIns, lnkInsErr := link.SDKFunc.CreateMetaData().Denormalize()(normalized.Link)
	if lnkInsErr != nil {
		return nil, lnkInsErr
	}

	if lnk, ok := lnkIns.(link.Link); ok {
		ip := net.ParseIP(normalized.IP)
		out := createSeed(&seedID, lnk, ip, normalized.Port)
		return out, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Seed instance", lnkIns.ID().String())
	return nil, errors.New(str)

}

// ID returns the ID
func (obj *seed) ID() *uuid.UUID {
	return obj.UUID
}

// IP returns the seed's IP
func (obj *seed) IP() net.IP {
	return obj.IPAddress
}

// Port returns the seed's port
func (obj *seed) Port() int {
	return obj.Prt
}

// Link returns the link instance
func (obj *seed) Link() link.Link {
	return obj.Lnk
}
