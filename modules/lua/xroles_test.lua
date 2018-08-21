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
        usrs = xusers.load()

        rols = xroles.load()
        amountAdded = rols:add(key, firstPK:pubKey(), secondPK:pubKey(), thirdPK:pubKey())
        amountDel = rols:del(key, secondPK:pubKey())

        firstAmountEnabled = rols:enableWriteAccess(key, usrs:key(firstPK:pubKey()), usrs:key(thirdPK:pubKey()))
        firstAmountDisabled = rols:disableWriteAccess(key, usrs:key(firstPK:pubKey()))

        -- verify:
        assertEquals(type(amountAdded), "number")
        assertEquals(amountAdded, 3)
        assertEquals(amountDel, 1)
        assertEquals(firstAmountEnabled, 2)
        assertEquals(firstAmountDisabled, 1)
    end

-- class TestRoles

LuaUnit:run()
