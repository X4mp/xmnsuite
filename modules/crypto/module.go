package crypto

import (
	"encoding/hex"
	"fmt"

	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	applications "github.com/xmnservices/xmnsuite/applications"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

const luaPrivKey = "privkey"

type privKey struct {
	pk crypto.PrivKey
}

type resourcePointer struct {
	ptr applications.ResourcePointer
}

type resource struct {
	ptr applications.Resource
}

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
	checkFn := func(l *lua.LState) *privKey {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(*privKey); ok {
			return v
		}

		l.ArgError(1, "private key expected")
		return nil
	}

	// create a new crypto instance:
	newPrivKey := func(l *lua.LState) int {

		if l.GetTop() == 1 {
			privKeyAsString := l.CheckString(1)
			privKeyAsBytes, privKeyAsBytesErr := hex.DecodeString(privKeyAsString)
			if privKeyAsBytesErr != nil {
				str := fmt.Sprintf("the given private key could not be converted from hex to []byte: %s", privKeyAsBytesErr.Error())
				l.ArgError(1, str)
				return 1
			}

			newPrivKey := new(ed25519.PrivKeyEd25519)
			unmarshalErr := cdc.UnmarshalBinaryBare(privKeyAsBytes, newPrivKey)
			if unmarshalErr != nil {
				str := fmt.Sprintf("the given private key could not be converted from []byte to PrivateKey:  %s", unmarshalErr.Error())
				l.ArgError(1, str)
				return 1
			}

			ud := l.NewUserData()
			ud.Value = &privKey{
				pk: newPrivKey,
			}

			l.SetMetatable(ud, l.GetTypeMetatable(luaPrivKey))
			l.Push(ud)
			return 1

		}

		ud := l.NewUserData()
		ud.Value = &privKey{
			pk: ed25519.GenPrivKey(),
		}

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

		pubKeyAsBytes, pubKeyAsBytesErr := cdc.MarshalBinaryBare(p.pk.PubKey())
		if pubKeyAsBytesErr != nil {
			l.ArgError(1, "the public key of the private key is invalid")
			return 1
		}

		pubKey := hex.EncodeToString(pubKeyAsBytes)
		l.Push(lua.LString(pubKey))
		return 1
	}

	//execute the signFn command on the privkey instance:
	signFn := func(l *lua.LState) int {
		p := checkFn(l)
		if l.GetTop() == 2 {
			encodedHash := l.CheckString(2)
			msg, msgErr := hex.DecodeString(encodedHash)
			if msgErr != nil {
				str := fmt.Sprintf("the hash could not be decoded: %s", msgErr.Error())
				l.ArgError(1, str)
				return 1
			}

			sig, sigErr := p.pk.Sign(msg)
			if sigErr != nil {
				str := fmt.Sprintf("there was an error while signing the given message: %s", sigErr.Error())
				l.ArgError(1, str)
				return 1
			}

			encocedSig := hex.EncodeToString(sig)
			l.Push(lua.LString(encocedSig))
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
