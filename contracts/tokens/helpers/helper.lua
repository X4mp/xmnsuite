function generateKeyname(entityName, pkName, pkValue)
    return entityName .. ":" .. "by_" .. pkName .. ":" .. pkValue
end
