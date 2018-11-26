package sell

// Daemon represents the sell daemon that tries to match the sell orders of this blockchain to the ones on other blockchains
type Daemon interface {
	Start() error
	Stop() error
}
