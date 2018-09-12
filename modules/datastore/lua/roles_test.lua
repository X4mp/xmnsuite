-- load the modules:
require("datastore")
require("crypto")

-- variables:
key = "this-is-my-role-key"
onKey = "this-is-a-key-to-enable-write-access-on"
firstPK = privkey.new()
secondPK = privkey.new()
thirdPK = privkey.new()

-- execute:
usrs = users.load()

rols = roles.load()
amountAdded = rols:add(key, firstPK:pubKey(), secondPK:pubKey(), thirdPK:pubKey())
amountDel = rols:del(key, secondPK:pubKey())

firstAmountEnabled = rols:enableWriteAccess(onKey, usrs:key(firstPK:pubKey()), usrs:key(thirdPK:pubKey()))
firstAmountWriteAccess = rols:hasWriteAccess(onKey, usrs:key(firstPK:pubKey()), usrs:key(secondPK:pubKey()), usrs:key(thirdPK:pubKey()))
firstAmountDisabled = rols:disableWriteAccess(onKey, usrs:key(firstPK:pubKey()))
secondAmountWriteAccess = rols:hasWriteAccess(onKey, usrs:key(firstPK:pubKey()), usrs:key(secondPK:pubKey()), usrs:key(thirdPK:pubKey()))

-- verify:
assert(type(amountAdded) == "number")
assert(amountAdded == 3)
assert(amountDel == 1)
assert(firstAmountEnabled == 2)

assert(firstAmountWriteAccess[1] == usrs:key(firstPK:pubKey()))
assert(firstAmountWriteAccess[2] == usrs:key(thirdPK:pubKey()))

assert(firstAmountDisabled == 1)

assert(secondAmountWriteAccess[1] == usrs:key(thirdPK:pubKey()))
