package node

type storableNode struct {
	ID     string `json:"id"`
	PubKey string `json:"pubkey"`
	Pow    int    `json:"power"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
}

func createStorableNode(ins Node) *storableNode {
	out := storableNode{
		ID:     ins.ID().String(),
		PubKey: string(ins.PublicKey().Bytes()),
		Pow:    ins.Power(),
		IP:     ins.IP().String(),
		Port:   ins.Port(),
	}

	return &out
}
