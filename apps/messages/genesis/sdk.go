package genesis

import "encoding/json"

// SDKFunc represents the genesis SDK func
var SDKFunc = struct {
	Retrieve func(input json.RawMessage) (interface{}, error)
}{
	Retrieve: func(input json.RawMessage) (interface{}, error) {
		return retrieve(input)
	},
}
