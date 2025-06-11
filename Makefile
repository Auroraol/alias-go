# Makefile for alias-go

# 变量定义
BINARY_NAME=als
MAIN_PACKAGE=.
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# 默认目标
.PHONY: all
all: build

# 构建二进制文件
.PHONY: build
build:
	@echo "构建 ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} ${MAIN_PACKAGE}

# 构建发布版本（多平台）
.PHONY: build-all
build-all: clean
	@echo "构建多平台版本..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ${MAIN_PACKAGE}
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 ${MAIN_PACKAGE}
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 ${MAIN_PACKAGE}
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 ${MAIN_PACKAGE}
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ${MAIN_PACKAGE}

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	go test -v ./...

# 运行测试（带覆盖率）
.PHONY: test-coverage
test-coverage:
	@echo "运行测试（带覆盖率）..."
	go test -v -cover ./...

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	rm -f ${BINARY_NAME}
	rm -rf ${BUILD_DIR}

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 检查代码
.PHONY: vet
vet:
	@echo "检查代码..."
	go vet ./...

# 安装依赖
.PHONY: deps
deps:
	@echo "安装依赖..."
	go mod tidy
	go mod download

# 运行 linter
.PHONY: lint
lint:
	@echo "运行 linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，使用 go vet 代替"; \
		go vet ./...; \
	fi

# 安装到系统
.PHONY: install
install: build
	@echo "安装 ${BINARY_NAME} 到系统..."
	cp ${BINARY_NAME} /usr/local/bin/

# 开发模式：监视文件变化并重新构建
.PHONY: dev
dev:
	@echo "开发模式启动..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air 未安装，请运行: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

# 显示帮助
.PHONY: help
help:
	@echo "可用的命令："
	@echo "  build        - 构建二进制文件"
	@echo "  build-all    - 构建多平台版本"
	@echo "  test         - 运行测试"
	@echo "  test-coverage- 运行测试（带覆盖率）"
	@echo "  clean        - 清理构建文件"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 检查代码"
	@echo "  lint         - 运行 linter"
	@echo "  deps         - 安装依赖"
	@echo "  install      - 安装到系统"
	@echo "  dev          - 开发模式（需要 air）"
	@echo "  help         - 显示此帮助信息"

# 示例用法
.PHONY: example
example: build
	@echo "运行示例..."
	@echo "生成 bash 初始化脚本："
	./${BINARY_NAME} init bash 