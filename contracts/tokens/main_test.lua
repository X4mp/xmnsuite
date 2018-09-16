require("crypto")
local sdk = require("sdk")
local json = require("json")

function save(fromPK, path, obj)
    -- create the resource pointer:
    ptr = rpointer.new({
        from = fromPK:pubKey(),
        path = path
    })

    -- create the resource:
    res = resource.new({
        pointer = ptr,
        data = json.encode(obj)
    })

    -- sign the resource:
    sig = fromPK:sign(res:hash())

    -- execte the transaction:
    resp = sdk.service().transact({
        resource = res,
        sig = sig
    })

    return resp
end

function retrieve(fromPK, path)

    -- create the resource pointer:
    ptr = rpointer.new({
        from = fromPK:pubKey(),
        path = path
    })

    -- sign the resource:
    sig = fromPK:sign(ptr:hash())

    -- execte the transaction:
    resp = sdk.service().query({
        rpointer = ptr,
        sig = sig
    })

    return resp
end

function delete(fromPK, path)

    -- create the resource pointer:
    ptr = rpointer.new({
        from = fromPK:pubKey(),
        path = path
    })

    -- sign the resource:
    sig = fromPK:sign(ptr:hash())

    -- execte the transaction:
    resp = sdk.service().transact({
        rpointer = ptr,
        sig = sig
    })

    return resp
end

-- create the node pk:
local nodePK = privkey.new("a3288910407b8954f7afab990ce902a33c290a73fd399fec9d96e1ff221826ac05e88c7a56bfbac5062c2b21060d94af33e72fcc64bcd599f3af596c96d859b31418adb0de")

-- generate a new wallet pk:
local walletPK = privkey.new()

-- create a new wallet:
local wallet = {
    pub_key = walletPK:pubKey(),
    created_on = os.time()
}

-- insert a new wallet:
insWalletResp = save(nodePK, "/wallets", wallet)
assert(insWalletResp:code() == 0)
assert(insWalletResp:log() == "success")
assert(insWalletResp:gazUsed() == 0)

-- retrieve a wallet:
retResp = retrieve(nodePK, "/wallets/" .. wallet.pub_key)
assert(retResp:code() == 0)
assert(retResp:key() == "/wallets/" .. wallet.pub_key)
assert(retResp:log() == "success")

retWallet = json.decode(retResp:value())
assert(retWallet.pub_key == wallet.pub_key)
assert(tonumber(retWallet.created_on) == wallet.created_on)

-- create a new token
local tok = {
    id = "24f1267a-931d-4560-8a10-650b0b83d81a",
    amount = math.pow(2, 62),
    wallet = wallet.pub_key,
    symbol = "XMND",
    name = "XMN Dollars",
    description = "This is the XMN dollars",
    created_on = os.time()
}

-- save the new token:
insTokResp = save(nodePK, "/tokens", tok)
assert(insTokResp:code() == 0)
assert(insTokResp:log() == "success")
assert(insTokResp:gazUsed() == 0)

-- retrieve a token:
retTokResp = retrieve(nodePK, "/tokens/" .. tok.id)
assert(retTokResp:code() == 0)
assert(retTokResp:key() == "/tokens/" .. tok.id)
assert(retTokResp:log() == "success")

retToken = json.decode(retTokResp:value())
assert(retToken.id == tok.id)
assert(retToken.symbol == tok.symbol)
assert(retToken.name == tok.name)
assert(retToken.description == tok.description)
assert(tonumber(retToken.created_on) == tok.created_on)

-- retrieve the deposits related to the wallet:

-- retrieve the deposits related to the token:

-- delete a token:
delResp = delete(nodePK, "/tokens/" .. tok.id)
assert(delResp:code() == 0)
assert(delResp:log() == "success")

-- delete a wallet:
delResp = delete(nodePK, "/wallets/" .. wallet.pub_key)
assert(delResp:code() == 0)
assert(delResp:log() == "success")
