---@param hurtData noteOnUnitHurtData
---@param effectiveData noteOnAbilityEffectiveData
TPL_ABILITY.ZZJY = AbilityTpl()
    :name("自在极意被动")
    :targetType(ABILITY_TARGET_TYPE.pas)
    :icon("UI\\Widgets\\Console\\Human\\infocard-heroattributes-str.blp")
    :coolDownAdv(10, 0)
    :mpCostAdv(100, 0)
    :levelMax(10)
    :levelUpNeedPoint(2)
    :onUnitEvent(EVENT.Unit.Hurt, function(hurtData) hurtData.triggerAbility:effective() end)
    :onEvent(EVENT.Ability.Effective,
    function(effectiveData)
        -- 技能被触发的效果
        local tu = effectiveData.triggerUnit
        tu:buff("自在极意被动")
          :duration(3)
          :icon(effectiveData.triggerAbility:icon())
          :purpose(
            function(buffObj)
                buffObj:attach("DivineShieldTarget", "origin", -1)
                buffObj:hurtReduction("+=100"):hurtRebound("+=100"):odds("hurtRebound", "+=100")
            end)
          :rollback(
            function(buffObj)
                buffObj:detach("DivineShieldTarget")
                buffObj:hurtReduction("-=100"):hurtRebound("-=100"):odds("hurtRebound", "-=100")
            end)
          :run()
    end)