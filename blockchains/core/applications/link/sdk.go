package link

// Daemon represents the link daemon that download the nodes related to links on the linked blockchains
type Daemon interface {
	Start() error
	Stop() error
}
