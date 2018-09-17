-- Unit testing starts
require("crypto")
local json = require("json")
local sdk = require("sdk")

function insert(pk, from, path, data)
    -- create the resource pointer:
    ptr = rpointer.new({
        from = from,
        path = path
    })

    -- create the resource:
    res = resource.new({
        pointer = ptr,
        data = data
    })

    -- sign the resource:
    sig = pk:sign(res:hash())

    -- execte the transaction:
    resp = sdk.service().transact({
        resource = res,
        sig = sig
    })

    return resp
end

function delete(pk, from, path)
    -- create the resource pointer:
    ptr = rpointer.new({
        from = from,
        path = path
    })

    -- sign the resource pointer:
    sig = pk:sign(ptr:hash())

    -- execute the transaction:
    resp = sdk.service().transact({
        rpointer = ptr,
        sig = sig
    })

    return resp
end

function retrieve(pk, from, path)
    -- create the resource pointer:
    resPointer = rpointer.new({
        from = from,
        path = path
    })

    -- sign the resource pointer:
    sig = pk:sign(resPointer:hash())

    -- execute the query:
    resp = sdk.service().query({
        rpointer = resPointer,
        sig = sig
    })

    return resp
end

-- create the private key:
pk = privkey.new("60699820e4ba230a2590255678f1ce701f3943e2f7388854ae0af169b9871c0a")
from = pk:pubKey()
path = "/messages"
id = "c52e8cdb-3fb2-4c0a-b4e7-90c677870774"
retrievalPath = path .. "/" .. id
message = {
    id = id,
    title = "This is my message title",
    description = "this is a message description.  Oh yes!"
}

data = json.encode(message)

-- retrieve the resource, should not be found:
retNotFoundResp = retrieve(pk, from, retrievalPath)
assert(retNotFoundResp:code() == 1)
assert(retNotFoundResp:log() == "not found")

-- delete the resource before inserting, not found:
insResp = delete(pk, from, retrievalPath)
assert(insResp:code() == 1)
assert(insResp:log() == "not found")

-- insert the resource:
insResp = insert(pk, from, path, data)
assert(insResp:code() == 0)
assert(insResp:log() == "success")
assert(insResp:gazUsed() == 1205)

-- retrieve the resource:
retResp = retrieve(pk, from, retrievalPath)
retrievedMsg = json.decode(retResp:value())
assert(retResp:code() == 0)
assert(retResp:log() == "success")
assert(retResp:key() == retrievalPath)
assert(retrievedMsg.id == message.id)
assert(retrievedMsg.title == message.title)
assert(retrievedMsg.description == message.description)

-- delete the resource:
delResp = delete(pk, from, retrievalPath)
assert(delResp:code() == 0)
assert(delResp:log() == "success")
