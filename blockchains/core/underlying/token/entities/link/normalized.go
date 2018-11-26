package link

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/entities/node"
)

type normalizedLink struct {
	ID          string            `json:"id"`
	Keyname     string            `json:"keyname"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Nodes       []node.Normalized `json:"nodes"`
}

func createNormalizedLink(link Link) (*normalizedLink, error) {
	nodes := link.Nodes()
	normalizedNodes := []node.Normalized{}
	nodeMetaData := node.SDKFunc.CreateMetaData()
	for _, oneNode := range nodes {
		oneNormalized, oneNormalizedErr := nodeMetaData.Normalize()(oneNode)
		if oneNormalizedErr != nil {
			return nil, oneNormalizedErr
		}

		normalizedNodes = append(normalizedNodes, oneNormalized)
	}

	out := normalizedLink{
		ID:          link.ID().String(),
		Keyname:     link.Keyname(),
		Title:       link.Title(),
		Description: link.Description(),
		Nodes:       normalizedNodes,
	}

	return &out, nil
}
