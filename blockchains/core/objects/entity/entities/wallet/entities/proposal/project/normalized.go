package project

import (
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	approved_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
)

type normalizedProject struct {
	ID        string                      `json:"id"`
	Project   approved_project.Normalized `json:"project"`
	Owner     wallet.Normalized           `json:"owner"`
	Mgr       wallet.Normalized           `json:"manager"`
	MgrShares int                         `json:"manager_shares"`
	Lnk       wallet.Normalized           `json:"linker"`
	LnkShares int                         `json:"linker_shares"`
	WrkShares int                         `json:"worker_shares"`
}

func createNormalizedProject(ins Project) (*normalizedProject, error) {
	proj, projErr := approved_project.SDKFunc.CreateMetaData().Normalize()(ins.Project())
	if projErr != nil {
		return nil, projErr
	}

	owner, ownerErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.Owner())
	if ownerErr != nil {
		return nil, ownerErr
	}

	mgr, mgrErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.Manager())
	if mgrErr != nil {
		return nil, mgrErr
	}

	linker, linkerErr := wallet.SDKFunc.CreateMetaData().Normalize()(ins.Linker())
	if linkerErr != nil {
		return nil, linkerErr
	}

	out := normalizedProject{
		ID:        ins.ID().String(),
		Project:   proj,
		Owner:     owner,
		Mgr:       mgr,
		MgrShares: ins.ManagerShares(),
		Lnk:       linker,
		LnkShares: ins.LinkerShares(),
		WrkShares: ins.WorkerShares(),
	}

	return &out, nil
}
