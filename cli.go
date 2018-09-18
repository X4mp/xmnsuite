package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	cliapp "github.com/urfave/cli"
	tendermint "github.com/xmnservices/xmnsuite/blockchains/tendermint"
	crypto "github.com/xmnservices/xmnsuite/crypto"
	datastore "github.com/xmnservices/xmnsuite/datastore"
	modules "github.com/xmnservices/xmnsuite/modules"
	module_chain "github.com/xmnservices/xmnsuite/modules/chain"
	module_crypto "github.com/xmnservices/xmnsuite/modules/crypto"
	module_datastore "github.com/xmnservices/xmnsuite/modules/datastore"
	json_module "github.com/xmnservices/xmnsuite/modules/json"
	module_sdk "github.com/xmnservices/xmnsuite/modules/sdk"
	uuid_module "github.com/xmnservices/xmnsuite/modules/uuid"
	lua "github.com/yuin/gopher-lua"
)

type loadModuleFn func() error

type cli struct {
	luaContext      *lua.LState
	cliContext      *cliapp.Context
	loadedModules   map[string]modules.Module
	possibleModules map[string]loadModuleFn
}

func createCLI(context *cliapp.Context) (*cli, error) {
	// retrieve the context size:
	ccSizeAsString := context.String("ccsize")
	ccSize, ccSizeErr := strconv.Atoi(ccSizeAsString)
	if ccSizeErr != nil {
		// log:
		log.Printf("there was an error while converting a string to an int: %s", ccSizeErr.Error())

		// output error:
		str := fmt.Sprintf("the ccsize param (%s) must be an int", ccSizeAsString)
		return nil, errors.New(str)
	}

	// retrieve the registry size:
	rSizeAsString := context.String("rsize")
	rSize, rSizeErr := strconv.Atoi(rSizeAsString)
	if rSizeErr != nil {
		// log:
		log.Printf("there was an error while converting a string to an int: %s", rSizeErr.Error())

		// output error:
		str := fmt.Sprintf("the rsize param (%s) must be an int", rSizeAsString)
		return nil, errors.New(str)
	}

	out := cli{
		cliContext: context,
		luaContext: lua.NewState(lua.Options{
			CallStackSize: ccSize,
			RegistrySize:  rSize,
		}),
		loadedModules:   map[string]modules.Module{},
		possibleModules: map[string]loadModuleFn{},
	}

	// load the modules:
	moduleNames := strings.Split(context.String("modules"), ",")
	out.possibleModules = map[string]loadModuleFn{
		"chain":     out.loadChainModule,
		"crypto":    out.loadCryptoModule,
		"datastore": out.loadDatastoreModule,
		"sdk":       out.loadSDKModule,
		"json":      out.loadJSONModule,
		"uuid":      out.loadUUIDModule,
	}

	for _, oneModuleName := range moduleNames {
		out.loadModuleByName(oneModuleName)
	}

	return &out, nil
}

func (app *cli) getModuleByName(moduleName string) modules.Module {
	if _, ok := app.loadedModules[moduleName]; ok {
		return app.loadedModules[moduleName]
	}

	return nil
}

func (app *cli) execute(scriptPath string) error {
	doFileErr := app.luaContext.DoFile(scriptPath)
	if doFileErr != nil {
		return doFileErr
	}

	return nil
}

func (app *cli) loadModuleByName(moduleName string) {
	if loadModFn, ok := app.possibleModules[moduleName]; ok {
		log.Printf("loading module: %s...", moduleName)
		loadErr := loadModFn()
		if loadErr != nil {
			log.Printf("error while loading module (%s): %s", moduleName, loadErr.Error())
		}

		return
	}

	log.Printf("the module name (%s) is invalid, skip...", moduleName)
}

func (app *cli) loadUUIDModule() error {
	if _, ok := app.loadedModules["uuid"]; ok {
		return nil
	}

	app.loadedModules["uuid"] = uuid_module.SDKFunc.Create(uuid_module.CreateParams{
		Context: app.luaContext,
	})

	return nil
}

func (app *cli) loadJSONModule() error {
	if _, ok := app.loadedModules["json"]; ok {
		return nil
	}

	app.loadedModules["json"] = json_module.SDKFunc.Create(json_module.CreateParams{
		Context: app.luaContext,
	})

	return nil
}

