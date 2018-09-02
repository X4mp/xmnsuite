-- Unit testing starts
require('luaunit')

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
