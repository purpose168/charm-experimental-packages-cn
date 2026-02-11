# Makefile for exp/toner package
# 本 Makefile 用于管理 exp/toner 包的构建、测试和依赖管理

# 默认目标
.DEFAULT_GOAL := help

# Go 命令
GO := go

# 代码检查工具
GOLANGCI_LINT := golangci-lint

# 格式化工具
GOFMT := gofmt

# 包路径
PACKAGE := github.com/purpose168/charm-experimental-packages-cn/exp/toner

# 帮助信息
.PHONY: help
help:
	@echo "exp/toner 包的 Makefile 命令："
	@echo "  make init      - 初始化包的模块"
	@echo "  make test      - 运行包的测试"
	@echo "  make tidy      - 整理包的依赖"
	@echo "  make format    - 格式化包的代码"
	@echo "  make lint      - 检查包的代码质量"
	@echo "  make clean     - 清理包的构建产物"
	@echo "  make help      - 显示此帮助信息"

# 初始化模块
.PHONY: init
init:
	@echo "初始化 exp/toner 包的模块..."
	@if [ -f "go.mod" ]; then \
		echo "go.mod 文件已存在，跳过初始化..."; \
		echo "运行 make tidy 整理依赖..."; \
		$(GO) mod tidy; \
	else \
		$(GO) mod init $(PACKAGE); \
		echo "模块初始化完成，运行 make tidy 添加依赖..."; \
		$(GO) mod tidy; \
	fi

# 运行测试
.PHONY: test
test:
	@echo "运行 exp/toner 包的测试..."
	@$(GO) test -v ./...

# 整理依赖
.PHONY: tidy
tidy:
	@echo "整理 exp/toner 包的依赖..."
	@$(GO) mod tidy

# 格式化代码
.PHONY: format
format:
	@echo "格式化 exp/toner 包的代码..."
	@$(GOFMT) -s -w .

# 检查代码质量
.PHONY: lint
lint:
	@echo "检查 exp/toner 包的代码质量..."
	@if command -v $(GOLANGCI_LINT) > /dev/null; then \
		$(GOLANGCI_LINT) run; \
	else \
		echo "警告: golangci-lint 未安装，跳过代码检查"; \
		$(GO) vet ./...; \
	fi

# 清理构建产物
.PHONY: clean
clean:
	@echo "清理 exp/toner 包的构建产物..."
	@$(GO) clean ./...

# 验证依赖
.PHONY: verify
verify:
	@echo "验证 exp/toner 包的依赖..."
	@$(GO) mod verify

# 显示依赖树
.PHONY: deptree
deptree:
	@echo "显示 exp/toner 包的依赖树..."
	@$(GO) list -m -u -json all | gojq -r '(.Path + " " + .Version)'
