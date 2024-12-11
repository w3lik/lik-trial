local process = Process("start")

process:onStart(function(self)

    -- 调试自动去除迷雾
    Game():fog(not DEBUGGING):mark(not DEBUGGING)

    self:next("test")

end)
