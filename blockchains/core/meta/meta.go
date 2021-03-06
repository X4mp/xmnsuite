package meta

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
)

type meta struct {
	gen                            entity.Representation
	wal                            entity.Representation
	req                            entity.Representation
	vot                            entity.Representation
	ret                            map[string]entity.MetaData
	wr                             map[string]entity.Representation
	allWriteOnEntReqRepresentation map[string]entity.Representation
	wrOnEntReq                     map[string]EntityRequest
}

func createMeta(
	gen entity.Representation,
	wal entity.Representation,
	req entity.Representation,
	vot entity.Representation,
	ret map[string]entity.MetaData,
	wr map[string]entity.Representation,
	wrOnEntReq map[string]EntityRequest,
) Meta {

	allWriteOnEntReqRepresentation := map[string]entity.Representation{}
	for _, oneReq := range wrOnEntReq {
		representations := oneReq.Map()
		for _, oneRepresentation := range representations {
			keyname := oneRepresentation.MetaData().Keyname()
			if _, ok := allWriteOnEntReqRepresentation[keyname]; !ok {
				allWriteOnEntReqRepresentation[keyname] = oneRepresentation
			}

			// register:
			request.SDKFunc.Register(request.RegisterParams{
				EntityMetaData: oneRepresentation.MetaData(),
			})
		}
	}

	out := meta{
		gen: gen,
		wal: wal,
		req: req,
		vot: vot,
		ret: ret,
		wr:  wr,
		allWriteOnEntReqRepresentation: allWriteOnEntReqRepresentation,
		wrOnEntReq:                     wrOnEntReq,
	}

	return &out
}

// Genesis returns the genesis representation
func (obj *meta) Genesis() entity.Representation {
	return obj.gen
}

// Wallet returns the wallet representation
func (obj *meta) Wallet() entity.Representation {
	return obj.wal
}

// Request returns the request representation
func (obj *meta) Request() entity.Representation {
	return obj.req
}

// Vote returns the vote representation
func (obj *meta) Vote() entity.Representation {
	return obj.vot
}

// Retrieval returns the retrieval metadata
func (obj *meta) Retrieval() map[string]entity.MetaData {
	return obj.ret
}

// Write returns the write representation
func (obj *meta) Write() map[string]entity.Representation {
	return obj.wr
}

// WriteOnAllEntityRequest returns all the write on entity representation
func (obj *meta) WriteOnAllEntityRequest() map[string]entity.Representation {
	return obj.allWriteOnEntReqRepresentation
}

// WriteOnEntityRequest returns the write on entity request representation
func (obj *meta) WriteOnEntityRequest() map[string]EntityRequest {
	return obj.wrOnEntReq
}

// AddToWriteOnEntityRequest adds a new entity that can be voted by the requestedBy entity
func (obj *meta) AddToWriteOnEntityRequest(requestedBy entity.MetaData, rep entity.Representation) error {
	keyname := requestedBy.Keyname()
	if _, ok := obj.wrOnEntReq[keyname]; !ok {
		str := fmt.Sprintf("the requestedBy entity (Keyname: %s) is not a valid entity requester", keyname)
		return errors.New(str)
	}

	// add to the list:
	metKeyname := rep.MetaData().Keyname()
	obj.wrOnEntReq[keyname].Add(rep)
	if _, ok := obj.allWriteOnEntReqRepresentation[metKeyname]; !ok {
		obj.allWriteOnEntReqRepresentation[metKeyname] = rep
	}

	// register:
	request.SDKFunc.Register(request.RegisterParams{
		EntityMetaData: rep.MetaData(),
	})

	return nil
}
