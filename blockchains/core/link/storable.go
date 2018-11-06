package link

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

type storableLink struct {
	ID          string          `json:"id"`
	Keyname     string          `json:"keyname"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Nodes       []*storableNode `json:"nodes"`
}

func createStorableLink(ins Link) *storableLink {

	nodes := ins.Nodes()
	storableNodes := []*storableNode{}
	for _, oneNode := range nodes {
		storableNodes = append(storableNodes, createStorableNode(oneNode))
	}

	out := storableLink{
		ID:          ins.ID().String(),
		Keyname:     ins.Keyname(),
		Title:       ins.Title(),
		Description: ins.Description(),
		Nodes:       storableNodes,
	}

	return &out
}
