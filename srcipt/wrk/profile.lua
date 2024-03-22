method="Post"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "PostmanRuntime/7.32.3"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
--wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4zMi4zIiwiZXhwIjoxNjkwMjczNjUwfQ.qmZ2jwT-JxDy4uGpuKJLSudEDpoxC1FDOe_KciNZbO8"


wrk.headers["User-Agent"] = ""
request = function()
    body = '{"email":"%s%d@qq.com", "password":"hello#world123", "confirmPassword": "hello#world123"}'
    return wrk.format(method, wrk.path, wrk.headers, body)
end


response = function(status, headers, body)
    if not token and status == 200 then
        token = headers["X-Jwt-Token"]
        path="users/profile"
        method = "GET"
        wrk.headers["Authorization"] = string.format("Bearer %s", token)
    end
end