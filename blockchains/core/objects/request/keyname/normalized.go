package keyname

import "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"

type normalizedKeyname struct {
	ID    string           `json:"id"`
	Group group.Normalized `json:"group"`
	Name  string           `json:"name"`
}

func createNormalizedKeyname(ins Keyname) (*normalizedKeyname, error) {
	grp, grpErr := group.SDKFunc.CreateMetaData().Normalize()(ins.Group())
	if grpErr != nil {
		return nil, grpErr
	}

	out := normalizedKeyname{
		ID:    ins.ID().String(),
		Group: grp,
		Name:  ins.Name(),
	}

	return &out, nil
}
