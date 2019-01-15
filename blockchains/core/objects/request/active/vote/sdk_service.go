package vote

import (
	"errors"
	"fmt"

	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/crypto"
	"github.com/xmnservices/xmnsuite/routers"
)

type outgoingVote struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Reason     string `json:"reason"`
	IsNeutral  bool   `json:"is_neutral"`
	IsApproved bool   `json:"is_approved"`
}

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
func (app *sdkService) Save(ins Vote, rep entity.Representation) error {
	// make sure the voter matches the pk:
	if !ins.Voter().PubKey().Equals(app.pk.PublicKey()) {
		str := fmt.Sprintf("the Voter PubKey was not created by the service's PK")
		return errors.New(str)
	}

	// create the vote:
	outVote := outgoingVote{
		ID:         ins.ID().String(),
		UserID:     ins.Voter().ID().String(),
		Reason:     ins.Reason(),
		IsNeutral:  ins.IsNeutral(),
		IsApproved: ins.IsApproved(),
	}

	// marshals to JSON:
	js, jsErr := cdc.MarshalJSON(&outVote)
	if jsErr != nil {
		return jsErr
	}

	// create the resource:
	route := fmt.Sprintf("/%s/requests/%s", rep.MetaData().Keyname(), ins.Request().ID().String())
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
