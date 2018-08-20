-- Unit testing starts
require('luaunit')

TestObjects = {} --class
    function TestObjects:testSaveThenRetrieve_Success()
        -- variables:
        firstKey = "steve-rodrigue"
        firstObject = {}
        firstObject.firstName = "Steve"
        firstObject.lastName = "Rodrigue"
        firstObject.id = {}
        firstObject.id.type = "person"
        firstObject.id.number = "q24234"
        firstParam = {}
        firstParam.key = firstKey
        firstParam.object = firstObject

        secondKey = "roger-cyr"
        secondObject = {}
        secondObject.firstName = "Roger"
        secondObject.lastName = "Cyr"
        secondParam = {}
        secondParam.key = secondKey
        secondParam.object = secondObject

        -- execute:
        x = xobjects.load()
        retAmountSaved = x:save(firstParam, secondParam)
        retFirstObj = x:retrieve(firstKey)
        retSecondObj = x:retrieve(secondKey)

        -- verify:
        assertEquals(retAmountSaved, 2)
        assertEquals(retFirstObj, firstObject)
        assertEquals(retSecondObj, secondObject)
    end

-- class TestObjects

LuaUnit:run()
