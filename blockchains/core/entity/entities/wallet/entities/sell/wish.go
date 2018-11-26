package sell

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/underlying/external"
)

type wish struct {
	Tok external.External `json:"external_token"`
	Am  int               `json:"amount"`
}

func createWish(tok external.External, amount int) Wish {
	out := wish{
		Tok: tok,
		Am:  amount,
	}

	return &out
}

// Token returns the external token
func (obj *wish) Token() external.External {
	return obj.Tok
}

// Amount returns the amount
func (obj *wish) Amount() int {
	return obj.Am
}
