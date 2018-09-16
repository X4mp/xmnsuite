package crypto

import (
	crypto "github.com/xmnservices/xmnsuite/crypto"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

const luaPrivKey = "privkey"

type module struct {
	context *lua.LState
}

func createModule(context *lua.LState) Crypto {
	out := module{
		context: context,
	}

	out.register()

	return &out
}

func (app *module) register() {
	// preload JSON:
	luajson.Preload(app.context)

	// preload XMN:
	app.context.PreloadModule("crypto", func(context *lua.LState) int {
		app.registerPrivKey(context)
		return 1
	})
}

func (app *module) registerPrivKey(context *lua.LState) {
	//verifies that the given type is a Crypto instance:
	checkFn := func(l *lua.LState) crypto.PrivateKey {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(crypto.PrivateKey); ok {
			return v
		}

		l.ArgError(1, "private key expected")
		return nil
	}

	// create a new crypto instance:
	newPrivKey := func(l *lua.LState) int {

		if l.GetTop() == 1 {
			privKeyAsString := l.CheckString(1)
			privKey := crypto.SDKFunc.CreatePK(crypto.CreatePKParams{
				PKAsString: privKeyAsString,
			})

			ud := l.NewUserData()
			ud.Value = privKey

			l.SetMetatable(ud, l.GetTypeMetatable(luaPrivKey))
			l.Push(ud)
			return 1

		}

		ud := l.NewUserData()
		ud.Value = crypto.SDKFunc.GenPK()

		l.SetMetatable(ud, l.GetTypeMetatable(luaPrivKey))
		l.Push(ud)
		return 1
	}

	//execute the pubKey command on the privkey instance:
	pubKeyFn := func(l *lua.LState) int {
		p := checkFn(l)
		if p == nil {
			return 1
		}

		pubKey := p.PublicKey()
		l.Push(lua.LString(pubKey.String()))
		return 1
	}

	//execute the signFn command on the privkey instance:
	signFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() == 2 {
			msg := l.CheckString(2)
			sig := p.Sign(msg)
			l.Push(lua.LString(sig.String()))
			return 1
		}

		l.ArgError(1, "the save func expected 1 parameter")
		return 1
	}

	// the users methods:
	var methods = map[string]lua.LGFunction{
		"pubKey": pubKeyFn,
		"sign":   signFn,
	}

	mt := context.NewTypeMetatable(luaPrivKey)
	context.SetGlobal(luaPrivKey, mt)

	// static attributes
	context.SetField(mt, "new", context.NewFunction(newPrivKey))

	// methods
	context.SetField(mt, "__index", context.SetFuncs(context.NewTable(), methods))
}
