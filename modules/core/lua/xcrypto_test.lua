-- Unit testing starts
require('luaunit')

TestCrypto = {} --class
    function TestCrypto:testGeneratePK_pubKey_Success()
        -- execute:
        x = xcrypto.new()
        pubKey = x:pubKey()

        -- verify:
        assertEquals(type(pubKey), "string")
    end

-- class TestCrypto

LuaUnit:run()
