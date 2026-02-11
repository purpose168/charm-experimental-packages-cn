# Makefile for charm-experimental-packages-cn
# 本 Makefile 用于管理多个实验性 Go 包的构建、测试、打包和安装

# 默认目标
.DEFAULT_GOAL := help

# 项目根目录
ROOT_DIR := $(shell pwd)

# Go 命令
GO := go

# 代码检查工具
GOLANGCI_LINT := golangci-lint

# 格式化工具
GOFMT := gofmt

# 现代化工具
MODERNIZE := modernize

# 包列表
PACKAGES := \
	ansi \
	cellbuf \
	colors \
	conpty \
	editor \
	errors \
	examples \
	exp/golden \
	exp/higherorder \
	exp/maps \
	exp/open \
	exp/ordered \
	exp/slice \
	exp/strings \
	exp/teatest \
	exp/teatest/v2 \
  exp/charmtone \
  exp/toner \
	input \
	json \
	sshkey \
	term \
	termios \
	vt \
	wcwidth \
	windows \
	xpty

# 检测操作系统
OS := $(shell uname -s)

# 操作系统特定的包
ifeq ($(OS),Windows_NT)
	WINDOWS_PACKAGES := conpty windows
else ifeq ($(OS),Darwin)
	# macOS 特定的包
	MACOS_PACKAGES :=
else ifeq ($(OS),Linux)
	# Linux 特定的包
	LINUX_PACKAGES :=
endif

# 帮助信息
.PHONY: help
help:
	@echo "charm-experimental-packages-cn 构建工具"
	@echo ""
	@echo "使用方法: make [目标] [包名]"
	@echo ""
	@echo "主要目标:"
	@echo "  all          - 构建所有包"
	@echo "  build        - 构建指定的包 (默认构建所有包)"
	@echo "  test         - 运行指定包的测试 (默认运行所有包的测试)"
	@echo "  test-all     - 运行所有包的测试，包括集成测试"
	@echo "  lint         - 运行代码检查"
	@echo "  tidy         - 清理所有包的依赖"
	@echo "  format       - 格式化所有包的代码"
	@echo "  clean        - 清理构建产物"
	@echo "  install      - 安装指定的包"
	@echo "  help         - 显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build ansi cellbuf  # 构建 ansi 和 cellbuf 包"
	@echo "  make test input          # 运行 input 包的测试"
	@echo "  make install term        # 安装 term 包"

# 构建所有包
.PHONY: all
all:
	@echo "\n\n构建所有包..."
	@for pkg in $(PACKAGES); do \
		$(MAKE) build $$pkg; \
	done

# 构建指定的包
.PHONY: build
build:
	@if [ "$(filter-out build,$(MAKECMDGOALS))" = "" ]; then \
		$(MAKE) all; \
	else \
		for pkg in $(filter-out build,$(MAKECMDGOALS)); do \
			echo "\n\n构建 $$pkg 包..."; \
			cd $(ROOT_DIR)/$$pkg && $(GO) build ./...; \
			done; \
	fi

# 运行所有包的测试
.PHONY: test
test:
	@if [ "$(filter-out test,$(MAKECMDGOALS))" = "" ]; then \
		echo "运行所有包的测试..."; \
		for pkg in $(PACKAGES); do \
			echo "测试 $$pkg 包..."; \
			cd $(ROOT_DIR)/$$pkg && $(GO) test ./...; \
			done; \
	else \
		for pkg in $(filter-out test,$(MAKECMDGOALS)); do \
			echo "测试 $$pkg 包..."; \
			cd $(ROOT_DIR)/$$pkg && $(GO) test ./...; \
			done; \
	fi

# 运行所有包的测试，包括集成测试
.PHONY: test-all
test-all:
	@echo "运行所有包的测试，包括集成测试..."
	@for pkg in $(PACKAGES); do \
		echo "测试 $$pkg 包..."; \
		cd $(ROOT_DIR)/$$pkg && $(GO) test -v ./...; \
	 done

# 运行代码检查
.PHONY: lint
lint:
	@echo "运行代码检查..."
	@for pkg in $(PACKAGES); do \
		echo "检查 $$pkg 包..."; \
		cd $(ROOT_DIR)/$$pkg && $(GOLANGCI_LINT) run; \
	 done

# 整理所有包的依赖（添加缺失的依赖并删除未使用的依赖）
.PHONY: tidy
tidy:
	@echo "整理所有包的依赖..."
	@for pkg in $(PACKAGES); do \
		echo "整理 $$pkg 包的依赖..."; \
		cd $(ROOT_DIR)/$$pkg && $(GO) mod tidy; \
	 done

# 重新初始化所有包的模块
.PHONY: reinit
reinit:
	@echo "重新初始化所有包的模块..."
	@for pkg in $(PACKAGES); do \
		echo "\n\n重新初始化 $$pkg 包的模块..."; \
		cd $(ROOT_DIR)/$$pkg && \
		if [ -f "go.mod" ]; then \
			echo "删除 go.mod 和 go.sum 文件..."; \
			rm -f go.mod go.sum; \
			echo "重新初始化模块..."; \
			module_path="github.com/purpose168/charm-experimental-packages-cn/$$pkg"; \
			$(GO) mod init $$module_path; \
			echo "运行 go mod tidy 重新添加依赖..."; \
			$(GO) mod tidy; \
			echo "$$pkg 包的模块重新初始化完成..."; \
		else \
			echo "$$pkg 包没有 go.mod 文件，跳过..."; \
		fi; \
	 done

# 格式化所有包的代码
.PHONY: format
format:
	@echo "格式化所有包的代码..."
	@for pkg in $(PACKAGES); do \
		echo "格式化 $$pkg 包的代码..."; \
		cd $(ROOT_DIR)/$$pkg && $(GOFMT) -s -w .; \
	 done

# 清理构建产物
.PHONY: clean
clean:
	@echo "清理构建产物..."
	@for pkg in $(PACKAGES); do \
		echo "清理 $$pkg 包的构建产物..."; \
		cd $(ROOT_DIR)/$$pkg && $(GO) clean ./...; \
	 done

# 安装指定的包
.PHONY: install
install:
	@if [ "$(filter-out install,$(MAKECMDGOALS))" = "" ]; then \
		echo "请指定要安装的包名"; \
		echo "示例: make install ansi"; \
	else \
		for pkg in $(filter-out install,$(MAKECMDGOALS)); do \
			echo "安装 $$pkg 包..."; \
			cd $(ROOT_DIR)/$$pkg && $(GO) install ./...; \
			done; \
	fi

# 处理命令行参数
%:
	@:
