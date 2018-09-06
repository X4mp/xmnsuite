-- Unit testing starts
require('luaunit')
local xmn = require("xmn")

TestPrivKey = {} --class
    function TestPrivKey:testGeneratePK_pubKey_Success()
        -- execute:
        x = privkey.new()
        pubKey = x:pubKey()

        -- verify:
        assertEquals(type(pubKey), "string")
    end

-- class TestPrivKey

LuaUnit:run()
