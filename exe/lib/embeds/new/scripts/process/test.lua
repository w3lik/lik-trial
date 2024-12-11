local process = Process("test")

process:onStart(function(this)

    print("Hello, Lik!")
    echo("你好，Lik！")

    Bgm():play("lik")

    local bubble = this:bubble()
    bubble.eff = Effect("Echo", 0, 0, 0, -1):rotateZ(270)
    bubble.unit = Unit(TPL_UNIT.DEMO, Player(1), 0, -100, 270)
    bubble.unit2 = Unit(TPL_UNIT.DEMO2, Player(2), 0, -100, 270)

end)

process:onOver(function(this)

    Bgm():stop()

end)
