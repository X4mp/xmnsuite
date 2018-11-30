package node

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/link"
)

type node struct {
	UUID      *uuid.UUID `json:"id"`
	Lnk       link.Link  `json:"link"`
	Pow       int        `json:"power"`
	IPAddress net.IP     `json:"ip"`
	Prt       int        `json:"port"`
}

func createNode(id *uuid.UUID, lnk link.Link, power int, ip net.IP, port int) Node {
	out := node{
		UUID:      id,
		Lnk:       lnk,
		Pow:       power,
		IPAddress: ip,
		Prt:       port,
	}

	return &out
}

func createNodeFromNormalized(normalized *normalizedNode) (Node, error) {
	nodeID, nodeIDErr := uuid.FromString(normalized.ID)
	if nodeIDErr != nil {
		return nil, nodeIDErr
	}

	lnkIns, lnkInsErr := link.SDKFunc.CreateMetaData().Denormalize()(normalized.Link)
	if lnkInsErr != nil {
		return nil, lnkInsErr
	}

	if lnk, ok := lnkIns.(link.Link); ok {
		ip := net.ParseIP(normalized.IP)
		out := createNode(&nodeID, lnk, normalized.Pow, ip, normalized.Port)
		return out, nil
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Node instance", lnkIns.ID().String())
	return nil, errors.New(str)

}

// ID returns the ID
func (obj *node) ID() *uuid.UUID {
	return obj.UUID
}

// Power returns the node's power
func (obj *node) Power() int {
	return obj.Pow
}

// IP returns the node's IP
func (obj *node) IP() net.IP {
	return obj.IPAddress
}

// Port returns the node's port
func (obj *node) Port() int {
	return obj.Prt
}

// Link returns the link instance
func (obj *node) Link() link.Link {
	return obj.Lnk
}
