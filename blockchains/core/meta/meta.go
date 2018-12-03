package meta

import "github.com/xmnservices/xmnsuite/blockchains/core/entity"

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
			allWriteOnEntReqRepresentation[keyname] = oneRepresentation
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
