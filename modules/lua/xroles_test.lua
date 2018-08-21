-- Unit testing starts
require('luaunit')

TestRoles = {} --class
    function TestRoles:testAdd_Success()
        -- variables:
        key = "this-is-my-role-key"
        firstPK = xcrypto.new()
        secondPK = xcrypto.new()
        thirdPK = xcrypto.new()

        -- execute:
        rols = xroles.load()
        amountAdded = rols:add(key, firstPK:pubKey(), secondPK:pubKey(), thirdPK:pubKey())
        amountDel = rols:del(key, secondPK:pubKey())

        -- verify:
        assertEquals(type(amountAdded), "number")
        assertEquals(amountAdded, 3)
        assertEquals(amountDel, 1)
    end

-- class TestRoles

LuaUnit:run()
