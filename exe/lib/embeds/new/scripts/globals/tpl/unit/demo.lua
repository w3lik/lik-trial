--- 单位例子
TPL_UNIT.DEMO = UnitTpl("Footman")
    :speechExtra("avatar")
    :name("步兵")
    :modelAlias("Footman")
    :abilitySlot({ TPL_ABILITY.ZZJY })
    :itemSlot({ TPL_ITEM.DEMO })
    :weaponSound("metal_slice_medium")
    :move(200)
    :attack(10)
    :hp(150)
    :hpRegen(3)
    :reborn(3)

--- 单位敌人
TPL_UNIT.DEMO2 = UnitTpl("Knight")
    :speechExtra("avatar")
    :name("骑士")
    :modelAlias("Knight")
    :itemSlot({ TPL_ITEM.DEMO })
    :weaponSound("metal_slice_medium")
    :move(200)
    :attack(20)
    :hp(2000)

--- 无语音单位
TPL_UNIT.DEMO3 = UnitTpl()

--- 固定不可移动型单位
TPL_UNIT.DEMO4 = UnitTpl(""):speechExtra("immovable"):modelAlias("Rifleman")

--- 建筑型单位
TPL_UNIT.DEMO5 = UnitTpl(""):speechExtra("building"):modelAlias("MortarTeam")