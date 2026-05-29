.PHONY: build run test clean migrate-up migrate-down docker-build lint fmt

# Go 参数
BINARY_NAME=server
CMD_PATH=cmd/server/main.go
BUILD_DIR=bin

# 构建
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行（开发模式）
run:
	go run $(CMD_PATH)

# 热重载（需安装 air）
dev:
	@which air > /dev/null 2>&1 || go install github.com/cosmtrek/air@latest
	air

# 测试
test:
	go test -v -race ./...

# 测试覆盖率
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 清理
 clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# 代码格式化
fmt:
	gofmt -w -s .
	goimports -w .

# 代码检查
lint:
	@which golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

# 数据库迁移
MIGRATE_URL=postgres://lagom:lagom@localhost:5432/lagom?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(MIGRATE_URL)" up

migrate-down:
	migrate -path migrations -database "$(MIGRATE_URL)" down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Docker
docker-build:
	docker build -t lagom-server:latest .

docker-run:
	docker-compose up --build

# 依赖管理
deps:
	go mod tidy
	go mod download

deps-update:
	go get -u ./...
	go mod tidy

# 代码生成（SQLC）
generate:
	sqlc generate

# 帮助
help:
	@echo "Lagom Server - Make targets:"
	@echo "  build          编译二进制文件"
	@echo "  run            直接运行（开发）"
	@echo "  dev            热重载开发（air）"
	@echo "  test           运行测试"
	@echo "  test-coverage  生成测试覆盖率报告"
	@echo "  fmt            格式化代码"
	@echo "  lint           代码静态检查"
	@echo "  migrate-up     执行数据库迁移"
	@echo "  migrate-down   回滚数据库迁移"
	@echo "  migrate-create 创建新迁移文件"
	@echo "  docker-build   构建 Docker 镜像"
	@echo "  docker-run     Docker Compose 启动"
	@echo "  deps           整理依赖"
	@echo "  generate       生成代码（sqlc）"
	@echo "  clean          清理构建产物"

.DEFAULT_GOAL := help
