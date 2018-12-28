package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"strings"
	"unsafe"

	amino "github.com/tendermint/go-amino"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/work"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/configs"
	"github.com/xmnservices/xmnsuite/crypto"
)

// CreateMiddlewareParams represents the CreateMiddleware params
type CreateMiddlewareParams struct {
	Next                    http.Handler
	CookieName              string
	MaxLoginFormSizeInBytes int64
	LoginTemplate           *template.Template
	LoggedInTemplate        *template.Template
}

// RouteRegisterParams represents the RouteRegister params
type RouteRegisterParams struct {
	RegisteredTemplate    *template.Template
	RegisterTemplate      *template.Template
	NewWalletAmountShares int
	Codec                 *amino.Codec
	PK                    crypto.PrivateKey
	Client                applications.Client
}

// SDKFunc represents the web SDK func
var SDKFunc = struct {
	CreateMiddleware func(params CreateMiddlewareParams) http.Handler
	RouteRegister    func(params RouteRegisterParams) func(w http.ResponseWriter, r *http.Request)
}{
	CreateMiddleware: func(params CreateMiddlewareParams) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// if the file is static, bypass the middleware:
			if strings.HasPrefix(r.RequestURI, "/static") {
				// call the next handler:
				params.Next.ServeHTTP(w, r)
				return
			}

			// if the requestURI is set to register:
			if strings.HasPrefix(r.RequestURI, "/register") {
				// call the next handler:
				params.Next.ServeHTTP(w, r)
				return
			}

			conf := getConfigsFromCookie(params.CookieName, r)
			if conf == nil {
				if parseFormErr := r.ParseMultipartForm(params.MaxLoginFormSizeInBytes); parseFormErr != nil {
					// render:
					w.WriteHeader(http.StatusOK)
					params.LoginTemplate.Execute(w, nil)
					return
				}

				pass := r.FormValue("pass")
				if pass != "" {
					// read the uploaded file:
					var encryptedCoinfigData bytes.Buffer
					xmnFile, _, xmnFileErr := r.FormFile("xmnfile")
					if xmnFileErr != nil {
						w.WriteHeader(http.StatusInternalServerError)
						str := fmt.Sprintf("there was an error while reading the uploaded file: %s", xmnFileErr.Error())
						w.Write([]byte(str))
						return
					}
					defer xmnFile.Close()
					io.Copy(&encryptedCoinfigData, xmnFile)

					// decrypt the configs:
					decrypted := configs.SDKFunc.Decrypt(configs.DecryptParams{
						Data: encryptedCoinfigData.String(),
						Pass: pass,
					})

					// set the cookie:
					http.SetCookie(w, &http.Cookie{
						Name:  params.CookieName,
						Value: decrypted.String(),
					})

					// render:
					w.WriteHeader(http.StatusOK)
					params.LoggedInTemplate.Execute(w, nil)
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("the given password is invalid")
				w.Write([]byte(str))
				return
			}

			// call the next handler:
			params.Next.ServeHTTP(w, r)
			return
		})
	},
	RouteRegister: func(params RouteRegisterParams) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			if parseFormErr := r.ParseForm(); parseFormErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("there was an error while parsing form elements: %s", parseFormErr.Error())
				w.Write([]byte(str))
				return
			}

			pass := r.FormValue("pass")
			rpass := r.FormValue("rpass")
			if pass != "" && rpass != "" {
				// generate the configs:
				conf := configs.SDKFunc.Generate()

				// create the repositories:
				entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
					PK:          params.PK,
					Client:      params.Client,
					RoutePrefix: "",
				})

				genesisRepository := genesis.SDKFunc.CreateRepository(genesis.CreateRepositoryParams{
					EntityRepository: entityRepository,
				})

				accountService := account.SDKFunc.CreateSDKService(account.CreateSDKServiceParams{
					PK:          params.PK,
					Client:      params.Client,
					RoutePrefix: "",
				})

				// retrieve the genesis:
				gen, genErr := genesisRepository.Retrieve()
				if genErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
					w.Write([]byte(str))
					return
				}

				// create the user:
				usr := user.SDKFunc.Create(user.CreateParams{
					PubKey: conf.WalletPK().PublicKey(),
					Shares: params.NewWalletAmountShares,
					Wallet: wallet.SDKFunc.Create(wallet.CreateParams{
						Creator:         conf.WalletPK().PublicKey(),
						ConcensusNeeded: params.NewWalletAmountShares,
					}),
				})

				// convert the user to json:
				jsUserData, jsUserErr := params.Codec.MarshalJSON(usr)
				if jsUserErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while converting a user to json: %s", jsUserErr.Error())
					w.Write([]byte(str))
					return
				}

				// calculate the gaz price:
				gazPrice := int(unsafe.Sizeof(jsUserData)) * gen.GazPriceInMatrixWorkKb()

				// create the account:
				ac := account.SDKFunc.Create(account.CreateAccountParams{
					User: usr,
					Work: work.SDKFunc.Generate(work.GenerateParams{
						MatrixSize:   gazPrice,
						MatrixAmount: gazPrice,
					}),
				})

				// save the account:
				saveErr := accountService.Save(ac, int(math.Ceil(float64(gazPrice/100))))
				if saveErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("there was an error while saving an Account instance: %s", saveErr.Error())
					w.Write([]byte(str))
					return
				}

				// encrypt the conf:
				encryptedConf := configs.SDKFunc.Encrypt(configs.EncryptParams{
					Conf:        conf,
					Pass:        pass,
					RetypedPass: rpass,
				})

				// render:
				w.WriteHeader(http.StatusOK)
				params.RegisteredTemplate.Execute(w, encryptedConf)
				return
			}

			// render:
			w.WriteHeader(http.StatusOK)
			params.RegisterTemplate.Execute(w, nil)
			return
		}
	},
}
