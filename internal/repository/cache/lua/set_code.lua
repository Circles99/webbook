-- 你的验证码在 redis上的key
-- phone:code:login:1234xxx
local key = KEYS[1]
-- 验证次数，一个验证码 最多3次，记录了验证了几次
-- phone:code:login:1234xxx:cnt
local cntKey = key..":cnt"
-- 你的验证码
local val = ARGV[1]

-- 获取过期时间
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- key 存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl < 540  then
    -- key 不存在，或者过期时间小于9分钟
    -- 设置key 和 过期时间
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set",cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送频繁
    return -1
end