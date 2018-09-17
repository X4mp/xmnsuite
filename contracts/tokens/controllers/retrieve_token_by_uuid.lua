local json = require("json")
local uuid = require("uuid")

-- retrieve wallet by its UUID
function retrieveTokenByUUID(from, path, params, sig)
    local rep = Repository:create()
    local token = rep:retrieveTokenByUUID(uuid.new(params.uid))
    if token == null then
        return {
            code = 1,
            log="the token (path: ".. path ..") was not found",
        }
    end

    return {
        code = 0,
        log="success",
        key=path,
        value=json.encode(token:toData())
    }
end
