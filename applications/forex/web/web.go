package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/applications/forex/web/controllers/banks"
	"github.com/xmnservices/xmnsuite/blockchains/applications"
	"github.com/xmnservices/xmnsuite/blockchains/core/meta"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account"
	walletpkg "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/work"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/request/vote"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/underlying/token/balance"
	"github.com/xmnservices/xmnsuite/configs"
)

const (
	amountOfElementsPerList = 20
	maxLoginFormSizeInBytes = 1000 * 1000
	loginCookieName         = "login"
)

type web struct {
	rter                   *mux.Router
	port                   int
	templateDir            string
	staticFilesDir         string
	meta                   meta.Meta
	client                 applications.Client
	entityService          entity.Service
	accountService         account.Service
	requestService         request.Service
	requestRepository      request.Repository
	voteRepository         vote.Repository
	voteService            vote.Service
	userRepository         user.Repository
	balanceRepository      balance.Repository
	genesisRepository      genesis.Repository
	walletRepository       walletpkg.Repository
	categoryRepository     category.Repository
	currencyRepository     currency.Repository
	walletRepresentation   entity.Representation
	categoryRepresentation entity.Representation
}

func createWeb(
	port int,
	meta meta.Meta,
	client applications.Client,
	entityService entity.Service,
	accountService account.Service,
	userRepository user.Repository,
	balanceRepository balance.Repository,
	genesisRepository genesis.Repository,
	walletRepository walletpkg.Repository,
	categoryRepository category.Repository,
	currencyRepository currency.Repository,
) Web {

	templateDir := "./applications/forex/web/templates"
	r := mux.NewRouter()

	app := web{
		port:                   port,
		templateDir:            templateDir,
		staticFilesDir:         "./applications/forex/web/static",
		meta:                   meta,
		client:                 client,
		entityService:          entityService,
		accountService:         accountService,
		requestService:         nil,
		requestRepository:      nil,
		voteRepository:         nil,
		voteService:            nil,
		userRepository:         userRepository,
		balanceRepository:      balanceRepository,
		genesisRepository:      genesisRepository,
		walletRepository:       walletRepository,
		categoryRepository:     categoryRepository,
		currencyRepository:     currencyRepository,
		categoryRepresentation: category.SDKFunc.CreateRepresentation(),
		walletRepresentation:   walletpkg.SDKFunc.CreateRepresentation(),
		rter:                   r,
	}

	app.rter.HandleFunc("/", app.home)
	app.rter.HandleFunc("/register", app.register)
	app.rter.HandleFunc("/genesis", app.genesis)
	app.rter.HandleFunc("/users", app.users)
	app.rter.HandleFunc("/wallets", app.wallets)
	app.rter.HandleFunc("/wallets/{id}", app.walletSingle)
	app.rter.HandleFunc("/categories", app.categories)
	app.rter.HandleFunc("/categories/new", app.newCategoriesForm)
	app.rter.HandleFunc("/requests", app.requests)
	app.rter.HandleFunc("/requests/{id}", app.requestSingle)
	app.rter.HandleFunc("/requests/{id}/{action}", app.requestSingleVote)

	// bank controllers:
	banks.SDKFunc.ShowBanks(banks.ShowBankParams{
		Router:      app.rter,
		TemplateDir: templateDir,
	})

	banks.SDKFunc.NewBankForm(banks.NewBankFormParams{
		Router:             app.rter,
		TemplateDir:        templateDir,
		CurrencyRepository: currencyRepository,
	})

	// add the login middleware:
	app.rter.Use(app.middlewareVerifyConfigsInCookie)

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

func (app *web) middlewareVerifyConfigsInCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if the file is static, bypass the middleware:
		if strings.HasPrefix(r.RequestURI, "/static") {
			// call the next handler:
			next.ServeHTTP(w, r)
			return
		}

		// if the requestURI is set to register:
		if strings.HasPrefix(r.RequestURI, "/register") {
			// call the next handler:
			next.ServeHTTP(w, r)
			return
		}

		conf := getConfigsFromCookie(loginCookieName, r)
		if conf == nil {
			if parseFormErr := r.ParseMultipartForm(maxLoginFormSizeInBytes); parseFormErr != nil {
				// retrieve the html page:
				content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "login.html"))
				if contentErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
					w.Write([]byte(str))
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write(content)
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
					Name:  loginCookieName,
					Value: decrypted.String(),
				})

				// render the continue page:
				content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "loggedin.html"))
				if contentErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
					w.Write([]byte(str))
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write(content)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the given password is invalid")
			w.Write([]byte(str))
			return
		}

		// create the repository/services:
		app.requestService = request.SDKFunc.CreateSDKService(request.CreateSDKServiceParams{
			PK:          conf.WalletPK(),
			Client:      app.client,
			RoutePrefix: "",
		})

		entityRepository := entity.SDKFunc.CreateSDKRepository(entity.CreateSDKRepositoryParams{
			PK:          conf.WalletPK(),
			Client:      app.client,
			RoutePrefix: "",
		})

		app.requestRepository = request.SDKFunc.CreateRepository(request.CreateRepositoryParams{
			EntityRepository: entityRepository,
		})

		app.voteRepository = vote.SDKFunc.CreateRepository(vote.CreateRepositoryParams{
			EntityRepository: entityRepository,
		})

		app.voteService = vote.SDKFunc.CreateSDKService(vote.CreateSDKServiceParams{
			PK:     conf.WalletPK(),
			Client: app.client,
			CreateRouteFunc: func(ins vote.Vote, rep entity.Representation) (string, error) {
				keyname := rep.MetaData().Keyname()
				entRequests := app.meta.WriteOnEntityRequest()
				for _, oneReq := range entRequests {
					reqBy := oneReq.RequestedBy()
					mp := oneReq.Map()
					if _, ok := mp[keyname]; ok {
						return fmt.Sprintf("%s/requests/%s/%s", rep.MetaData().Keyname(), ins.Request().ID().String(), reqBy.MetaData().Keyname()), nil
					}
				}

				str := fmt.Sprintf("the keyname (Keyname: %s) cannot be voted on", keyname)
				return "", errors.New(str)
			},
		})

		// call the next handler:
		next.ServeHTTP(w, r)
		return
	})
}

