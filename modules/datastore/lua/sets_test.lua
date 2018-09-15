-- load the modules:
require("datastore")

-- variables:
key = "my-key"
secondKey = "my-second-key"
unionStoreKey = "this-is-the-union-store-key"
walkStoreKey = "this-is-the-walkstore-key"
interStoreKey = "this-is-the-interstore-key"
firstValue = "this is the first value"
secondValue = "this is the second value"

function walkFn(index, value)
    return "works!"
end

-- execute:
sts = sets.load()

-- add:
retAmountAdded = sts:add(key, firstValue, secondValue)
assert(type(retAmountAdded) == "number")
assert(retAmountAdded == 2)

-- del:
retAmountDeleted = sts:del(key, firstValue)
assert(type(retAmountDeleted) == "number")
assert(retAmountDeleted == 1)

-- retrieve:
retValues = sts:retrieve(key, 0, -1)
assert(type(retValues) == "table")
assert(table.getn(retValues) == 1)
assert(retValues[1] == secondValue)

-- len:
length = sts:len(key)
assert(type(length) == "number")
assert(length == 1)

-- union:
sts:add(secondKey, firstValue, secondValue)
retUnion = sts:union(key, secondKey)
assert(type(retUnion) == "table")
assert(table.getn(retUnion) == 2)
assert(retUnion[1] == secondValue)
assert(retUnion[2] == firstValue)

-- unionstore:
retAmountUnionStore = sts:unionstore(unionStoreKey, key, secondKey)
retUnionStoreValues = sts:retrieve(unionStoreKey, 0, -1)
assert(type(retAmountUnionStore) == "number")
assert(retAmountUnionStore == 2)
assert(type(retUnion) == "table")
assert(table.getn(retUnion) == 2)
assert(retUnion[1] == secondValue)
assert(retUnion[2] == firstValue)

-- inter:
retInter = sts:inter(key, secondKey)
assert(type(retInter) == "table")
assert(table.getn(retInter) == 1)
assert(retInter[1] == secondValue)

-- interstore:
retAmountInterStore = sts:interstore(interStoreKey, key, secondKey)
retInterStoreValues = sts:retrieve(interStoreKey, 0, -1)
assert(type(retAmountInterStore) == "number")
assert(retAmountInterStore == 1)
assert(type(retInterStoreValues) == "table")
assert(table.getn(retInterStoreValues) == 1)
assert(retInterStoreValues[1] == secondValue)

-- trim:
retRemaining = sts:trim(secondKey, 1, 1)
assert(type(retInter) == "table")
assert(table.getn(retInter) == 1)
assert(retInter[1] == secondValue)

-- walk:
retWalk = sts:walk(key, walkFn)
assert(type(retWalk) == "table")
assert(table.getn(retWalk) == 1)
assert(retWalk[1] == "works!")

-- walkstore
retAmountWalk = sts:walkstore(walkStoreKey, key, walkFn)
retRetrieveWalkValues = sts:retrieve(walkStoreKey, 0, -1)
assert(type(retAmountWalk) == "number")
assert(retAmountWalk == 1)
assert(type(retRetrieveWalkValues) == "table")
assert(table.getn(retRetrieveWalkValues) == 1)
assert(retRetrieveWalkValues[1] == "works!")
