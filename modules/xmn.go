package modules

import (
	"encoding/gob"
	"encoding/hex"
	"errors"

	keys "github.com/XMNBlockchain/datamint/keys"
	objects "github.com/XMNBlockchain/datamint/objects"
	roles "github.com/XMNBlockchain/datamint/roles"
	users "github.com/XMNBlockchain/datamint/users"
	amino "github.com/tendermint/go-amino"
	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	lua "github.com/yuin/gopher-lua"
)

var cdc = amino.NewCodec()

func init() {
	Register(cdc)
}

// Register registers all the interface -> struct to amino
func Register(codec *amino.Codec) {
	// crypto.PubKey
	func() {
		defer func() {
			recover()
		}()
		codec.RegisterInterface((*crypto.PubKey)(nil), nil)
		codec.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.Ed25519PubKeyAminoRoute, nil)
	}()
}

const luaKey = "xkeys"
const luaLists = "xlists"
const luaSets = "xsets"
const luaObjs = "xtables"
const luaUsers = "xusers"
const luaRoles = "xroles"
const luaCrypto = "xcrypto"

// XMN represents the XMN module:
type XMN struct {
	l    *lua.LState
	k    keys.Keys
	obj  objects.Objects
	usr  users.Users
	rols roles.Roles
}

// Crypto represents a crypto instance
type Crypto struct {
	pk crypto.PrivKey
}

func createXMN(l *lua.LState) *XMN {

	//create the instance:
	out := XMN{
		l:    l,
		k:    keys.SDKFunc.Create(),
		obj:  objects.SDKFunc.Create(),
		usr:  users.SDKFunc.Create(),
		rols: roles.SDKFunc.Create(),
	}

	//register the module on the lua state:
	out.register()

	//return the instance:
	return &out
}

func (app *XMN) register() {
	app.registerCrypto()
	app.registerKeys()
	app.registerObjects()
	app.registerUsers()
	app.registerRoles()
}

func (app *XMN) registerCrypto() {
	//verifies that the given type is a Crypto instance:
	checkFn := func(l *lua.LState) *Crypto {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(*Crypto); ok {
			return v
		}

		l.ArgError(1, "users expected")
		return nil
	}

	// create a new crypto instance:
	newCrypto := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = &Crypto{
			pk: ed25519.GenPrivKey(),
		}

		l.SetMetatable(ud, l.GetTypeMetatable(luaCrypto))
		l.Push(ud)
		return 1
	}

	//execute the pubKey command on the objects instance:
	pubKeyFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() == 1 {
			pubKeyAsBytes, pubKeyAsBytesErr := cdc.MarshalBinary(p.pk.PubKey())
			if pubKeyAsBytesErr != nil {
				l.ArgError(1, "the public key of the private key is invalid")
				return 1
			}

			pubKey := hex.EncodeToString(pubKeyAsBytes)
			l.Push(lua.LString(pubKey))
			return 1
		}

		l.ArgError(1, "the save func expected 0 parameter")
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"pubKey": pubKeyFn,
	}

	mt := app.l.NewTypeMetatable(luaCrypto)
	app.l.SetGlobal(luaCrypto, mt)

	// static attributes
	app.l.SetField(mt, "new", app.l.NewFunction(newCrypto))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}

