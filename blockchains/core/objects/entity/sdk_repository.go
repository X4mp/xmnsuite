package entity

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

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
	queryResp, queryRespErr := app.execute(queryPath)
	if queryRespErr != nil {
		return nil, queryRespErr
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
	// create the comma-separated list:
	keynamesList := strings.Join(keynames, ",")

	// encode the keynames:
	encodedKeynames := base64.StdEncoding.EncodeToString([]byte(keynamesList))

	// create the resource pointer:
	queryPath := fmt.Sprintf("%s/%s/%s/intersect", app.routePrefix, met.Keyname(), encodedKeynames)
	queryResp, queryRespErr := app.execute(queryPath)
	if queryRespErr != nil {
		return nil, queryRespErr
	}

	// convert to an entity:
	ins, insErr := met.ToEntity()(app, queryResp.Value())
	if insErr != nil {
		return nil, insErr
	}

	return ins, nil
}

// RetrieveSetByKeyname retrieves an entity set by using a keyname
func (app *sdkRepository) RetrieveSetByKeyname(met MetaData, keyname string, index int, amount int) (PartialSet, error) {
	return app.RetrieveSetByIntersectKeynames(met, []string{keyname}, index, amount)
}

// RetrieveSetByIntersectKeynames retrieves an entity set by intersecting keynames
func (app *sdkRepository) RetrieveSetByIntersectKeynames(met MetaData, keynames []string, index int, amount int) (PartialSet, error) {
	// create the comma-separated list:
	keynamesList := strings.Join(keynames, ",")

	// encode the keynames:
	encodedKeynames := base64.StdEncoding.EncodeToString([]byte(keynamesList))

	// create the resource pointer:
	queryPath := fmt.Sprintf("%s/%s/%s/set/intersect", app.routePrefix, met.Keyname(), encodedKeynames)
	queryResp, queryRespErr := app.execute(queryPath)
	if queryRespErr != nil {
		return nil, queryRespErr
	}

	// unmarshal the normalized partial set:
	ptr := new(normalizedPartialSet)
	jsErr := cdc.UnmarshalJSON(queryResp.Value(), ptr)
	if jsErr != nil {
		return nil, jsErr
	}

	// denormalize:
	ps, psErr := createEntityPartialSetFromNormalized(ptr, met)
	if psErr != nil {
		return nil, psErr
	}

	return ps, nil
}

func (app *sdkRepository) execute(path string) (routers.QueryResponse, error) {
	queryResPtr := routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
		From: app.pk.PublicKey(),
		Path: path,
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

	return queryResp, nil
}
