---@param prevUnit Unit
---@param sourceUnit Unit
---@param targetUnit Unit
local _call = function(index, lightningType, sourceUnit, prevUnit, targetUnit, dmg, dmgSrc, dmgType, dmgTypeLv, breakArmor)
    local sp
    if (isClass(prevUnit, UnitClass)) then
        sp = { prevUnit:x(), prevUnit:y(), prevUnit:h() }
    elseif (isClass(sourceUnit, UnitClass)) then
        sp = { sourceUnit:x(), sourceUnit:y(), sourceUnit:h() }
    else
        sp = { targetUnit:x(), targetUnit:y(), targetUnit:h() + 1000 } --头顶劈下
    end
    lightning(lightningType, sp[1], sp[2], sp[3], targetUnit:x(), targetUnit:y(), targetUnit:h(), 0.25)
    sp = nil
    ability.damage({
        sourceUnit = sourceUnit,
        targetUnit = targetUnit,
        damage = dmg,
        damageSrc = dmgSrc,
        damageType = dmgType,
        damageTypeLevel = dmgTypeLv,
        breakArmor = breakArmor,
    })
    if (sourceUnit ~= nil) then
        event.syncTrigger(sourceUnit, EVENT.Unit.LightningChain, { targetUnit = targetUnit, index = index })
    end
    event.syncTrigger(targetUnit, EVENT.Unit.Be.LightningChain, { sourceUnit = sourceUnit, index = index })
end

--[[
    模版技能 闪电链
    options = {
        sourceUnit 伤害来源
        targetUnit 目标单位
        lightningType = LIGHTNING_TYPE, -- 闪电效果类型(可选 详情查看 LIGHTNING_TYPE)
        qty = 1, --传递的最大单位数（可选，默认1）
        rate = 0, --增减率%（可选，默认不增不减为0，范围建议[-100,100]）
        radius = 500, --寻找下一目标的作用半径范围（可选，默认500）
        isRepeat = false, --是否允许同一个单位重复打击（最近2次打击不会是同一个，repeat也不能打击同一个单体单位多次）
        damage 伤害
        damageSrc 伤害来源 DAMAGE_SRC
        damageType 伤害类型 DAMAGE_TYPE
        damageTypeLevel 伤害等级默认1
        breakArmor 破防类型 BREAK_ARMOR
        prevUnit = [unit], --隐藏的参数，上一个的目标单位（必须有，用于构建两点间闪电特效）
        index = 1,--隐藏的参数，用于暗地里记录是第几个被电到的单位
        repeatGroup = [group],--隐藏的参数，用于暗地里记录单位是否被电过
    }
]]
---@param options {sourceUnit:Unit,targetUnit:Unit,qty:number,rate:number,radius:number,damage:number,damageSrc:table,damageType:table,damageTypeLevel:number,breakArmor:table}
function ability.lightningChain(options)
    local sourceUnit = options.sourceUnit
    local targetUnit = options.targetUnit
    if (isClass(targetUnit, UnitClass) == false or targetUnit:isDead()) then
        return
    end
    local damage = options.damage or 0
    if (damage > 0) then
        local lightningType = options.lightningType or LIGHTNING_TYPE.thunder
        local qty = options.qty or 1
        local rate = 100
        local dmgSrc = options.damageSrc or DAMAGE_SRC.ability
        local dmgType = options.damageType or DAMAGE_TYPE.common
        local dmgTypeLv = options.damageTypeLevel or 1
        local breakArmor = options.breakArmor or { BREAK_ARMOR.avoid }
        local index = 1
        _call(index, lightningType, sourceUnit, nil, targetUnit, damage * rate * 0.01, dmgSrc, dmgType, dmgTypeLv, breakArmor)
        qty = qty - 1
        if (qty > 0) then
            local radius = options.radius or 600
            local isRepeat = options.isRepeat or false
            local repeatJudge = {}
            time.setInterval(0.25, function(curTimer)
                if (qty <= 0 or isDestroy(targetUnit)) then
                    destroy(curTimer)
                    return
                end
                if (options.rate) then
                    rate = rate + options.rate
                end
                index = index + 1
                local nextUnit = Group():closest(UnitClass, {
                    circle = {
                        x = targetUnit:x(),
                        y = targetUnit:y(),
                        radius = radius,
                    },
                    ---@param enumUnit Unit
                    filter = function(enumUnit)
                        if (repeatJudge[enumUnit:id()] ~= nil) then
                            return false
                        end
                        if (isDestroy(enumUnit) or enumUnit:isDead()) then
                            return false
                        end
                        if (sourceUnit == nil or isDestroy(sourceUnit)) then
                            return enumUnit:isAlly(targetUnit:owner())
                        else
                            return enumUnit:isEnemy(sourceUnit:owner())
                        end
                    end,
                })
                if (nextUnit == nil) then
                    destroy(curTimer)
                    return
                end
                if (isRepeat ~= true) then
                    repeatJudge[nextUnit:id()] = 1
                end
                _call(index, lightningType, sourceUnit, targetUnit, nextUnit, damage * rate * 0.01, dmgSrc, dmgType, dmgTypeLv, breakArmor)
                qty = qty - 1
                targetUnit = nextUnit
            end)
        end
    end
end