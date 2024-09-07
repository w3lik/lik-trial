---@class alerter
alerter = alerter or {}

--- 异步玩家警告提示
---@param whichPlayer Player
---@param vcm boolean 是否播放音效
---@param msg string 警告信息
---@param red number 红0-255
---@param green number 绿0-255
---@param blue number 蓝0-255
---@param alpha number 透明0-255
---@return void
function alerter.message(whichPlayer, vcm, msg, red, green, blue, alpha)
    if (false == isClass(whichPlayer, PlayerClass)) then
        return
    end
    async.call(whichPlayer, function()
        if (type(vcm) ~= "boolean") then
            vcm = true
        end
        if (type(msg) == "string" and string.len(msg) > 0) then
            if (vcm) then
                audio(Vcm("war3_Error"))
            end
            -- 默认金色
            alpha = alpha or 255
            red = red or 255
            green = green or 215
            blue = blue or 0
            local dur = math.max(3, 0.2 * mbstring.len(msg))
            japi.DZ_SimpleMessageFrameAddMessage(japi.DZ_FrameGetWorldFrameMessage(), msg, japi.DZ_GetColor(alpha, red, green, blue), dur, false)
        end
    end)
end