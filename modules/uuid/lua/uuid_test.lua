-- load the modules:
local uuid = require("uuid")

-- execute:
id = uuid.new()
existingID = uuid.new("e2a07177-5f96-4c66-a997-11284fb3cb00")

-- string:
idAsString = id:string()
assert(type(idAsString) == "string")
assert(idAsString == string.match(idAsString, "%x%x%x%x%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%-%x%x%x%x%x%x%x%x%x%x%x%x"))
assert(existingID:string() == "e2a07177-5f96-4c66-a997-11284fb3cb00")