func (app *XMN) registerKeys() {
	//verifies that the given type is a keys instance:
	checkKeys := func(l *lua.LState) keys.Keys {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(keys.Keys); ok {
			return v
		}

		l.ArgError(1, "keys expected")
		return nil
	}

	// load the Keys instance:
	loadKeys := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.k
		l.SetMetatable(ud, l.GetTypeMetatable(luaKey))
		l.Push(ud)
		return 1
	}

	// execute the retrieve command on the keys instance:
	keysRetrieve := func(l *lua.LState) int {
		p := checkKeys(l)
		amount := l.GetTop()
		if amount != 2 {
			l.ArgError(1, "the retrieve func expected 1 parameter")
			return 1
		}

		key := l.CheckString(2)
		value := p.Retrieve(key)
		if value == nil {
			l.Push(&lua.LNilType{})
			return 1
		}

		l.Push(lua.LString(value.(string)))
		return 1
	}

	//execute the save command on the keys instance:
	keysSave := func(l *lua.LState) int {
		p := checkKeys(l)
		if l.GetTop() == 3 {
			key := l.CheckString(2)
			value := l.CheckString(3)
			p.Save(key, value)
			return 0
		}

		l.ArgError(1, "the save func expected 2 parameters")
		return 1
	}

	// the keys methods:
	var keysMethods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkKeys(l)
			return app.lenFn(p)(l)
		},
		"exists": func(l *lua.LState) int {
			p := checkKeys(l)
			return app.existsFn(p)(l)
		},
		"retrieve": keysRetrieve,
		"search": func(l *lua.LState) int {
			p := checkKeys(l)
			return app.searchFn(p)(l)
		},
		"save": keysSave,
		"delete": func(l *lua.LState) int {
			p := checkKeys(l)
			return app.delFn(p)(l)
		},
	}

	mt := app.l.NewTypeMetatable(luaKey)
	app.l.SetGlobal(luaKey, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadKeys))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), keysMethods))
}

func (app *XMN) registerObjects() {

	//gob register:
	gob.Register(map[string]interface{}{})

	//verifies that the given type is a keys instance:
	checkObjects := func(l *lua.LState) objects.Objects {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(objects.Objects); ok {
			return v
		}

		l.ArgError(1, "objects expected")
		return nil
	}

	// load the Objects instance:
	loadObjects := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.obj
		l.SetMetatable(ud, l.GetTypeMetatable(luaObjs))
		l.Push(ud)
		return 1
	}

	//execute the save command on the objects instance:
	saveFn := func(l *lua.LState) int {
		p := checkObjects(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the save func expected at least 2 parameters")
			return 1
		}

		params := []*objects.ObjInKey{}
		for i := 2; i <= amount; i++ {
			oneObjInKey := objects.ObjInKey{}
			oneParam := l.CheckTable(i)
			oneParam.ForEach(func(name lua.LValue, value lua.LValue) {
				valueType := value.Type()
				nameAsString := name.String()
				if nameAsString == "table" && valueType == lua.LTTable {
					oneObjInKey.Obj = app.convertLTableToHashMap(value.(*lua.LTable))
				}

				if nameAsString == "key" && valueType == lua.LTString {
					oneObjInKey.Key = value.String()
				}
			})

			params = append(params, &oneObjInKey)
		}

		amountSaved := p.Save(params...)
		l.Push(lua.LNumber(amountSaved))
		return 1
	}

	//execute the retrieve command on the objects instance:
	retrieveFn := func(l *lua.LState) int {
		p := checkObjects(l)
		if l.GetTop() != 2 {
			l.ArgError(1, "the save func expected 1 parameter")
			return 1
		}

		objInKey := objects.ObjInKey{
			Key: l.CheckString(2),
			Obj: new(map[string]interface{}),
		}

		p.Retrieve(&objInKey)
		mapResult := objInKey.Obj.(*map[string]interface{})
		ltable := app.convertHashMapToLTable(*mapResult)
		l.Push(ltable)
		return 1
	}

	// the objects methods:
	var methods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkObjects(l)
			return app.lenFn(p.Keys())(l)
		},
		"exists": func(l *lua.LState) int {
			p := checkObjects(l)
			return app.existsFn(p.Keys())(l)
		},
		"retrieve": retrieveFn,
		"search": func(l *lua.LState) int {
			p := checkObjects(l)
			return app.searchFn(p.Keys())(l)
		},
		"save": saveFn,
		"delete": func(l *lua.LState) int {
			p := checkObjects(l)
			return app.delFn(p.Keys())(l)
		},
	}

	mt := app.l.NewTypeMetatable(luaObjs)
	app.l.SetGlobal(luaObjs, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadObjects))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}

