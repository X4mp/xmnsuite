package web

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/address"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/deposit"
	"github.com/xmnservices/xmnsuite/applications/cryptocurrency/objects/offer"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/transfer"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	active_request "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote"
	active_vote "github.com/xmnservices/xmnsuite/blockchains/core/objects/request/active/vote/active"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/group"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/keyname"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
	core_web "github.com/xmnservices/xmnsuite/blockchains/core/web"
	"github.com/xmnservices/xmnsuite/crypto"
)

const (
	amountOfElementsPerList = 20
	maxLoginFormSizeInBytes = 1000 * 1000
	newWalletAmountShares   = 100
	loginCookieName         = "login"
)

type web struct {
	rter           *mux.Router
	port           int
	templateDir    string
	staticFilesDir string
	meta           meta.Meta
	client         applications.Client
}

func createWeb(
	port int,
	meta meta.Meta,
	client applications.Client,
	pk crypto.PrivateKey,
) Web {

	templateDir := "./applications/cryptocurrency/web/templates"
	r := mux.NewRouter()

	app := web{
		port:           port,
		templateDir:    templateDir,
		staticFilesDir: "./applications/cryptocurrency/web/static",
		meta:           meta,
		client:         client,
		rter:           r,
	}

	entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
		PK:          pk,
		Client:      client,
		RoutePrefix: "",
	})

	app.rter.HandleFunc("/", app.home)

	app.rter.HandleFunc("/register", core_web.SDKFunc.RouteRegister(core_web.RouteRegisterParams{
		RegisteredTemplate:    app.createTemplate("registered.html", "registered"),
		RegisterTemplate:      app.createTemplate("register.html", "register"),
		NewWalletAmountShares: newWalletAmountShares,
		Codec:  cdc,
		PK:     pk,
		Client: client,
	}))

	app.rter.HandleFunc("/address", address.SDKFunc.RouteSet(address.RouteSetParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("specific/addresses.html", "addresses"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/address/new", address.SDKFunc.RouteNew(address.RouteNewParams{
		PK:                      pk,
		Client:                  client,
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("specific/address_new.html", "address_new"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/offers", offer.SDKFunc.RouteSet(offer.RouteSetParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("specific/offers.html", "offers"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/offers/new", offer.SDKFunc.RouteNew(offer.RouteNewParams{
		PK:                      pk,
		Client:                  client,
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("specific/offers_new.html", "offers_new"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/offers/{id}", deposit.SDKFunc.RouteSet(deposit.RouteSetParams{
		PK:               pk,
		Client:           client,
		Tmpl:             app.createTemplate("specific/offers_single.html", "offers"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/requests", group.SDKFunc.RouteSet(group.RouteSetParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("core/request_groups.html", "request_groups"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/requests/{name}", keyname.SDKFunc.RouteSetOfGroup(keyname.RouteSetOfGroupParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("core/request_groups_keynames.html", "request_groups_keynames"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/requests/{groupname}/{keyname}", active_request.SDKFunc.RouteSetOfKeyname(active_request.RouteSetOfKeynameParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("core/requests.html", "requests"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/requests/{groupname}/{keyname}/{id}", active_vote.SDKFunc.RouteSetOfRequest(active_vote.RouteSetOfRequestParams{
		PK: pk,
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("core/requests_single.html", "requests_single"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/requests/{groupname}/{keyname}/{id}/vote", vote.SDKFunc.RouteNew(vote.RouteNewParams{
		PK:               pk,
		Client:           client,
		EntityRepository: entityRepository,
		FetchRepresentation: func(groupName string, keyname string) (entity.Representation, error) {
			reqs := app.meta.WriteOnEntityRequest()
			if req, ok := reqs[groupName]; ok {
				mps := req.Map()
				if rep, ok := mps[keyname]; ok {
					return rep, nil
				}

				str := fmt.Sprintf("the given keyname (%s) cannot be voted on by the group (%s)", keyname, groupName)
				return nil, errors.New(str)
			}

			str := fmt.Sprintf("the given group (%s) cannot vote", groupName)
			return nil, errors.New(str)
		},
	}))

	app.rter.HandleFunc("/genesis", genesis.SDKFunc.Route(genesis.RouteParams{
		Tmpl:             app.createTemplate("core/genesis.html", "genesis"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/wallets", balance.SDKFunc.RouteList(balance.RouteListParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("core/wallets.html", "wallets"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/wallets/{id}", balance.SDKFunc.Route(balance.RouteParams{
		Tmpl:             app.createTemplate("core/wallet_single.html", "wallet_single"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/users", user.SDKFunc.RouteWalletList(user.RouteWalletListParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("core/users.html", "users"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/users/{wallet_id}", user.SDKFunc.RouteUserSetInWallet(user.RouteUserSetInWalletParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("core/users_of_wallet.html", "users_of_wallet"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/users/{wallet_id}/{pubkey}", user.SDKFunc.RouteUserInWallet(user.RouteUserInWalletParams{
		Tmpl:             app.createTemplate("core/user_single.html", "user_single"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/transfers", transfer.SDKFunc.RouteSet(transfer.RouteSetParams{
		AmountOfElementsPerList: amountOfElementsPerList,
		Tmpl:             app.createTemplate("core/transfers.html", "transfers"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/transfers/new", transfer.SDKFunc.RouteNew(transfer.RouteNewParams{
		PK:               pk,
		Client:           app.client,
		Tmpl:             app.createTemplate("core/transfers_new.html", "transfers_new"),
		EntityRepository: entityRepository,
	}))

	app.rter.HandleFunc("/transfers/{id}", transfer.SDKFunc.Route(transfer.RouteParams{
		Tmpl:             app.createTemplate("core/transfer_single.html", "transfer_single"),
		EntityRepository: entityRepository,
	}))

	// add the login middleware:
	app.rter.Use(func(next http.Handler) http.Handler {
		return core_web.SDKFunc.CreateMiddleware(core_web.CreateMiddlewareParams{
			Next:                    next,
			CookieName:              loginCookieName,
			MaxLoginFormSizeInBytes: maxLoginFormSizeInBytes,
			LoginTemplate:           app.createTemplate("login.html", "login"),
			LoggedInTemplate:        app.createTemplate("loggedin.html", "loggedin"),
		})
	})

	// setup the static files:
	app.rter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(app.staticFilesDir))))

	http.Handle("/", app.rter)
	return &app
}

// Start starts the web server
func (app *web) Start() error {
	addr := fmt.Sprintf(":%d", app.port)
	srv := &http.Server{
		Addr: addr,
		// Avoid Slowloris attacks...
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      app.rter,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

// Stop stops the web server
func (app *web) Stop() error {
	return nil
}

func (app *web) createTemplate(fileName string, name string) *template.Template {
	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, fileName))
	if contentErr != nil {
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		panic(errors.New(str))
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
		panic(errors.New(str))
	}

	return tmpl
}

func (app *web) home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	str := fmt.Sprintf("home")
	w.Write([]byte(str))
	return
}
