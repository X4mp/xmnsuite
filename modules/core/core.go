package core

import lua "github.com/yuin/gopher-lua"

// Register registers the base module
func Register(l *lua.LState) {
	// crypto:
	CreateXCrypto(l)

	// datastore:
	CreateXKeys(l)
	CreateXTables(l)

	// roles and users:
	CreateXUsers(l)
	CreateXRoles(l)

	// router:
	CreateXRoute(l)
	CreateXRouter(l)
}
