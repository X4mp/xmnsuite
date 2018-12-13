package commands

import (
	"fmt"
	"net"
)

type node struct {
	ip   net.IP
	port int
}

func createNode(ip net.IP, port int) (Node, error) {
	out := node{
		ip:   ip,
		port: port,
	}

	return &out, nil
}

// IP returns the IP
func (obj *node) IP() net.IP {
	return obj.ip
}

// Port returns the port
func (obj *node) Port() int {
	return obj.port
}

// String returns the node string
func (obj *node) String() string {
	return fmt.Sprintf("tcp://%s:%d", obj.IP().String(), obj.Port())
}
