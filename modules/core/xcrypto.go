package core

import (
	"encoding/hex"

	crypto "github.com/tendermint/tendermint/crypto"
	ed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	lua "github.com/yuin/gopher-lua"
)

const luaCrypto = "xcrypto"

// XCrypto represents the xcrypto instance
type XCrypto struct {
	l  *lua.LState
	pk crypto.PrivKey
}

// CreateXCrypto creates a new XCrypto instance:
func CreateXCrypto(l *lua.LState) *XCrypto {
	// create the instance:
	out := XCrypto{
		l: l,
	}

	//registers the xcrypto module on the current lua state:
	out.register()

	//returns:
	return &out
}

func (app *XCrypto) register() {
	//verifies that the given type is a Crypto instance:
	checkFn := func(l *lua.LState) *XCrypto {
		ud := l.CheckUserData(1)
		if v, ok := ud.Value.(*XCrypto); ok {
			return v
		}

		l.ArgError(1, "users expected")
		return nil
	}

	// create a new crypto instance:
	newCrypto := func(l *lua.LState) int {
		ud := l.NewUserData()
		ud.Value = &XCrypto{
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
