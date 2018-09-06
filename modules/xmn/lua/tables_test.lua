-- Unit testing starts
require('luaunit')
local xmn = require("xmn")

TestTables = {} --class
    function TestTables:testSaveThenRetrieve_Success()
        -- variables:
        firstKey = "steve-rodrigue"
        firstTable = {firstName="Steve", lastName=Rodrigue, id={type="person", number="q24234"}}
        firstParam = {key=firstKey, table=firstTable}

        secondKey = "roger-cyr"
        secondTable = {firstName="Roger", lastName="Cyr"}
        secondParam = {key=secondKey, table=secondTable}

        -- execute:
        x = tables.load()
        reload = tables.load()
        retAmountSaved = x:save(firstParam, secondParam)
        retFirstObj = x:retrieve(firstKey)
        retSecondObj = x:retrieve(secondKey)
        retInvalidObj = x:retrieve("invalidkey")
        retLen = x:len()
        retAmountExists = x:exists(firstKey, "invalid-key", secondKey)
        retAmountDel = reload:delete(firstKey, "invalid")

        -- verify:
        assertEquals(retAmountSaved, 2)
        assertEquals(retFirstObj, firstTable)
        assertEquals(retSecondObj, secondTable)
        assertEquals(retInvalidObj, null)
        assertEquals(retLen, 2)
        assertEquals(retAmountExists, 2)
        assertEquals(retAmountDel, 1)
    end

-- class TestTables

LuaUnit:run()
