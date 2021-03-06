package node

type storableNode struct {
	ID     string `json:"id"`
	LinkID string `json:"link_id"`
	Pow    int    `json:"power"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
}

func createStorableNode(ins Node) *storableNode {
	out := storableNode{
		ID:     ins.ID().String(),
		LinkID: ins.Link().ID().String(),
		Pow:    ins.Power(),
		IP:     ins.IP().String(),
		Port:   ins.Port(),
	}

	return &out
}
