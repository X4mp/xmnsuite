function generateKeyname(entityName, pkName, pkValue)
    return entityName .. ":" .. "by_" .. pkName .. ":" .. pkValue
end

function toData(tb)
    data = {}
    for index, obj in pairs(tb) do
        oneData = obj:toData()
        table.insert(data, oneData)
    end

    return data
end
