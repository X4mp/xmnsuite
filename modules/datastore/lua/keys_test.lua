-- load the modules:
require("datastore")

-- variables:
key = "my-key"
value = "this is some value"

-- execute:
x = keys.load()
x:save(key, value)
retValue = x:retrieve(key)
retLen = x:len()
retSearch = x:search("[a-z-]+")
retExistsAmount = x:exists(key, "invalid", "another-invalid")
retAmountDeleted = x:delete(key, "another-invalid-key")
retAfterDel = x:retrieve(key)

-- verify:
assert(type(retValue) == "string")
assert(retValue == value)
assert(retLen == 1)
assert(retSearch[0] == key)
assert(retExistsAmount == 1)
assert(retAmountDeleted == 1)
assert(retAfterDel == null)
