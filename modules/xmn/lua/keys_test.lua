-- Unit testing starts
require('luaunit')
local xmn = require("xmn")

TestKeys = {} --class
    function TestKeys:testSaveThenRetrieve_Success()
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
        assertEquals(type(retValue), "string")
        assertEquals(retValue, value)
        assertEquals(retLen, 1)
        assertEquals(retSearch[0], key)
        assertEquals(retExistsAmount, 1)
        assertEquals(retAmountDeleted, 1)
        assertEquals(retAfterDel, nil)
    end

-- class TestKeys

LuaUnit:run()
