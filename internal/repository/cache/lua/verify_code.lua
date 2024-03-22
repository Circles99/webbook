local key = KEYS[1]

-- 用户输入的 code
local expectedCode = ARGV[1]

local code =  redis.call("get", key)

local cntKey = key..":cnt"
-- 转成数字
local cnt = tonumber(redis.call("get",cntKey))


if cnt <= 0 then
    -- 次数耗尽
    return -1
elseif expectedCode == code then
    redis.call("set", cntKey, -1)
    return 0
else
    -- 可用次数 -1
    redis.call("decr", cntKey)
    -- 验证码错误
    return -2
end