func (app *web) home(w http.ResponseWriter, r *http.Request) {

	formatWalletPS := func(walPS entity.PartialSet, gen genesis.Genesis) *homeWalletList {
		walsIns := walPS.Instances()
		creatorWallets := []*homeWallet{}
		for _, oneWalletIns := range walsIns {
			if wal, ok := oneWalletIns.(walletpkg.Wallet); ok {
				// retrieve the wallet balance:
				bal, balErr := app.balanceRepository.RetrieveByWalletAndToken(wal, gen.Deposit().Token())
				if balErr != nil {
					log.Printf("there was an error while retrieving the wallet (ID: %s) balance of the given Token (ID: %s): %s", wal.ID().String(), gen.Deposit().Token().ID().String(), balErr.Error())
					continue
				}

				// retrieve the users:

				creatorWallets = append(creatorWallets, &homeWallet{
					ID:              wal.ID().String(),
					Creator:         wal.Creator().String(),
					ConcensusNeeded: wal.ConcensusNeeded(),
					TokenAmount:     bal.Amount(),
				})

				continue
			}

			log.Printf("the given entity (ID: %s) is not a valid Wallet instance", oneWalletIns.ID().String())
			continue
		}

		return &homeWalletList{
			Index:       walPS.Index(),
			Amount:      walPS.Amount(),
			TotalAmount: walPS.TotalAmount(),
			IsLast:      walPS.IsLast(),
			Wallets:     creatorWallets,
		}
	}

	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// retrieve the wallets created by our conf PK:
	walPS, walPSErr := app.walletRepository.RetrieveSetByCreatorPublicKey(conf.WalletPK().PublicKey(), 0, amountOfElementsPerList)
	if walPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the Wallet entity set from creator's public key (PubKey: %s): %s", conf.WalletPK().PublicKey().String(), walPSErr.Error())
		w.Write([]byte(str))
		return
	}

	// retrieve the users associated with our conf PK:
	usrPS, usrPSErr := app.userRepository.RetrieveSetByPubKey(conf.WalletPK().PublicKey(), 0, amountOfElementsPerList)
	if usrPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the users entity set from creator's public key (PubKey: %s): %s", conf.WalletPK().PublicKey().String(), usrPSErr.Error())
		w.Write([]byte(str))
		return
	}

	// retrieve all the wallets:
	allWalPS, allWalPSErr := app.walletRepository.RetrieveSet(0, amountOfElementsPerList)
	if allWalPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the wallet entity set: %s", allWalPSErr.Error())
		w.Write([]byte(str))
		return
	}

	// retrieve the genesis:
	gen, genErr := app.genesisRepository.Retrieve()
	if genErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
		w.Write([]byte(str))
		return
	}

	homeUsers := []*homeUser{}
	usersIns := usrPS.Instances()
	for _, oneUserIns := range usersIns {
		if usr, ok := oneUserIns.(user.User); ok {
			homeUsers = append(homeUsers, &homeUser{
				ID:       usr.ID().String(),
				Shares:   usr.Shares(),
				WalletID: usr.Wallet().ID().String(),
			})
		}

		log.Printf("the given entity (ID: %s) is not a valid User instance", oneUserIns.ID().String())
		continue
	}

	// template structure:
	templateData := home{
		WalletPS:    formatWalletPS(walPS, gen),
		AllWalletPS: formatWalletPS(allWalPS, gen),
		UserPS: &homeUserList{
			Index:       usrPS.Index(),
			Amount:      usrPS.Amount(),
			TotalAmount: usrPS.TotalAmount(),
			IsLast:      usrPS.IsLast(),
			Users:       homeUsers,
		},
		Genesis: &homeGenesis{
			ID:                     gen.ID().String(),
			GazPricePerKb:          gen.GazPricePerKb(),
			GazPriceInMatrixWorkKb: gen.GazPriceInMatrixWorkKb(),
			MaxAmountOfValidators:  gen.MaxAmountOfValidators(),
			UserID:                 gen.User().ID().String(),
			DepositID:              gen.Deposit().ID().String(),
		},
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "index.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	tmpl, err := template.New("home").Parse(string(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
		w.Write([]byte(str))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, templateData)
}

