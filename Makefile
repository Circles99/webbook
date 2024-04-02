.PHONY: docker

docker:
	# 把上次编译的删除掉
	@rm webook || true
	# 防止go mod文件不对 便是失败
	@go mod tidy
 	# 打包项目，指定平台和 架构
	@GOOS=linux GOARCH=amd64 go build -o webook .
	# 删除之前的镜像
	@docker rmi -f circles99/webook:v0.0.1
	# 打包到docker镜像
	@docker build -t circles99/webook:v0.0.1 .


.PHONY: mock
mock:
    # mock service.go
	@mockgen -source=./internal/service/user.go -destination=./internal/service/mocks/user_mock.go -package=svcmocks
	@mockgen -source=./internal/service/code.go -destination=./internal/service/mocks/code_mock.go -package=svcmocks
    # mock repository
	@mockgen -source=./internal/repository/code.go -destination=./internal/repository/mocks/code_mock.go -package=repmocks
	@mockgen -source=./internal/repository/user.go -destination=./internal/repository/mocks/user_mock.go -package=repmocks
	# mock dao
	@mockgen -source=./internal/repository/dao/user.go -destination=./internal/repository/dao/mocks/user_mock.go -package=daomocks
	@mockgen -source=./internal/repository/cache/user.go -destination=./internal/repository/cache/mocks/user_mock.go -package=cachemocks
	# mock redis
	@mockgen -package=redismocks -destination=./internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy