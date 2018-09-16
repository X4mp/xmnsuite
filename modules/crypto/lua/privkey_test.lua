-- load the modules:
require("crypto")

-- execute:
local x = privkey.new()
local pubKey = x:pubKey()

local anotherPrivKey = privkey.new("a32889104085887e230c406f1b0c5852239b603b5e7e8488ad3e8a3e3f83cdc8cbb7943213b07c81a0cf215c86ea46a697d6dff85b98335733038f2d3a64388306e8603951")
local anotherPubKey = anotherPrivKey:pubKey()

local ringSign = anotherPrivKey:ringSign("mymessage", {pubKey, anotherPubKey})

-- verify:
assert(type(pubKey) == "string")
assert(type(anotherPubKey) == "string")
assert(anotherPubKey == "aab735bdbbe6f03fd37f020738a76c4c92a5f00ec82d69f6c129f97b46112a59")
assert(type(ringSign) == "string")
