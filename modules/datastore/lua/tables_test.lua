-- load the modules:
require("datastore")

-- variables:
firstKey = "steve-rodrigue"
firstTable = {firstName="Steve", lastName=Rodrigue, id={type="person", number="q24234"}}
firstParam = {key=firstKey, table=firstTable}

secondKey = "roger-cyr"
secondTable = {firstName="Roger", lastName="Cyr"}
secondParam = {key=secondKey, table=secondTable}

-- execute:
x = tables.load()
reload = tables.load()
retAmountSaved = x:save(firstParam, secondParam)
retFirstObj = x:retrieve(firstKey)
retSecondObj = x:retrieve(secondKey)
retInvalidObj = x:retrieve("invalidkey")
retLen = x:len()
retAmountExists = x:exists(firstKey, "invalid-key", secondKey)
retAmountDel = reload:delete(firstKey, "invalid")

-- verify:
assert(retAmountSaved == 2)
assert(retFirstObj.firstName == firstTable.firstName)
assert(retFirstObj.lastName == firstTable.lastName)
assert(retFirstObj.id.type == firstTable.id.type)
assert(retFirstObj.id.number == firstTable.id.number)

assert(retSecondObj.firstName == secondTable.firstName)
assert(retSecondObj.lastName == secondTable.lastName)

assert(retInvalidObj == null)
assert(retLen == 2)
assert(retAmountExists == 2)
assert(retAmountDel == 1)
