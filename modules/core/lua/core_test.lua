-- Unit testing starts
require('luaunit')

TestCore = {} --class
    function TestCore:testCreate_Success()
        -- func handlers:
        function saveVideo(from, path, params, data, sig)
            return {code = 0, log="video saved"}
        end

        function saveProfile(from, path, params, data, sig)
            keyname = "profile-path:"..path

            return {
                code = 0,
                log="profile saved",
                gazUsed=0.675,
                tags={
                    keyname = "someValue"
                }
            }
        end

        -- create the application
        first = xapp.new("17.03.09", 0, -1, xrouter.new(
            xroute.new("save", "/videos", saveVideo),
            xroute.new("save", "/profiles", saveProfile)
        ))

        xchain.load(first)

        -- verify:
        -- assertEquals(type(firstRte), "xroute")
    end

-- class TestCore

LuaUnit:run()
