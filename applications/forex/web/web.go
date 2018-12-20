package web

import (
	"bytes"
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
	"github.com/xmnservices/xmnsuite/applications/forex/objects/category"
	"github.com/xmnservices/xmnsuite/applications/forex/objects/currency"
	"github.com/xmnservices/xmnsuite/applications/forex/web/controllers/banks"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account"
	walletpkg "github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/wallet/entities/user"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/account/work"
	"github.com/xmnservices/xmnsuite/blockchains/core/objects/entity/entities/genesis"
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
	entityService          entity.Service
	accountService         account.Service
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
		entityService:          entityService,
		accountService:         accountService,
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
	app.rter.HandleFunc("/categories", app.categories)
	app.rter.HandleFunc("/categories/new", app.newCategoriesForm)

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

				log.Printf("\n login: %s \n", decrypted.String())

				// set the cookie:
				http.SetCookie(w, &http.Cookie{
					Name:  loginCookieName,
					Value: decrypted.String(),
				})

				// render the continue page:
				// retrieve the html page:
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

func (app *web) categories(w http.ResponseWriter, r *http.Request) {
	// retrieve the categories:
	/*catPS, catPSErr := app.categoryRepository.RetrieveSet(0, amountOfElementsPerList)
	if catPSErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retrieving a Category instances: %s", catPSErr.Error())
		w.Write([]byte(str))
		return
	}*/

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "categories.html"))
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

	categoryName := r.FormValue("name")
	categoryDescription := r.FormValue("description")
	if categoryName != "" && categoryDescription != "" {
		// create the new category instance:
		cat := category.SDKFunc.Create(category.CreateParams{
			Name:        categoryName,
			Description: categoryDescription,
		})

		// create the request:
		/*catRequest := request.SDKFunc.Create(request.CreateParams{
			FromUser:       genIns.User(),
			NewEntity:      cat,
			EntityMetaData: category.SDKFunc.CreateMetaData(),
		})*/

		// save the instance:
		saveErr := app.entityService.Save(cat, app.categoryRepresentation)
		if saveErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while saving a Category instance: %s", saveErr.Error())
			w.Write([]byte(str))
			return
		}

		// redirect to the votes:
		http.Redirect(w, r, "/requests", http.StatusSeeOther)
	}

	// retrieve the html page:
	content, contentErr := ioutil.ReadFile(filepath.Join(app.templateDir, "categories_new.html"))
	if contentErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("the template could not be rendered: %s", contentErr.Error())
		w.Write([]byte(str))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
