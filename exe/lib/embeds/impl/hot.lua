function FWHT()
    local create = function(name)
        if package.loaded[name] ~= nil then
            return
        end
        local status, err = xpcall(require, debug.traceback, name)
        if (status ~= true) then
            print("[HOT]err: " .. err)
        else
            print("[HOT]load: " .. name)
        end
    end
    local reload = function(name)
        if package.loaded[name] then
            package.loaded[name] = nil
        end
        local status, err = xpcall(require, debug.traceback, name)
        if (status ~= true) then
            print("[HOT]err: " .. err)
        else
            print("[HOT]reload: " .. name)
        end
    end
    local fwhc = function()
        local cr, err1 = io.open("fwhc.txt", 'r+')
        if (err1 == nil) then
            local str = cr:read('a')
            cr:close()
            if (string.len(str) > 0) then
                local strs = string.explode('|', str)
                for _, s in ipairs(strs) do
                    create(s)
                end
            end
        end
    end
    local fwht = function()
        local rl, err2 = io.open("fwht.txt", 'r+')
        if (err2 == nil) then
            local str = rl:read('a')
            rl:close()
            rl, _ = io.open("fwht.txt", 'w')
            rl:close()
            if (string.len(str) > 0) then
                local strs = string.explode('|', str)
                for _, s in ipairs(strs) do
                    reload(s)
                end
            end
        end
    end
    fwhc()
    time.setInterval(0.5, function()
        fwhc()
        fwht()
    end)
end
Game():onEvent(EVENT.Game.Start, "FWHT", function() FWHT() end)