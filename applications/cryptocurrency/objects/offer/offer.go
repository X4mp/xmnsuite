package offer

import (
	"errors"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/pledge"
)

type offer struct {
	UUID   *uuid.UUID      `json:"id"`
	Pldge  pledge.Pledge   `json:"pledge"`
	Addr   address.Address `json:"to_address"`
	Conf   int             `json:"confirmations"`
	Am     int             `json:"amount"`
	Prce   int             `json:"price"`
	IPAddr net.IP          `json:"ip_address"`
	Prt    int             `json:"port"`
}

func createOffer(id *uuid.UUID, pldge pledge.Pledge, to address.Address, conf int, amount int, price int, ip net.IP, port int) (Offer, error) {
	out := offer{
		UUID:   id,
		Pldge:  pldge,
		Addr:   to,
		Conf:   conf,
		Am:     amount,
		Prce:   price,
		IPAddr: ip,
		Prt:    port,
	}

	return &out, nil
}

func createOfferFromNormalized(normalized *normalizedOffer) (Offer, error) {
	id, idErr := uuid.FromString(normalized.ID)
	if idErr != nil {
		return nil, idErr
	}

	pldgeIns, pldgeInsErr := pledge.SDKFunc.CreateMetaData().Denormalize()(normalized.Pledge)
	if pldgeInsErr != nil {
		return nil, pldgeInsErr
	}

	toAddrIns, toAddrInsErr := address.SDKFunc.CreateMetaData().Denormalize()(normalized.To)
	if toAddrInsErr != nil {
		return nil, toAddrInsErr
	}

	if pldge, ok := pldgeIns.(pledge.Pledge); ok {
		if toAddr, ok := toAddrIns.(address.Address); ok {
			ip := net.ParseIP(normalized.IP)
			return createOffer(&id, pldge, toAddr, normalized.Confirmations, normalized.Amount, normalized.Price, ip, normalized.Port)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Address instance", toAddrIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pldgeIns.ID().String())
	return nil, errors.New(str)

}

func createOfferFromStorable(storable *storableOffer, rep entity.Repository) (Offer, error) {
	id, idErr := uuid.FromString(storable.ID)
	if idErr != nil {
		return nil, idErr
	}

	pldgeID, pldgeIDErr := uuid.FromString(storable.PledgeID)
	if pldgeIDErr != nil {
		return nil, pldgeIDErr
	}

	toAddrID, toAddrIDErr := uuid.FromString(storable.ToID)
	if toAddrIDErr != nil {
		return nil, toAddrIDErr
	}

	pldgeIns, pldgeInsErr := rep.RetrieveByID(pledge.SDKFunc.CreateMetaData(), &pldgeID)
	if pldgeInsErr != nil {
		return nil, pldgeInsErr
	}

	toAddrIns, toAddrInsErr := rep.RetrieveByID(address.SDKFunc.CreateMetaData(), &toAddrID)
	if toAddrInsErr != nil {
		return nil, toAddrInsErr
	}

	if pldge, ok := pldgeIns.(pledge.Pledge); ok {
		if toAddr, ok := toAddrIns.(address.Address); ok {
			ip := net.ParseIP(storable.IP)
			return createOffer(&id, pldge, toAddr, storable.Confirmations, storable.Amount, storable.Price, ip, storable.Port)
		}

		str := fmt.Sprintf("the entity (ID: %s) is not a valid Address instance", toAddrIns.ID().String())
		return nil, errors.New(str)
	}

	str := fmt.Sprintf("the entity (ID: %s) is not a valid Pledge instance", pldgeIns.ID().String())
	return nil, errors.New(str)
}

// ID returns the ID
func (obj *offer) ID() *uuid.UUID {
	return obj.UUID
}

// Pedge returns the pledge
func (obj *offer) Pledge() pledge.Pledge {
	return obj.Pldge
}

// To returns the to Address
func (obj *offer) To() address.Address {
	return obj.Addr
}

// Amount returns the amount
func (obj *offer) Amount() int {
	return obj.Am
}

// Confirmations returns the amount of confirmations needed
func (obj *offer) Confirmations() int {
	return obj.Conf
}

// Price returns the price
func (obj *offer) Price() int {
	return obj.Prce
}

// IP returns the ip address
func (obj *offer) IP() net.IP {
	return obj.IPAddr
}

// Port returns the port
func (obj *offer) Port() int {
	return obj.Prt
}
