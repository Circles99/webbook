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
	@mockgen -source=./internal/service/user.go -destination=./internal/service/user_mock.go -package=user
	@mockgen -source=./internal/service/code.go -destination=./internal/service/code_mock.go -package=svcmocks
	@go mod tidy