package request

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/routers"
)

type outgoingRequest struct {
	ID               string `json:"id"`
	Reason           string `json:"reason"`
	WalletID         string `json:"wallet_id"`
	SaveEntityJSON   []byte `json:"save_entity_json"`
	DeleteEntityJSON []byte `json:"delete_entity_json"`
}

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

// Save saves a request instance to the service
func (app *sdkService) Save(req Request, rep entity.Representation) error {

	var toSaveEntity []byte
	var toDeleteEntity []byte

	if req.HasSave() {
		normalized, normalizedErr := rep.MetaData().Normalize()(req.Save())
		if normalizedErr != nil {
			return normalizedErr
		}

		insJS, insJSErr := cdc.MarshalJSON(normalized)
		if insJSErr != nil {
			return insJSErr
		}

		toSaveEntity = insJS
	}

	if req.HasDelete() {
		normalized, normalizedErr := rep.MetaData().Normalize()(req.Delete())
		if normalizedErr != nil {
			return normalizedErr
		}

		insJS, insJSErr := cdc.MarshalJSON(normalized)
		if insJSErr != nil {
			return insJSErr
		}

		toDeleteEntity = insJS
	}

	outReq := outgoingRequest{
		ID:               req.ID().String(),
		Reason:           req.Reason(),
		WalletID:         req.From().Wallet().ID().String(),
		SaveEntityJSON:   toSaveEntity,
		DeleteEntityJSON: toDeleteEntity,
	}

	js, jsErr := cdc.MarshalJSON(&outReq)
	if jsErr != nil {
		return jsErr
	}

	// create the resource:
	route := fmt.Sprintf("%s/%s/requests", app.routePrefix, req.Keyname().Name())
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
