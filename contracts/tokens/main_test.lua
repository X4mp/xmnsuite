require("crypto")
local sdk = require("sdk")
local json = require("json")
local uuid = require("uuid")

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

-- create the user pk:
local userPK = privkey.new("dfab9ff67f646eb235e4aa9b0474a8ba6a987cf905f1e5a04ac3e4d99168730f")

-- generate a new wallet pk:
local walletPK = privkey.new()

-- create a new wallet:
local wallet = {
    pub_key = walletPK:pubKey(),
    created_on = os.time()
}

-- insert a new wallet:
insWalletResp = save(userPK, "/wallets", wallet)
assert(insWalletResp:code() == 0)
assert(insWalletResp:log() == "success")
assert(insWalletResp:gazUsed() == 0)

-- retrieve a wallet:
retResp = retrieve(userPK, "/wallets/" .. wallet.pub_key)
assert(retResp:code() == 0)
assert(retResp:key() == "/wallets/" .. wallet.pub_key)
assert(retResp:log() == "success")

retWallet = json.decode(retResp:value())
assert(retWallet.pub_key == wallet.pub_key)
assert(tonumber(retWallet.created_on) == wallet.created_on)

-- create a new token
local tok = {
    uid = uuid.new():string(),
    amount = math.pow(2, 62),
    wallet = wallet.pub_key,
    symbol = "XMND",
    name = "XMN Dollars",
    description = "This is the XMN dollars",
    created_on = os.time()
}

-- save the new token:
insTokResp = save(userPK, "/tokens", tok)
assert(insTokResp:code() == 0)
assert(insTokResp:log() == "success")
assert(insTokResp:gazUsed() == 0)

-- retrieve a token:
retTokResp = retrieve(userPK, "/tokens/" .. tok.uid)
assert(retTokResp:code() == 0)
assert(retTokResp:key() == "/tokens/" .. tok.uid)
assert(retTokResp:log() == "success")

retToken = json.decode(retTokResp:value())
assert(retToken.uid == tok.uid)
assert(retToken.symbol == tok.symbol)
assert(retToken.name == tok.name)
assert(retToken.description == tok.description)
assert(tonumber(retToken.created_on) == tok.created_on)

-- retrieve the deposits related to the wallet:
retDepsByWalletPubKeyResp = retrieve(userPK, "/deposits/" .. wallet.pub_key)
assert(retDepsByWalletPubKeyResp:code() == 0)
assert(retDepsByWalletPubKeyResp:log() == "success")
assert(retDepsByWalletPubKeyResp:key() == "/deposits/" .. wallet.pub_key)

retDepsByWalletPubKey = json.decode(retDepsByWalletPubKeyResp:value())
assert(table.getn(retDepsByWalletPubKey) == 1)
assert(retDepsByWalletPubKey[1].wallet_pub_key == wallet.pub_key)
assert(retDepsByWalletPubKey[1].token_uuid == tok.uid)
assert(tonumber(retDepsByWalletPubKey[1].amount) == tok.amount)
assert(tonumber(retDepsByWalletPubKey[1].created_on) == tok.created_on)

-- retrieve the deposits related to the token:
retDepsByTokenUUIDResp = retrieve(userPK, "/deposits/" .. tok.uid)
assert(retDepsByTokenUUIDResp:code() == 0)
assert(retDepsByTokenUUIDResp:log() == "success")
assert(retDepsByTokenUUIDResp:key() == "/deposits/" .. tok.uid)

retDepsByTokenUUID = json.decode(retDepsByTokenUUIDResp:value())
assert(table.getn(retDepsByTokenUUID) == 1)
assert(retDepsByTokenUUID[1].wallet_pub_key == wallet.pub_key)
assert(retDepsByTokenUUID[1].token_uuid == tok.uid)
assert(tonumber(retDepsByTokenUUID[1].amount) == tok.amount)
assert(tonumber(retDepsByTokenUUID[1].created_on) == tok.created_on)

-- retrieve the deposits related to the token AND wallet pub key:
retDepsByWalletAndTokenUUIDResp = retrieve(userPK, "/deposits/" .. wallet.pub_key .. "/" .. tok.uid)
assert(retDepsByWalletAndTokenUUIDResp:code() == 0)
assert(retDepsByWalletAndTokenUUIDResp:log() == "success")
assert(retDepsByWalletAndTokenUUIDResp:key() == "/deposits/" .. wallet.pub_key .. "/" .. tok.uid)

retDepsByWalletAndTokenUUID = json.decode(retDepsByWalletAndTokenUUIDResp:value())
assert(table.getn(retDepsByWalletAndTokenUUID) == 1)
assert(retDepsByWalletAndTokenUUID[1].wallet_pub_key == wallet.pub_key)
assert(retDepsByWalletAndTokenUUID[1].token_uuid == tok.uid)
assert(tonumber(retDepsByWalletAndTokenUUID[1].amount) == tok.amount)
assert(tonumber(retDepsByWalletAndTokenUUID[1].created_on) == tok.created_on)

-- delete a token:
--delResp = delete(userPK, "/tokens/" .. tok.uid)
--assert(delResp:code() == 0)
--assert(delResp:log() == "success")

-- delete a wallet:
--delResp = delete(userPK, "/wallets/" .. wallet.pub_key)
--assert(delResp:code() == 0)
--assert(delResp:log() == "success")
