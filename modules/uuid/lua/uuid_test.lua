-- load the modules:
local uuid = require("uuid")

-- execute:
id = uuid.new()

-- string:
idAsString = id:string()
assert(type(idAsString) == "string")
assert(idAsString == string.match(idAsString, "%x%x%x%x%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%x%x%x%x%x%x%x%x"))