func (app *cli) loadChainModule() error {
	if _, ok := app.loadedModules["chain"]; ok {
		return nil
	}

	if _, ok := app.loadedModules["json"]; !ok {
		app.loadModuleByName("json")
	}

	if _, ok := app.loadedModules["datastore"]; !ok {
		app.loadModuleByName("datastore")
	}

	// dbpath:
	dbPath := app.cliContext.String("dbpath")

	// nodepk:
	nodePkAsString := app.cliContext.String("nodepk")
	nodePKAsBytes, nodePKAsBytesErr := hex.DecodeString(nodePkAsString)
	if nodePKAsBytesErr != nil {
		// log:
		log.Printf("there was an error while decoding a string to hex: %s", nodePKAsBytesErr.Error())

		// output error:
		str := fmt.Sprintf("the given nodepk (%s) is not a valid private key", nodePkAsString)
		return errors.New(str)
	}

	nodePK := new(ed25519.PrivKeyEd25519)
	nodePKErr := cdc.UnmarshalBinaryBare(nodePKAsBytes, nodePK)
	if nodePKErr != nil {
		// log:
		log.Printf("there was an error while Unmarshalling []byte to PrivateKey instance: %s", nodePKErr.Error())

		// output error:
		str := fmt.Sprintf("the given nodepk (%s) is not a valid private key", nodePkAsString)
		return errors.New(str)
	}

	// id:
	idAsString := app.cliContext.String("id")
	id, idErr := uuid.FromString(idAsString)
	if idErr != nil {
		// log:
		log.Printf("there was an error while converting a string to an ID: %s", idErr.Error())

		// output error:
		str := fmt.Sprintf("the given id (%s) is not a valid ID", idAsString)
		return errors.New(str)
	}

	// rpubkeys:
	rootPubKeys := []crypto.PublicKey{}
	rootPubKeyAsCommaSeperatedString := app.cliContext.String("rpubkeys")
	if rootPubKeyAsCommaSeperatedString != "" {
		rootPubKeysAsString := strings.Split(rootPubKeyAsCommaSeperatedString, ",")
		for _, oneRootPubKeyAsString := range rootPubKeysAsString {
			pubKey := crypto.SDKFunc.CreatePubKey(crypto.CreatePubKeyParams{
				PubKeyAsString: oneRootPubKeyAsString,
			})

			// add the pubkey to our list:
			rootPubKeys = append(rootPubKeys, pubKey)

			// log:
			log.Printf("adding root pub key: %s", oneRootPubKeyAsString)
		}
	}

	// create module:
	app.loadedModules["chain"] = module_chain.SDKFunc.Create(module_chain.CreateParams{
		Context:     app.luaContext,
		DBPath:      dbPath,
		ID:          &id,
		RootPubKeys: rootPubKeys,
		NodePK:      nodePK,
		Datastore:   app.loadedModules["datastore"].(module_datastore.Datastore),
	})

	return nil
}

func (app *cli) loadCryptoModule() error {
	if _, ok := app.loadedModules["crypto"]; ok {
		return nil
	}

	app.loadedModules["crypto"] = module_crypto.SDKFunc.Create(module_crypto.CreateParams{
		Context: app.luaContext,
	})

	return nil
}

func (app *cli) loadDatastoreModule() error {
	if _, ok := app.loadedModules["datastore"]; ok {
		return nil
	}

	if _, ok := app.loadedModules["json"]; !ok {
		app.loadModuleByName("json")
	}

	// create the datastore:
	ds := datastore.SDKFunc.Create()

	// create the module:
	app.loadedModules["datastore"] = module_datastore.SDKFunc.Create(module_datastore.CreateParams{
		Context:   app.luaContext,
		Datastore: ds,
	})
	return nil
}

func (app *cli) loadSDKModule() error {
	if _, ok := app.loadedModules["sdk"]; ok {
		return nil
	}

	if _, ok := app.loadedModules["json"]; !ok {
		app.loadModuleByName("json")
	}

	// connect to:
	connectorAsString := app.cliContext.String("connector")
	if connectorAsString != "" {
		appService := tendermint.SDKFunc.CreateApplicationService()
		client, clientErr := appService.Connect(connectorAsString)
		if clientErr != nil {
			// log:
			log.Printf("there was an error while connecting to the given host: %s", clientErr.Error())

			// output error:
			str := fmt.Sprintf("the given connector (%s) is not a valid blockchain host", connectorAsString)
			return errors.New(str)
		}

		app.loadedModules["sdk"] = module_sdk.SDKFunc.Create(module_sdk.CreateParams{
			Context: app.luaContext,
			Client:  client,
		})

		return nil
	}

	return errors.New("the connector param is mandatory in order to load the sdk module")
}
