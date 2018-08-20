-- Unit testing starts
require('luaunit')

TestUsers = {} --class
    function TestUsers:testSave_Keys_Success()
        -- variables:
        x = xcrypto.new()
        pubKey = x:pubKey()

        -- execute:
        usrs = xusers.load()
        invalidUsrKey = usrs:key(pubKey)
        exists = usrs:exists(pubKey)

        -- verify:
        assertEquals(type(invalidUsrKey), "string")
        assertEquals(exists, false)
    end

-- class TestUsers

LuaUnit:run()
