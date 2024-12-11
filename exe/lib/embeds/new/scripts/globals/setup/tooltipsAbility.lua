---@param whichAbility AbilityTpl|Ability
---@return table
function tooltipsAbility(whichAbility, lvOffset)
    if (instanceof(whichAbility, AbilityFuncClass) == false) then
        return nil
    end
    local types = {
        --{ "worth" },
        { "coolDown", "15DF89", "秒" },
        { "hpCost", "DE5D43", "血" },
        { "mpCost", "83B3E4", "蓝" },
    }
    lvOffset = lvOffset or 0
    local lv = lvOffset + whichAbility:level()
    if (lv > whichAbility:levelMax()) then
        return nil
    end
    local tips
    if (lvOffset > 0) then
        tips = Game():combineDescription(whichAbility, { level = lv }, "abilityBase", SYMBOL_D, "attributes", "abilityLvPoint")
    else
        tips = Game():combineDescription(whichAbility, nil, "abilityBase", SYMBOL_D, "attributes")
    end
    local content = {
        icons = {},
        bars = {},
        tips = tips,
    }
    --for _, c in ipairs(types) do
    --    local method = c[1]
    --    if (method == "worth") then
    --        local wv = whichAbility:worthCost(lv)
    --        if (wv ~= nil) then
    --            local wk = {
    --                { "lumber", "C49D5A", "木" },
    --                { "gold", "ECD104", "金" },
    --                { "silver", "E3E3E3", "银" },
    --                { "copper", "EC6700", "铜" }
    --            }
    --            for _, w in ipairs(wk) do
    --                local key = w[1]
    --                local color = w[2]
    --                local uit = w[3]
    --                local val = math.floor(wv[key] or 0)
    --                if (val > 0) then
    --                    if (uit ~= nil) then
    --                        val = val .. " " .. uit
    --                    end
    --                    table.insert(content.icons, {
    --                        texture = "Framework\\ui\\" .. key .. ".tga",
    --                        text = colour.hex(color, val),
    --                    })
    --                end
    --            end
    --        end
    --    end
    --end
    for _, c in ipairs(types) do
        local method = c[1]
        if (method == "coolDown" or method == "hpCost" or method == "mpCost") then
            local color = c[2]
            local uit = c[3]
            local val = 0
            if (method == "coolDown") then
                val = whichAbility:coolDown(lv)
            elseif (method == "hpCost") then
                val = whichAbility:hpCost(lv)
            elseif (method == "mpCost") then
                val = whichAbility:mpCost(lv)
            end
            if (val > 0) then
                if (uit ~= nil) then
                    val = val .. " " .. uit
                end
                table.insert(content.icons, {
                    texture = "Framework\\ui\\" .. method .. ".tga",
                    text = colour.hex(color, val),
                })
            end
        end
    end
    if (isClass(whichAbility, AbilityClass)) then
        if (lv == whichAbility:level() and lv < whichAbility:levelMax()) then
            if (whichAbility:exp() > 0) then
                local cur = whichAbility:exp() or 0
                local prev = whichAbility:expNeed(lv) or 0
                local need = whichAbility:expNeed() or 0
                local percent = math.trunc((cur - prev) / (need - prev), 3)
                if (percent ~= nil) then
                    table.insert(content.bars, {
                        texture = "Framework\\ui\\tile_yellow.tga",
                        text = colour.hex("E2C306", "经验：" .. math.floor(cur - prev) .. "/" .. math.ceil(need - prev)),
                        ratio = percent,
                        width = 0.10,
                        height = 0.001,
                    })
                end
            end
        end
    end
    return content
end