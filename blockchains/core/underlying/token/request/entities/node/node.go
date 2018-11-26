package node

import (
	"net"

	uuid "github.com/satori/go.uuid"
)

type node struct {
	UUID      *uuid.UUID `json:"id"`
	Pow       int        `json:"power"`
	IPAddress net.IP     `json:"ip"`
	Prt       int        `json:"port"`
}

func createNode(id *uuid.UUID, power int, ip net.IP, port int) Node {
	out := node{
		UUID:      id,
		Pow:       power,
		IPAddress: ip,
		Prt:       port,
	}

	return &out
}

func createNodeFromStorable(storable *storableNode) (Node, error) {

	nodeID, nodeIDErr := uuid.FromString(storable.ID)
	if nodeIDErr != nil {
		return nil, nodeIDErr
	}

	ip := net.ParseIP(storable.IP)
	out := createNode(&nodeID, storable.Pow, ip, storable.Port)
	return out, nil
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
