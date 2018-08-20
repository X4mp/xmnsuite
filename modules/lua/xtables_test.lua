-- Unit testing starts
require('luaunit')

TestTables = {} --class
    function TestTables:testSaveThenRetrieve_Success()
        -- variables:
        firstKey = "steve-rodrigue"
        firstTable = {}
        firstTable.firstName = "Steve"
        firstTable.lastName = "Rodrigue"
        firstTable.id = {}
        firstTable.id.type = "person"
        firstTable.id.number = "q24234"
        firstParam = {}
        firstParam.key = firstKey
        firstParam.table = firstTable

        secondKey = "roger-cyr"
        secondTable = {}
        secondTable.firstName = "Roger"
        secondTable.lastName = "Cyr"
        secondParam = {}
        secondParam.key = secondKey
        secondParam.table = secondTable

        -- execute:
        x = xtables.load()
        retAmountSaved = x:save(firstParam, secondParam)
        retFirstObj = x:retrieve(firstKey)
        retSecondObj = x:retrieve(secondKey)
        retLen = x:len()
        retAmountExists = x:exists(firstKey, "invalid-key", secondKey)
        retAmountDel = x:delete(firstKey, "invalid")

        -- verify:
        assertEquals(retAmountSaved, 2)
        assertEquals(retFirstObj, firstTable)
        assertEquals(retSecondObj, secondTable)
        assertEquals(retLen, 2)
        assertEquals(retAmountExists, 2)
        assertEquals(retAmountDel, 1)
    end

-- class TestTables

LuaUnit:run()
