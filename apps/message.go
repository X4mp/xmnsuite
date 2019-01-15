package main

import (
	"encoding/json"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	genesis "github.com/xmnservices/xmnsuite/apps/messages/genesis"
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "genesis.retrieve":
		payload, err = genesis.SDKFunc.Retrieve(m.Payload)
		break
	case "event.name":
		// Unmarshal payload
		var s string
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
			return
		}
		payload = s + " world"
	}
	return
}
