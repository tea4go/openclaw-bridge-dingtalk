.PHONY: build run clean test

# 构建二进制文件
build:
	go build -o bin/dingtalk-bridge main.go handler.go session.go

# 运行开发服务器
run:
	go run main.go handler.go session.go

# 清理构建产物
clean:
	rm -rf bin/

# 运行测试
test:
	go test -v ./...

# 下载依赖
deps:
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 检查代码
lint:
	golangci-lint run

# Docker 构建
docker-build:
	docker build -t dingtalk-bridge:latest .

# Docker 运行
docker-run:
	docker run -p 3000:3000 -v $(PWD)/config.yaml:/app/config.yaml dingtalk-bridge:latest