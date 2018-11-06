package link

import (
	"net"

	uuid "github.com/satori/go.uuid"
	tcrypto "github.com/tendermint/tendermint/crypto"
)

type node struct {
	UUID      *uuid.UUID     `json:"id"`
	PKey      tcrypto.PubKey `json:"pubkey"`
	Pow       int            `json:"power"`
	IPAddress net.IP         `json:"ip"`
	Prt       int            `json:"port"`
}

func createNode(id *uuid.UUID, pubKey tcrypto.PubKey, power int, ip net.IP, port int) Node {
	out := node{
		UUID:      id,
		PKey:      pubKey,
		Pow:       power,
		IPAddress: ip,
		Prt:       port,
	}

	return &out
}

// ID returns the ID
func (obj *node) ID() *uuid.UUID {
	return obj.UUID
}

// PublicKey returns the node's public key
func (obj *node) PublicKey() tcrypto.PubKey {
	return obj.PKey
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
