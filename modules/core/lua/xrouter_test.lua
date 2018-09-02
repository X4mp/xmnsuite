-- Unit testing starts
require('luaunit')

testRouter = {} --class
    function testRouter:testCreate_Success()
        -- func handlers:
        function retrieveVideosByID(from, path, params, sig)
            return {}
        end

        function saveVideo(from, path, params, data, sig)
            return {}
        end

        function deleteVideo(from, path, params, sig)
            return {}
        end

        -- create router:
        xrouter.load(
            xroute.new("retrieve", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", retrieveVideosByID),
            xroute.new("delete", "/videos/<id|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}>", deleteVideo),
            xroute.new("save", "/videos", saveVideo)
        )

        -- verify:
        -- assertEquals(type(firstRte), "xroute")
    end

-- class testRouter

LuaUnit:run()
