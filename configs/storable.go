package configs

import "encoding/hex"

type storableConfigs struct {
	NodePK   string `json:"node_pk"`
	WalletPK string `json:"wallet_pk"`
}

func createStorableConfigs(ins Configs) *storableConfigs {
	out := storableConfigs{
		NodePK:   hex.EncodeToString(ins.NodePK().Bytes()),
		WalletPK: ins.WalletPK().String(),
	}

	return &out
}