func (app *XMN) registerUsers() {
	//verifies that the given type is a Users instance:
	checkFn := func(l *lua.LState) users.Users {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(users.Users); ok {
			return v
		}

		l.ArgError(1, "users expected")
		return nil
	}

	// load the Users instance:
	loadUsers := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.usr
		l.SetMetatable(ud, l.GetTypeMetatable(luaUsers))
		l.Push(ud)
		return 1
	}

	//execute the key command on the objects instance:
	keyFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() != 2 {
			l.ArgError(1, "the save func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := app.fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		key := p.Key(pubKey)
		l.Push(lua.LString(key))
		return 1
	}

	//execute the exists command on the objects instance:
	existsFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := app.fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		exists := p.Exists(pubKey)
		l.Push(lua.LBool(exists))
		return 1
	}

	//execute the insert command on the objects instance:
	insertFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := app.fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		isInserted := p.Insert(pubKey)
		l.Push(lua.LBool(isInserted))
		return 1
	}

	//execute the delete command on the objects instance:
	deleteFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() < 2 {
			l.ArgError(1, "the exists func expected 1 parameter")
			return 1
		}

		pubKeyAsString := l.CheckString(2)
		pubKey, pubKeyErr := app.fromStringToPubKey(pubKeyAsString)
		if pubKeyErr != nil {
			l.ArgError(1, pubKeyErr.Error())
			return 1
		}

		isDeleted := p.Delete(pubKey)
		l.Push(lua.LBool(isDeleted))
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"len": func(l *lua.LState) int {
			p := checkFn(l)
			return app.lenFn(p.Objects().Keys())(l)
		},
		"key":    keyFn,
		"exists": existsFn,
		"insert": insertFn,
		"delete": deleteFn,
	}

	mt := app.l.NewTypeMetatable(luaUsers)
	app.l.SetGlobal(luaUsers, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadUsers))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))
}

func (app *XMN) registerRoles() {
	//verifies that the given type is a Roles instance:
	checkFn := func(l *lua.LState) roles.Roles {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(roles.Roles); ok {
			return v
		}

		l.ArgError(1, "roles expected")
		return nil
	}

	// load the Roles instance:
	loadRoles := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = app.rols
		l.SetMetatable(ud, l.GetTypeMetatable(luaRoles))
		l.Push(ud)
		return 1
	}

	//execute the add command on the roles instance:
	addFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		pubKeys := []crypto.PubKey{}
		key := l.CheckString(2)
		for i := 3; i <= amount; i++ {
			pubKeyAsString := l.CheckString(i)
			pubKey, pubKeyErr := app.fromStringToPubKey(pubKeyAsString)
			if pubKeyErr != nil {
				l.ArgError(1, pubKeyErr.Error())
				return 1
			}

			pubKeys = append(pubKeys, pubKey)
		}

		amountAdded := p.Add(key, pubKeys...)
		l.Push(lua.LNumber(amountAdded))
		return 1
	}

	//execute the del command on the roles instance:
	delFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		pubKeys := []crypto.PubKey{}
		key := l.CheckString(2)
		for i := 3; i <= amount; i++ {
			pubKeyAsString := l.CheckString(i)
			pubKey, pubKeyErr := app.fromStringToPubKey(pubKeyAsString)
			if pubKeyErr != nil {
				l.ArgError(1, pubKeyErr.Error())
				return 1
			}

			pubKeys = append(pubKeys, pubKey)
		}

		amountAdded := p.Del(key, pubKeys...)
		l.Push(lua.LNumber(amountAdded))
		return 1
	}

	//execute the enableWriteAccess command on the roles instance:
	enableWriteAccessFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		key := l.CheckString(2)
		patterns := []string{}
		for i := 3; i <= amount; i++ {
			patterns = append(patterns, l.CheckString(i))
		}

		amountEnabled := p.EnableWriteAccess(key, patterns...)
		l.Push(lua.LNumber(amountEnabled))
		return 1
	}

	//execute the disableWriteAccess command on the roles instance:
	disableWriteAccessFn := func(l *lua.LState) int {
		p := checkFn(l)
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the exists func expected ast least 2 parameters")
			return 1
		}

		key := l.CheckString(2)
		patterns := []string{}
		for i := 3; i <= amount; i++ {
			patterns = append(patterns, l.CheckString(i))
		}

		amountEnabled := p.DisableWriteAccess(key, patterns...)
		l.Push(lua.LNumber(amountEnabled))
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"add":                addFn,
		"del":                delFn,
		"enableWriteAccess":  enableWriteAccessFn,
		"disableWriteAccess": disableWriteAccessFn,
	}

	mt := app.l.NewTypeMetatable(luaRoles)
	app.l.SetGlobal(luaRoles, mt)

	// static attributes
	app.l.SetField(mt, "load", app.l.NewFunction(loadRoles))

	// methods
	app.l.SetField(mt, "__index", app.l.SetFuncs(app.l.NewTable(), methods))

}

