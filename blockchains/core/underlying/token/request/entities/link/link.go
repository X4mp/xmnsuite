package link

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/token/request/entities/node"
)

type link struct {
	UUID *uuid.UUID  `json:"id"`
	Key  string      `json:"keyname"`
	Titl string      `json:"title"`
	Desc string      `json:"description"`
	Nods []node.Node `json:"nodes"`
}

func createLink(id *uuid.UUID, keyname string, title string, description string, nodes []node.Node) (Link, error) {

	if len(nodes) <= 0 {
		return nil, errors.New("the link must contain at least 1 Node")
	}

	out := link{
		UUID: id,
		Key:  keyname,
		Titl: title,
		Desc: description,
		Nods: nodes,
	}

	return &out, nil
}

func createLinkFromNormalized(normalized *normalizedLink) (Link, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	nodes := []node.Node{}
	nodeMetaData := node.SDKFunc.CreateMetaData()
	for _, oneNormalizedNode := range normalized.Nodes {
		oneNode, oneNodeErr := nodeMetaData.Denormalize()(oneNormalizedNode)
		if oneNodeErr != nil {
			return nil, oneNodeErr
		}

		if nod, ok := oneNode.(node.Node); ok {
			nodes = append(nodes, nod)
			continue
		}

		str := fmt.Sprintf("there is at least one entity (ID: %s) that was expected to be a node in the link (ID: %s), but is not", oneNode.ID().String(), id.String())
		return nil, errors.New(str)

	}

	return createLink(&id, normalized.Keyname, normalized.Title, normalized.Description, nodes)
}

// ID returns the ID
func (obj *link) ID() *uuid.UUID {
	return obj.UUID
}

// Keyname returns the keyname
func (obj *link) Keyname() string {
	return obj.Key
}

// Title returns the title
func (obj *link) Title() string {
	return obj.Titl
}

// Description returns the description
func (obj *link) Description() string {
	return obj.Desc
}

// Nodes returns the nodes
func (obj *link) Nodes() []node.Node {
	return obj.Nods
}
