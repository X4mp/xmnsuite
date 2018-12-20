package account

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/routers"
)

type sdkService struct {
	pk          crypto.PrivateKey
	client      applications.Client
	routePrefix string
}

func createSDKService(pk crypto.PrivateKey, client applications.Client, routePrefix string) Service {
	out := sdkService{
		pk:          pk,
		client:      client,
		routePrefix: routePrefix,
	}

	return &out
}

// Save saves a new account
func (app *sdkService) Save(ins Account, amountOfWorkToVerify int) error {
	normalized, normalizedErr := createNormalizedAccount(ins)
	if normalizedErr != nil {
		return normalizedErr
	}

	js, jsErr := cdc.MarshalJSON(normalized)
	if jsErr != nil {
		return jsErr
	}

	// create the resource:
	route := fmt.Sprintf("%s/account", app.routePrefix)
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