func (app *XMN) convertHashMapToLTable(hashmap map[string]interface{}) *lua.LTable {
	out := lua.LTable{}
	for keyname, value := range hashmap {
		if subHashMap, ok := value.(map[string]interface{}); ok {
			subLTable := app.convertHashMapToLTable(subHashMap)
			out.RawSet(lua.LString(keyname), subLTable)
			continue
		}

		out.RawSet(lua.LString(keyname), lua.LString(value.(string)))
	}

	return &out
}

func (app *XMN) convertLTableToHashMap(table *lua.LTable) map[string]interface{} {
	hashmap := map[string]interface{}{}
	table.ForEach(func(name lua.LValue, value lua.LValue) {
		if value.Type() == lua.LTTable {
			hashmap[name.String()] = app.convertLTableToHashMap(value.(*lua.LTable))
			return
		}

		hashmap[name.String()] = value.String()
	})

	return hashmap
}

func (app *XMN) existsFn(p keys.Keys) lua.LGFunction {
	fn := func(l *lua.LState) int {
		amount := l.GetTop()
		if amount < 2 {
			l.ArgError(1, "the save func expected 0 parameter")
			return 1
		}

		keys := []string{}
		for i := 2; i <= amount; i++ {
			oneKey := l.CheckString(i)
			keys = append(keys, oneKey)
		}

		existsAmount := p.Exists(keys...)
		l.Push(lua.LNumber(existsAmount))
		return 1
	}

	return fn
}

func (app *XMN) lenFn(p keys.Keys) lua.LGFunction {
	fn := func(l *lua.LState) int {
		if l.GetTop() == 1 {
			amount := p.Len()
			l.Push(lua.LNumber(amount))
			return 1
		}

		l.ArgError(1, "the save func expected 0 parameter")
		return 1
	}

	return fn
}

func (app *XMN) searchFn(p keys.Keys) lua.LGFunction {
	fn := func(l *lua.LState) int {
		amount := l.GetTop()
		if amount != 2 {
			l.ArgError(1, "the retrieve func expected 1 parameter")
			return 1
		}

		pattern := l.CheckString(2)
		results := p.Search(pattern)

		keys := lua.LTable{}
		for index, oneResult := range results {
			keys.Insert(index, lua.LString(oneResult))
		}

		l.Push(&keys)
		return 1
	}

	return fn
}

func (app *XMN) delFn(p keys.Keys) lua.LGFunction {
	fn := func(l *lua.LState) int {
		amount := l.GetTop()
		if amount < 1 {
			l.ArgError(1, "the retrieve func expected at least 1 parameter")
			return 1
		}

		keys := []string{}
		for i := 2; i <= amount; i++ {
			oneKey := l.CheckString(i)
			keys = append(keys, oneKey)
		}

		amountDeleted := p.Delete(keys...)
		l.Push(lua.LNumber(amountDeleted))
		return 1
	}

	return fn
}

func (app *XMN) fromStringToPubKey(pubKeyAsString string) (crypto.PubKey, error) {
	pubKeyAsBytes, pubKeyAsBytesErr := hex.DecodeString(pubKeyAsString)
	if pubKeyAsBytesErr != nil {
		return nil, errors.New("the encoded public key is invalid")
	}

	pubKey := new(ed25519.PubKeyEd25519)
	pubKeyErr := cdc.UnmarshalBinary(pubKeyAsBytes, pubKey)
	if pubKeyErr != nil {
		return nil, errors.New("the public key []byte is invalid")
	}

	return pubKey, nil
}
