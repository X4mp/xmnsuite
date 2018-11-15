package request

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/entity"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/routers"
)

type sdkService struct {
	pk     crypto.PrivateKey
	client applications.Client
}

func createSDKService(pk crypto.PrivateKey, client applications.Client) Service {
	out := sdkService{
		pk:     pk,
		client: client,
	}
	return &out
}

// Save saves a request instance to the service
func (app *sdkService) Save(req Request, rep entity.Representation) error {
	normalized, normalizedErr := rep.MetaData().Normalize()(req.New())
	if normalizedErr != nil {
		return normalizedErr
	}

	js, jsErr := cdc.MarshalJSON(normalized)
	if jsErr != nil {
		return jsErr
	}

	// create the resource:
	route := fmt.Sprintf("/requests/%s/%s/%s", req.From().Wallet().ID().String(), rep.MetaData().Keyname(), req.ID().String())
	firstRes := routers.SDKFunc.CreateResource(routers.CreateResourceParams{
		ResPtr: routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
			From: app.pk.PublicKey(),
			Path: route,
		}),
		Data: js,
	})

	// sign the resource:
	firstSig := app.pk.Sign(firstRes.Hash())

	// save the instance:
	trxResp, trxRespErr := app.client.Transact(routers.SDKFunc.CreateTransactionRequest(routers.CreateTransactionRequestParams{
		Res: firstRes,
		Sig: firstSig,
	}))

	if trxRespErr != nil {
		return trxRespErr
	}

	chk := trxResp.Check()
	chkCode := chk.Code()
	if chkCode != routers.IsSuccessful {
		trx := trxResp.Transaction()
		str := fmt.Sprintf("there was an error (Check Code: %d, Trx Code: %d) while executing the transaction: (Check Log: %s, Trx Log: %s, Route: %s)", chkCode, trx.Code(), chk.Log(), trx.Log(), route)
		return errors.New(str)
	}

	return nil
}
