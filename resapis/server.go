package restapis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	crypto "github.com/xmnservices/xmnsuite/crypto"
)

type server struct {
	pk              crypto.PrivateKey
	rep             *repository
	serv            *service
	port            int
	gracefulTimeout time.Duration
	server          *http.Server
}

func createServer(
	pk crypto.PrivateKey,
	rter *mux.Router,
	rep *repository,
	ser *service,
	port int,
	gracefulTimeout time.Duration,
) Server {

	serv := server{
		rep:             rep,
		serv:            ser,
		port:            port,
		server:          nil,
		gracefulTimeout: gracefulTimeout,
		pk:              pk,
	}

	// setup the routes:
	rter.HandleFunc("/", serv.home).Methods("GET")

	rter.HandleFunc("/accounts", serv.createAccount).Methods("POST")
	rter.HandleFunc("/accounts", serv.retrieveAccounts).Methods("GET")
	rter.HandleFunc("/accounts/{name:[a-zA-Z0-9-]+}", serv.retrieveAccountByName).Methods("GET")

	//rter.HandleFunc("/contracts", serv.home).Methods("POST")
	//rter.HandleFunc("/contracts", serv.home).Methods("GET")
	//rter.HandleFunc("/contracts/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", serv.home).Methods("GET")

	//rter.HandleFunc("/contracts/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/instances", serv.home).Methods("POST")
	//rter.HandleFunc("/contracts/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/instances", serv.home).Methods("GET")
	//rter.HandleFunc("/contracts/{contractID:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/instances/{inatanceID:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", serv.home).Methods("GET")
	//rter.HandleFunc("/contracts/{contractID:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/instances/{inatanceID:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", serv.home).Methods("DELETE")

	// middleware func to verify that the http call is authenticated:
	rter.Use(serv.authenticate)

	// server:
	serv.server = &http.Server{
		Handler: rter,
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &serv
}

// Start starts the server application
func (serv *server) Start() error {
	go func(srv *http.Server) {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}(serv.server)

	return nil
}

// Stop stops the server application
func (serv *server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), serv.gracefulTimeout)
	defer cancel()

	return serv.server.Shutdown(ctx)
}

func (serv *server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}()

		parseErr := r.ParseForm()
		if parseErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while parsing the form: %s", parseErr.Error())
			w.Write([]byte(str))
			return
		}

		sig := crypto.SDKFunc.CreateSig(crypto.CreateSigParams{
			SigAsString: r.Header.Get("X-Session-Token"),
		})

		token := createToken(r.Method, r.RequestURI, r.PostForm).Hash()
		if !sig.PublicKey(token).Equals(serv.pk.PublicKey()) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
		return
	})
}

func (serv *server) home(w http.ResponseWriter, r *http.Request) {
	return
}

func (serv *server) createAccount(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm()
	if parseErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while parsing the form: %s", parseErr.Error())
		w.Write([]byte(str))
		return
	}

	name := r.FormValue("name")
	if serv.rep.Exists(name) {
		w.WriteHeader(http.StatusConflict)
		str := fmt.Sprintf("the given account name (%s) already exists", name)
		w.Write([]byte(str))
		return
	}

	jsSeedWords := r.FormValue("seedwords")
	seedWords := new([]string)
	jsErr := json.Unmarshal([]byte(jsSeedWords), seedWords)
	if jsErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		str := fmt.Sprintf("there was an error while unmarshalling the given json seed words (%s) to []string: %s", jsSeedWords, jsErr.Error())
		w.Write([]byte(str))
		return
	}

	pk := crypto.SDKFunc.GenPK()
	encryptedPK := crypto.SDKFunc.Encrypt(crypto.EncryptParams{
		Pass: []byte(strings.Join(*seedWords, "|")),
		Msg:  []byte(pk.String()),
	})

	acc, accErr := createAccount(r.FormValue("name"), encryptedPK, time.Now().UTC())
	if accErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		str := fmt.Sprintf("there was an error while creating the account instance: %s", accErr.Error())
		w.Write([]byte(str))
		return
	}

	saveErr := serv.serv.Save(name, acc)
	if saveErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while marshalling an account to json: %s", saveErr.Error())
		w.Write([]byte(str))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
	return
}

func (serv *server) retrieveAccounts(w http.ResponseWriter, r *http.Request) {
	names, namesErr := serv.rep.RetrieveNames()
	if namesErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while retireving the names: %s", namesErr.Error())
		w.Write([]byte(str))
		return
	}

	accs := []Account{}
	for _, oneName := range names {
		acc := new(account)
		retErr := serv.rep.Retrieve(oneName, acc)
		if retErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the account: %s", retErr.Error())
			w.Write([]byte(str))
			return
		}

		accs = append(accs, acc)
	}

	js, jsErr := json.Marshal(accs)
	if jsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		str := fmt.Sprintf("there was an error while marshalling an account to json: %s", jsErr.Error())
		w.Write([]byte(str))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
	return
}

func (serv *server) retrieveAccountByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if name, ok := vars["name"]; ok {
		if !serv.rep.Exists(name) {
			w.WriteHeader(http.StatusNotFound)
			str := fmt.Sprintf("the account (name: %s) could not be found", name)
			w.Write([]byte(str))
			return
		}

		acc := new(account)
		retErr := serv.rep.Retrieve(name, acc)
		if retErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while retrieving the account: %s", retErr.Error())
			w.Write([]byte(str))
			return
		}

		js, jsErr := json.Marshal(acc)
		if jsErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			str := fmt.Sprintf("there was an error while marshalling an account to json: %s", jsErr.Error())
			w.Write([]byte(str))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(js)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("the name is mandatory in order to retrieve an account by name"))
	return
}