func (app *web) genesis(w http.ResponseWriter, r *http.Request) {
	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// retrieve the genesis:
	gen, genErr := app.genesisRepository.Retrieve()
	if genErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
		w.Write([]byte(str))
		return
	}

	// template structure:
	templateData := &homeGenesis{
		ID:                     gen.ID().String(),
		GazPricePerKb:          gen.GazPricePerKb(),
		GazPriceInMatrixWorkKb: gen.GazPriceInMatrixWorkKb(),
		MaxAmountOfValidators:  gen.MaxAmountOfValidators(),
		UserID:                 gen.User().ID().String(),
		DepositID:              gen.Deposit().ID().String(),
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "genesis.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	tmpl, err := template.New("genesis").Parse(string(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
		w.Write([]byte(str))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, templateData)
}

func (app *web) users(w http.ResponseWriter, r *http.Request) {
	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// retrieve the users associated with our conf PK:
	usrPS, usrPSErr := app.userRepository.RetrieveSetByPubKey(conf.WalletPK().PublicKey(), 0, amountOfElementsPerList)
	if usrPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the users entity set from creator's public key (PubKey: %s): %s", conf.WalletPK().PublicKey().String(), usrPSErr.Error())
		w.Write([]byte(str))
		return
	}

	usrs := []*homeUser{}
	usersIns := usrPS.Instances()
	for _, oneUserIns := range usersIns {
		if usr, ok := oneUserIns.(user.User); ok {
			usrs = append(usrs, &homeUser{
				ID:       usr.ID().String(),
				Shares:   usr.Shares(),
				WalletID: usr.Wallet().ID().String(),
			})
		}

		log.Printf("the given entity (ID: %s) is not a valid User instance", oneUserIns.ID().String())
		continue
	}

	// template structure:
	templateData := &homeUserList{
		Index:       usrPS.Index(),
		Amount:      usrPS.Amount(),
		TotalAmount: usrPS.TotalAmount(),
		IsLast:      usrPS.IsLast(),
		Users:       usrs,
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "users.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	tmpl, err := template.New("users").Parse(string(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
		w.Write([]byte(str))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, templateData)
}

