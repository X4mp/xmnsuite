package commands

type startConfigs struct {
	conf Configs
	prs  []Node
}

func createStartConfigs(conf Configs, prs []Node) (StartConfigs, error) {
	out := startConfigs{
		conf: conf,
		prs:  prs,
	}

	return &out, nil
}

// Configs returns the configs
func (obj *startConfigs) Configs() Configs {
	return obj.conf
}

// HasPeers return true if there is peers, false otherwise
func (obj *startConfigs) HasPeers() bool {
	return len(obj.prs) > 0
}

// Peers return the peers
func (obj *startConfigs) Peers() []Node {
	return obj.prs
}
