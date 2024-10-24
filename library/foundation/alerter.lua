---@class alerter
alerter = alerter or {}

---@type Timer
alerter._alerterMessageTimer = alerter._alerterMessageTimer or nil
--- 异步玩家警告提示
---@param whichPlayer Player
---@param vcm boolean 是否播放音效
---@param msg string 警告信息
---@return void
function alerter.message(whichPlayer, vcm, msg)
    if (false == isClass(whichPlayer, PlayerClass)) then
        return
    end
    if (type(vcm) ~= "boolean") then
        vcm = true
    end
    if (type(msg) == "string" and string.len(msg) > 0) then
        if (vcm) then
            async.call(whichPlayer, function()
                audio(Vcm("war3_Error"))
            end)
        end
        async.call(whichPlayer, function()
            UIKit("lik_msg"):alerterMessage(msg)
        end)
        local frames = math.max(3, 0.2 * mbstring.len(msg)) * 60
        async.call(whichPlayer, function()
            if (nil ~= alerter._alerterMessageTimer) then
                destroy(alerter._alerterMessageTimer)
                alerter._alerterMessageTimer = nil
            end
            alerter._alerterMessageTimer = async.setTimeout(frames, function()
                UIKit("lik_msg"):alerterMessage('')
            end)
        end)
    end
end