func (app *web) wallets(w http.ResponseWriter, r *http.Request) {
	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// retrieve all the wallets:
	allWalPS, allWalPSErr := app.walletRepository.RetrieveSet(0, amountOfElementsPerList)
	if allWalPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the wallet entity set: %s", allWalPSErr.Error())
		w.Write([]byte(str))
		return
	}

	// retrieve the genesis:
	gen, genErr := app.genesisRepository.Retrieve()
	if genErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
		w.Write([]byte(str))
		return
	}

	walsIns := allWalPS.Instances()
	creatorWallets := []*homeWallet{}
	for _, oneWalletIns := range walsIns {
		if wal, ok := oneWalletIns.(walletpkg.Wallet); ok {
			// retrieve the wallet balance:
			bal, balErr := app.balanceRepository.RetrieveByWalletAndToken(wal, gen.Deposit().Token())
			if balErr != nil {
				log.Printf("there was an error while retrieving the wallet (ID: %s) balance of the given Token (ID: %s): %s", wal.ID().String(), gen.Deposit().Token().ID().String(), balErr.Error())
				continue
			}

			creatorWallets = append(creatorWallets, &homeWallet{
				ID:              wal.ID().String(),
				Creator:         wal.Creator().String(),
				ConcensusNeeded: wal.ConcensusNeeded(),
				TokenAmount:     bal.Amount(),
			})

			continue
		}

		log.Printf("the given entity (ID: %s) is not a valid Wallet instance", oneWalletIns.ID().String())
		continue
	}

	// template structure:
	templateData := &homeWalletList{
		Index:       allWalPS.Index(),
		Amount:      allWalPS.Amount(),
		TotalAmount: allWalPS.TotalAmount(),
		IsLast:      allWalPS.IsLast(),
		Wallets:     creatorWallets,
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "wallets.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	tmpl, err := template.New("wallets").Parse(string(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
		w.Write([]byte(str))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, templateData)
}

func (app *web) walletSingle(w http.ResponseWriter, r *http.Request) {
	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// get the id from the uri:
	vars := mux.Vars(r)
	if idAsString, ok := vars["id"]; ok {
		// convert the string to an id:
		id, idErr := uuid.FromString(idAsString)
		if idErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the wallet ID (%s) is invalid", idAsString)
			w.Write([]byte(str))
			return
		}

		// retrieve the wallet by id:
		wal, walErr := app.walletRepository.RetrieveByID(&id)
		if idErr != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(walErr.Error()))
			return
		}

		// retrieve the genesis:
		gen, genErr := app.genesisRepository.Retrieve()
		if genErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
			w.Write([]byte(str))
			return
		}

		// retrieve the balance:
		bal, balErr := app.balanceRepository.RetrieveByWalletAndToken(wal, gen.Deposit().Token())
		if balErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the wallet (ID: %s) balance of the given Token (ID: %s): %s", wal.ID().String(), gen.Deposit().Token().ID().String(), balErr.Error())
			w.Write([]byte(str))
			return
		}

		// retrieve the users:
		usrsPS, usrsPSErr := app.userRepository.RetrieveSetByWallet(wal, 0, amountOfElementsPerList)
		if usrsPSErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the user set related to the wallet (ID: %s): %s", wal.ID().String(), usrsPSErr.Error())
			w.Write([]byte(str))
			return
		}

		usrs := []*homeUser{}
		usersIns := usrsPS.Instances()
		for _, oneUserIns := range usersIns {
			if usr, ok := oneUserIns.(user.User); ok {
				usrs = append(usrs, &homeUser{
					ID:       usr.ID().String(),
					Shares:   usr.Shares(),
					WalletID: usr.Wallet().ID().String(),
				})
			}
		}

		// template structure:
		templateData := &singleWallet{
			ID:              wal.ID().String(),
			ConcensusNeeded: wal.ConcensusNeeded(),
			TokenAmount:     bal.Amount(),
			Users: &homeUserList{
				Index:       usrsPS.Index(),
				Amount:      usrsPS.Amount(),
				TotalAmount: usrsPS.TotalAmount(),
				IsLast:      usrsPS.IsLast(),
				Users:       usrs,
			},
		}

		// retrieve the html page:
		content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "wallet_single.html"))
		if contentErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
			w.Write([]byte(str))
			return
		}

		tmpl, err := template.New("walletSingle").Parse(string(content))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
			w.Write([]byte(str))
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, templateData)
	}

	w.WriteHeader(http.StatusInternalServerError)
	str := fmt.Sprintf("the ID could not be found")
	w.Write([]byte(str))
}

