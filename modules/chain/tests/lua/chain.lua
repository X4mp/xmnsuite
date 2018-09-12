-- load the modules:
require("datastore")
local json = require("json")
local chain = require("chain")

-- func handlers:
function saveMessage(from, path, params, data, sig)
    local newMsg = json.decode(data)
    local msgPath = path .. "/" .. newMsg.id

    -- save the msg:
    local x = tables.load()
    local retAmountSaved = x:save({key=msgPath, table=newMsg})
    if retAmountSaved ~= 1 then
        return {
            code = 2,
            log="there was an error while saving the message",
        }
    end

    return {
        code = 0,
        log="success",
        gazUsed=1205,
        tags={
            {
                key=msgPath,
                value=json.encode(newMsg)
            }
        }
    }
end

function retrieveMessageByID(from, path, params, sig)
    local x = tables.load()
    local msg = x:retrieve(path)
    if msg == null then
        return {
            code = 1,
            log="not found",
            key=path,
            value=null
        }
    end

    return {
        code = 0,
        log="success",
        key=path,
        value=json.encode(msg)
    }
end

function deleteMessageByID(from, path, params, sig)
    local x = tables.load()
    local retAmountDeleted = x:delete(path)
    if retAmountDeleted ~= 1 then
        return {
            code = 1,
            log="not found",
        }
    end

    return {
        code = 0,
        log="success",
    }
end

chain.chain().load({
    namespace = "xmn",
    name = "messages",
    apps = {
        chain.app().new({
            version = "17.03.09",
            beginBlockIndex = 0,
            endBlockIndex = -1,
            router = chain.router().new({
                key = "this-is-the-router-key",
                routes = {
                    chain.route().new("save", "/messages", saveMessage),
                    chain.route().new("delete", "/messages/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteMessageByID),
                    chain.route().new("retrieve", "/messages/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveMessageByID),
                }
            })
        })
    }
})
