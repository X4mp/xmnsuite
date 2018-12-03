package entity

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/routers"
)

type sdkRepository struct {
	pk          crypto.PrivateKey
	client      applications.Client
	routePrefix string
}

func createSDKRepository(pk crypto.PrivateKey, client applications.Client, routePrefix string) Repository {
	out := sdkRepository{
		pk:          pk,
		client:      client,
		routePrefix: routePrefix,
	}
	return &out
}

// RetrieveByID retrieves an entity by its ID
func (app *sdkRepository) RetrieveByID(met MetaData, id *uuid.UUID) (Entity, error) {
	// create the resource pointer:
	queryPath := fmt.Sprintf("%s/%s/%s", app.routePrefix, met.Keyname(), id.String())
	queryResPtr := routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
		From: app.pk.PublicKey(),
		Path: queryPath,
	})

	// create the signature:
	querySig := app.pk.Sign(queryResPtr.Hash())

	// execute a query:
	queryResp, queryRespErr := app.client.Query(routers.SDKFunc.CreateQueryRequest(routers.CreateQueryRequestParams{
		Ptr: queryResPtr,
		Sig: querySig,
	}))

	if queryRespErr != nil {
		return nil, queryRespErr
	}

	if queryResp.Code() != routers.IsSuccessful {
		str := fmt.Sprintf("there was an error (Code: %d) while executing the query: %s", queryResp.Code(), queryResp.Log())
		return nil, errors.New(str)
	}

	// convert to an entity:
	ins, insErr := met.ToEntity()(app, queryResp.Value())
	if insErr != nil {
		return nil, insErr
	}

	return ins, nil
}

// RetrieveByIntersectKeynames retrieves an entity by intersecting keynames
func (app *sdkRepository) RetrieveByIntersectKeynames(met MetaData, keynames []string) (Entity, error) {
	return nil, nil
}

// RetrieveSetByKeyname retrieves an entity set by using a keyname
func (app *sdkRepository) RetrieveSetByKeyname(met MetaData, keyname string, index int, amount int) (PartialSet, error) {
	return nil, nil
}

// RetrieveSetByIntersectKeynames retrieves an entity set by intersecting keynames
func (app *sdkRepository) RetrieveSetByIntersectKeynames(met MetaData, keynames []string, index int, amount int) (PartialSet, error) {
	return nil, nil
}