func (app *web) categories(w http.ResponseWriter, r *http.Request) {
	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// retrieve the categories with no parent:
	catWithNoParentPS, catWithNoParentPSErr := app.categoryRepository.RetrieveSetWithNoParent(0, amountOfElementsPerList)
	if catWithNoParentPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the category entity set: %s", catWithNoParentPSErr.Error())
		w.Write([]byte(str))
		return
	}

	cats := []*homeCategory{}
	catsIns := catWithNoParentPS.Instances()
	for _, oneCatIns := range catsIns {
		if cat, ok := oneCatIns.(category.Category); ok {
			oneCat := &homeCategory{
				ID:          cat.ID().String(),
				Name:        cat.Name(),
				Description: cat.Description(),
			}

			if cat.HasParent() {
				oneCat.ParentID = cat.Parent().ID().String()
			}

			cats = append(cats, oneCat)
			continue
		}

		log.Printf("the category (ID: %s) is not a valid Category instance", oneCatIns.ID().String())

	}

	// template structure:
	templateData := &homeCategoryList{
		Index:       catWithNoParentPS.Index(),
		Amount:      catWithNoParentPS.Amount(),
		TotalAmount: catWithNoParentPS.TotalAmount(),
		IsLast:      catWithNoParentPS.IsLast(),
		Categories:  cats,
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "categories.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	tmpl, err := template.New("categories").Parse(string(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
		w.Write([]byte(str))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, templateData)
}

func (app *web) register(w http.ResponseWriter, r *http.Request) {
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

		// retrieve the genesis:
		gen, genErr := app.genesisRepository.Retrieve()
		if genErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the genesis instance: %s", genErr.Error())
			w.Write([]byte(str))
			return
		}

		// create the user:
		usr := user.SDKFunc.Create(user.CreateParams{
			PubKey: conf.WalletPK().PublicKey(),
			Shares: 100,
			Wallet: walletpkg.SDKFunc.Create(walletpkg.CreateParams{
				Creator:         conf.WalletPK().PublicKey(),
				ConcensusNeeded: 100,
			}),
		})

		// convert the user to json:
		jsUserData, jsUserErr := cdc.MarshalJSON(usr)
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
		saveErr := app.accountService.Save(ac, int(math.Ceil(float64(gazPrice/100))))
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

		// retrieve the html page:
		content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "registered.html"))
		if contentErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
			w.Write([]byte(str))
			return
		}

		tmpl, err := template.New("registered").Parse(string(content))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
			w.Write([]byte(str))
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, encryptedConf)
		return
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "register.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func (app *web) newCategoriesForm(w http.ResponseWriter, r *http.Request) {

	if parseFormErr := r.ParseForm(); parseFormErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while parsing form elements: %s", parseFormErr.Error())
		w.Write([]byte(str))
		return
	}

	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	categoryName := r.FormValue("name")
	categoryDescription := r.FormValue("description")
	fromWalletID := r.FormValue("fromwalletid")
	if categoryName != "" && categoryDescription != "" && fromWalletID != "" {
		// parse the walletID:
		frmWalletID, frmWalletIDErr := uuid.FromString(fromWalletID)
		if frmWalletIDErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the given WalletID (ID: %s) is invalid: %s", frmWalletID, frmWalletIDErr.Error())
			w.Write([]byte(str))
			return
		}

		// retrieve the wallet:
		wal, walErr := app.walletRepository.RetrieveByID(&frmWalletID)
		if walErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the given WalletID (ID: %s) could not be retrieved: %s", frmWalletID.String(), walErr.Error())
			w.Write([]byte(str))
			return
		}

		// retrieve the user:
		usr, usrErr := app.userRepository.RetrieveByPubKeyAndWallet(conf.WalletPK().PublicKey(), wal)
		if usrErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the user (Pubkey: %s, WalletID: %s): %s", conf.WalletPK().PublicKey().String(), wal.ID().String(), usrErr.Error())
			w.Write([]byte(str))
			return
		}

		// create the new category instance:
		cat := category.SDKFunc.Create(category.CreateParams{
			Name:        categoryName,
			Description: categoryDescription,
		})

		// create the request:
		catRequest := request.SDKFunc.Create(request.CreateParams{
			FromUser:       usr,
			NewEntity:      cat,
			EntityMetaData: app.categoryRepresentation.MetaData(),
		})

		// save the request:
		saveErr := app.requestService.Save(catRequest, app.categoryRepresentation)
		if saveErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while saving a Request (Category instance): %s", saveErr.Error())
			w.Write([]byte(str))
			return
		}

		// redirect:
		uri := fmt.Sprintf("/requests/%s", catRequest.ID().String())
		http.Redirect(w, r, uri, http.StatusPermanentRedirect)
		return
	}

	// retrieve the users associated with our conf PK:
	usrPS, usrPSErr := app.userRepository.RetrieveSetByPubKey(conf.WalletPK().PublicKey(), 0, amountOfElementsPerList)
	if usrPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving the users entity set from creator's public key (PubKey: %s): %s", conf.WalletPK().PublicKey().String(), usrPSErr.Error())
		w.Write([]byte(str))
		return
	}

	usrs := []*homeUser{}
	usersIns := usrPS.Instances()
	for _, oneUserIns := range usersIns {
		if usr, ok := oneUserIns.(user.User); ok {
			usrs = append(usrs, &homeUser{
				ID:       usr.ID().String(),
				Shares:   usr.Shares(),
				WalletID: usr.Wallet().ID().String(),
			})
		}

		log.Printf("the given entity (ID: %s) is not a valid User instance", oneUserIns.ID().String())
		continue
	}

	// template structure:
	templateData := &homeCategoryNew{
		Users: &homeUserList{
			Index:       usrPS.Index(),
			Amount:      usrPS.Amount(),
			TotalAmount: usrPS.TotalAmount(),
			IsLast:      usrPS.IsLast(),
			Users:       usrs,
		},
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "categories_new.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	tmpl, err := template.New("categories_new").Parse(string(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
		w.Write([]byte(str))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, templateData)
}

