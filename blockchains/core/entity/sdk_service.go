package entity

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
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

// Save saves an entity instance to the service
func (app *sdkService) Save(ins Entity, rep Representation) error {
	normalized, normalizedErr := rep.MetaData().Normalize()(ins)
	if normalizedErr != nil {
		return normalizedErr
	}

	js, jsErr := cdc.MarshalJSON(normalized)
	if jsErr != nil {
		return jsErr
	}

	// create the resource:
	firstRes := routers.SDKFunc.CreateResource(routers.CreateResourceParams{
		ResPtr: routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
			From: app.pk.PublicKey(),
			Path: fmt.Sprintf("/%s", rep.MetaData().Keyname()),
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
		str := fmt.Sprintf("there was an error (Check Code: %d, Trx Code: %d) while executing the transaction: (Check Log: %s, Trx Log: %s)", chkCode, trx.Code(), chk.Log(), trx.Log())
		return errors.New(str)
	}

	return nil
}

// Delete deletes the entity instance from the service
func (app *sdkService) Delete(ins Entity, rep Representation) error {
	// create the resource:
	respPtr := routers.SDKFunc.CreateResourcePointer(routers.CreateResourcePointerParams{
		From: app.pk.PublicKey(),
		Path: fmt.Sprintf("/%s/%s", rep.MetaData().Keyname(), ins.ID().String()),
	})

	// sign the resource:
	firstSig := app.pk.Sign(respPtr.Hash())

	// delete the instance:
	trxResp, trxRespErr := app.client.Transact(routers.SDKFunc.CreateTransactionRequest(routers.CreateTransactionRequestParams{
		Ptr: respPtr,
		Sig: firstSig,
	}))

	if trxRespErr != nil {
		return trxRespErr
	}

	chk := trxResp.Check()
	chkCode := chk.Code()
	if chkCode != routers.IsSuccessful {
		trx := trxResp.Transaction()
		str := fmt.Sprintf("there was an error (Check Code: %d, Trx Code: %d) while executing the transaction: (Check Log: %s, Trx Log: %s)", chkCode, trx.Code(), chk.Log(), trx.Log())
		return errors.New(str)
	}

	return nil
}
