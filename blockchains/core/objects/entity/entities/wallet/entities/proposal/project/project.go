package project

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/wallet"
	approved_project "github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/entities/project"
)

type project struct {
	UUID      *uuid.UUID               `json:"id"`
	Proj      approved_project.Project `json:"project"`
	Own       wallet.Wallet            `json:"owner"`
	Mgr       wallet.Wallet            `json:"manager"`
	MgrShares int                      `json:"manager_shares"`
	Lnk       wallet.Wallet            `json:"linker"`
	LnkShares int                      `json:"linker_shares"`
	WrkShares int                      `json:"worker_shares"`
}

func createProject(
	id *uuid.UUID,
	proj approved_project.Project,
	own wallet.Wallet,
	mgr wallet.Wallet,
	mgrShares int,
	lnk wallet.Wallet,
	lnkShares int,
	wrkShares int,
) (Project, error) {
	out := project{
		UUID:      id,
		Proj:      proj,
		Own:       own,
		Mgr:       mgr,
		MgrShares: mgrShares,
		Lnk:       lnk,
		LnkShares: lnkShares,
		WrkShares: wrkShares,
	}

	return &out, nil
}

func createProjectFromNormalized(normalized *normalizedProject) (Project, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	projIns, projInsErr := approved_project.SDKFunc.CreateMetaData().Denormalize()(normalized.Project)
	if projInsErr != nil {
		return nil, projInsErr
	}

	ownerIns, ownerInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(normalized.Owner)
	if ownerInsErr != nil {
		return nil, ownerInsErr
	}

	mgrIns, mgrInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(normalized.Mgr)
	if mgrInsErr != nil {
		return nil, mgrInsErr
	}

	linkerIns, linkerInsErr := wallet.SDKFunc.CreateMetaData().Denormalize()(normalized.Lnk)
	if linkerInsErr != nil {
		return nil, linkerInsErr
	}

	if proj, ok := projIns.(approved_project.Project); ok {
		if own, ok := ownerIns.(wallet.Wallet); ok {
			if mgr, ok := mgrIns.(wallet.Wallet); ok {
				if lnker, ok := linkerIns.(wallet.Wallet); ok {
					return createProject(&id, proj, own, mgr, normalized.MgrShares, lnker, normalized.LnkShares, normalized.WrkShares)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", linkerIns.ID().String())
				return nil, errors.New(str)
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", mgrIns.ID().String())
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", ownerIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", projIns.ID().String())
	return nil, errors.New(str)
}

func createProjectFromStorable(storable *storableProject, rep entity.Repository) (Project, error) {
	// create metadata:
	walletMetaData := wallet.SDKFunc.CreateMetaData()

	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	projID, projIDErr := uuid.FromString(storable.ProjectID)
	if projIDErr != nil {
		return nil, projIDErr
	}

	ownerID, ownerIDErr := uuid.FromString(storable.OwnerID)
	if ownerIDErr != nil {
		return nil, ownerIDErr
	}

	mgrID, mgrIDErr := uuid.FromString(storable.MgrID)
	if mgrIDErr != nil {
		return nil, mgrIDErr
	}

	lnkerID, lnkerIDErr := uuid.FromString(storable.LnkID)
	if lnkerIDErr != nil {
		return nil, lnkerIDErr
	}

	projIns, projInsErr := rep.RetrieveByID(approved_project.SDKFunc.CreateMetaData(), &projID)
	if projInsErr != nil {
		return nil, projInsErr
	}

	ownerIns, ownerInsErr := rep.RetrieveByID(walletMetaData, &ownerID)
	if ownerInsErr != nil {
		return nil, ownerInsErr
	}

	mgrIns, mgrInsErr := rep.RetrieveByID(walletMetaData, &mgrID)
	if mgrInsErr != nil {
		return nil, mgrInsErr
	}

	linkerIns, linkerInsErr := rep.RetrieveByID(walletMetaData, &lnkerID)
	if linkerInsErr != nil {
		return nil, linkerInsErr
	}

	if proj, ok := projIns.(approved_project.Project); ok {
		if own, ok := ownerIns.(wallet.Wallet); ok {
			if mgr, ok := mgrIns.(wallet.Wallet); ok {
				if lnker, ok := linkerIns.(wallet.Wallet); ok {
					return createProject(&id, proj, own, mgr, storable.MgrShares, lnker, storable.LnkShares, storable.WrkShares)
				}

				str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", linkerIns.ID().String())
				return nil, errors.New(str)
			}

			str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", mgrIns.ID().String())
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Wallet instance", ownerIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Project instance", projIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *project) ID() *uuid.UUID {
	return obj.UUID
}

// Project returns the project
func (obj *project) Project() approved_project.Project {
	return obj.Proj
}

// Owner returns the owner wallet
func (obj *project) Owner() wallet.Wallet {
	return obj.Own
}

// Manager returns the manager wallet
func (obj *project) Manager() wallet.Wallet {
	return obj.Mgr
}

// ManagerShares returns the manager shares
func (obj *project) ManagerShares() int {
	return obj.MgrShares
}

// Linker returns the linker wallet
func (obj *project) Linker() wallet.Wallet {
	return obj.Lnk
}

// LinkerShares returns the linker shares
func (obj *project) LinkerShares() int {
	return obj.LnkShares
}

// WorkerShares returns the worker shares
func (obj *project) WorkerShares() int {
	return obj.WrkShares
}
