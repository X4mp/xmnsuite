package link

type storableLink struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	NodeIDs     []string `json:"node_ids"`
}

func createStorableLink(ins Link) *storableLink {

	nodes := ins.Nodes()
	nodeIDs := []string{}
	for _, oneNode := range nodes {
		nodeIDs = append(nodeIDs, oneNode.ID().String())
	}

	out := storableLink{
		ID:          ins.ID().String(),
		Title:       ins.Title(),
		Description: ins.Description(),
		NodeIDs:     nodeIDs,
	}

	return &out
}