func (app *web) requests(w http.ResponseWriter, r *http.Request) {
	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// retrieve the requests:
	reqPS, reqPSErr := app.requestRepository.RetrieveSet(0, amountOfElementsPerList)
	if reqPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retireving requests: %s", reqPSErr.Error())
		w.Write([]byte(str))
		return
	}

	reqs := []*homeRequest{}
	reqsIns := reqPS.Instances()
	for _, oneReqIns := range reqsIns {
		if req, ok := oneReqIns.(request.Request); ok {
			reqs = append(reqs, &homeRequest{
				ID:         req.ID().String(),
				FromUserID: req.From().ID().String(),
				NewName:    req.NewName(),
			})

		}

		log.Printf("the given entity (ID: %s) is not a valid Request instance", oneReqIns.ID().String())
		continue
	}

	// template structure:
	templateData := &homeRequestList{
		Index:       reqPS.Index(),
		Amount:      reqPS.Amount(),
		TotalAmount: reqPS.TotalAmount(),
		IsLast:      reqPS.IsLast(),
		Requests:    reqs,
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "requests.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	tmpl, err := template.New("requests").Parse(string(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
		w.Write([]byte(str))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, templateData)
	return

}

func (app *web) requestSingle(w http.ResponseWriter, r *http.Request) {
	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// parse the id:
	vars := mux.Vars(r)
	if reqIDAsString, ok := vars["id"]; ok {
		// parse the ID:
		reqID, reqIDErr := uuid.FromString(reqIDAsString)
		if reqIDErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the given ID (%s) is invalid", reqIDAsString)
			w.Write([]byte(str))
			return
		}

		// retrieve the request:
		req, reqErr := app.requestRepository.RetrieveByID(&reqID)
		if reqErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the Request (ID: %s): %s", reqID.String(), reqErr.Error())
			w.Write([]byte(str))
			return
		}

		// retrieve the votes:
		votesPS, votesPSErr := app.voteRepository.RetrieveSetByRequest(req, 0, amountOfElementsPerList)
		if votesPSErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the Vote set by Request (ID: %s): %s", reqID.String(), votesPSErr.Error())
			w.Write([]byte(str))
			return
		}

		vots := []*homeVote{}
		votsIns := votesPS.Instances()
		for _, oneVoteIns := range votsIns {
			if vot, ok := oneVoteIns.(vote.Vote); ok {
				vots = append(vots, &homeVote{
					ID:               vot.ID().String(),
					UserVoterID:      vot.Voter().ID().String(),
					UserAmountShares: vot.Voter().Shares(),
					IsApproved:       vot.IsApproved(),
				})

			}

			log.Printf("the given entity (ID: %s) is not a valid Vote instance", oneVoteIns.ID().String())
			continue
		}

		jsData, jsDataErr := json.MarshalIndent(req.New(), "", "    ")
		if jsDataErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while converting the entity in request (ID: %s): %s", reqID.String(), jsDataErr.Error())
			w.Write([]byte(str))
			return
		}

		// template structure:
		templateData := &homeRequestSingle{
			ID:         req.ID().String(),
			FromUserID: req.From().ID().String(),
			NewName:    req.NewName(),
			NewJS:      string(jsData),
			Votes: &homeVoteList{
				Index:       votesPS.Index(),
				Amount:      votesPS.Amount(),
				TotalAmount: votesPS.TotalAmount(),
				IsLast:      votesPS.IsLast(),
				Votes:       vots,
			},
		}

		// retrieve the html page:
		content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "requests_single.html"))
		if contentErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
			w.Write([]byte(str))
			return
		}

		tmpl, err := template.New("requests_single").Parse(string(content))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the template could not be rendered: %s", err.Error())
			w.Write([]byte(str))
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, templateData)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("the ID could not be found"))
}

