-- load the modules:
require("datastore")
require("crypto")

-- variables:
x = privkey.new()
pubKey = x:pubKey()

-- execute:
usrs = users.load()
invalidUsrKey = usrs:key(pubKey)
firstExists = usrs:exists(pubKey)
isInserted = usrs:insert(pubKey)
amountUsers = usrs:len()
isInsertedAgain = usrs:insert(pubKey)
isDeleted = usrs:delete(pubKey)
isDeletedAgain = usrs:delete(pubKey)
secondAmountUsers = usrs:len()

-- verify:
assert(type(invalidUsrKey) == "string")
assert(firstExists == false)
assert(isInserted == true)
assert(amountUsers == 1)
assert(isInsertedAgain == false)
assert(isDeleted == true)
assert(isDeletedAgain == false)
assert(secondAmountUsers == 0)
