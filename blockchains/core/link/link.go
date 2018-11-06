package link

import uuid "github.com/satori/go.uuid"

type link struct {
	UUID *uuid.UUID `json:"id"`
	Key  string     `json:"keyname"`
	Titl string     `json:"title"`
	Desc string     `json:"description"`
	Nods []Node
}

func createLink(id *uuid.UUID, keyname string, title string, description string) Link {
	out := link{
		UUID: id,
		Key:  keyname,
		Titl: title,
		Desc: description,
		Nods: []Node{},
	}

	return &out
}
func createLinkWithNodes(id *uuid.UUID, keyname string, title string, description string, nodes []Node) Link {
	out := link{
		UUID: id,
		Key:  keyname,
		Titl: title,
		Desc: description,
		Nods: nodes,
	}

	return &out
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
func (obj *link) Nodes() []Node {
	return obj.Nods
}