func (app *web) requestSingleVote(w http.ResponseWriter, r *http.Request) {
	// retrieve the conf:
	conf := getConfigsFromCookie(loginCookieName, r)
	if conf == nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("cookie not found!")
		w.Write([]byte(str))
		return
	}

	// parse the id:
	vars := mux.Vars(r)
	if reqIDAsString, ok := vars["id"]; ok {
		// parse the ID:
		reqID, reqIDErr := uuid.FromString(reqIDAsString)
		if reqIDErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the given ID (%s) is invalid", reqIDAsString)
			w.Write([]byte(str))
			return
		}

		// retrieve the action:
		isApproved := false
		if action, ok := vars["action"]; ok {
			if action == "approved" {
				isApproved = true
			}
		}

		// retrieve the request:
		req, reqErr := app.requestRepository.RetrieveByID(&reqID)
		if reqErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the Request (ID: %s): %s", reqID.String(), reqErr.Error())
			w.Write([]byte(str))
			return
		}

		// retrieve all our users:
		usrPS, usrPSErr := app.userRepository.RetrieveSetByPubKey(conf.WalletPK().PublicKey(), 0, -1)
		if usrPSErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the user (Pubkey: %s): %s", conf.WalletPK().PublicKey().String(), usrPSErr.Error())
			w.Write([]byte(str))
			return
		}

		// we submit a vote request for each user we have:
		usrsIns := usrPS.Instances()
		for _, oneUsrIns := range usrsIns {
			if usr, ok := oneUsrIns.(user.User); ok {

				// create the vote:
				vot := vote.SDKFunc.Create(vote.CreateParams{
					Request:    req,
					Voter:      usr,
					IsApproved: isApproved,
				})

				// save the vote:
				reps := app.meta.WriteOnAllEntityRequest()
				if oneRep, ok := reps[req.NewName()]; ok {
					app.voteService.Save(vot, oneRep)
					continue
				}

				w.WriteHeader(http.StatusInternalServerError)
				str := fmt.Sprintf("the request entity (name: %s) cannot be voted on", req.NewName())
				w.Write([]byte(str))
				return

			}

			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("the entity (ID: %s) is not a valid User instance", oneUsrIns.ID().String())
			w.Write([]byte(str))
			return
		}

		// we redirect to the request:
		uri := fmt.Sprintf("/requests/%s", req.ID().String())
		http.Redirect(w, r, uri, http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("the ID could not be found"))
}
