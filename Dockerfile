# 基础镜像
FROM ubuntu:20.04

# 把编译后的打包进来这个镜像
COPY webook /app/webook

# 设定工作目录
WORKDIR /app


# 执行命令
ENTRYPOINT ["/app/webook"]


