-- Unit testing starts
require('luaunit')
local xmn = require("xmn")

TestUsers = {} --class
    function TestUsers:testSave_Keys_Success()
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
        assertEquals(type(invalidUsrKey), "string")
        assertEquals(firstExists, false)
        assertEquals(isInserted, true)
        assertEquals(amountUsers, 1)
        assertEquals(isInsertedAgain, false)
        assertEquals(isDeleted, true)
        assertEquals(isDeletedAgain, false)
        assertEquals(secondAmountUsers, 0)
    end

-- class TestUsers

LuaUnit:run()
