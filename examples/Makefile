# Makefile for examples directory
# 本 Makefile 用于管理 examples 目录下的多个示例程序

# 默认目标
.DEFAULT_GOAL := help

# Go 命令
GO := go

# 格式化工具
GOFMT := gofmt

# 示例列表
EXAMPLES := \
	charmtone \
	cellbuf \
	faketty \
	img2term \
	layout \
	mosaic \
	parserlog \
	parserlog2 \
	pen \
	toner

# 帮助信息
.PHONY: help
help:
	@echo "examples 目录的 Makefile 命令："
	@echo "  make init      - 初始化示例的模块"
	@echo "  make reinit    - 重新初始化示例的模块（备份并重置）"
	@echo "  make tidy      - 整理示例的依赖"
	@echo "  make format    - 格式化示例的代码"
	@echo "  make build     - 构建所有示例"
	@echo "  make run-<示例> - 运行指定的示例（例如：make run-toner）"
	@echo "  make clean     - 清理示例的构建产物"
	@echo "  make help      - 显示此帮助信息"
	@echo ""
	@echo "可用的示例："
	@for example in $(EXAMPLES); do \
		echo "  $$example"; \
	 done

# 初始化模块
.PHONY: init
init:
	@echo "初始化 examples 目录的模块..."
	@if [ -f "go.mod" ]; then \
		echo "go.mod 文件已存在，跳过初始化..."; \
		echo "运行 make tidy 整理依赖..."; \
		$(GO) mod tidy; \
	else \
		$(GO) mod init github.com/purpose168/charm-experimental-packages-cn/examples; \
		echo "模块初始化完成，运行 make tidy 添加依赖..."; \
		$(GO) mod tidy; \
	fi

# 重新初始化模块（备份并重置）
.PHONY: reinit
reinit:
	@echo "重新初始化 examples 目录的模块..."
	@if [ -f "go.mod" ]; then \
		echo "删除 go.mod 和 go.sum 文件..."; \
		rm -f go.mod go.sum; \
		echo "重新初始化模块..."; \
		$(GO) mod init github.com/purpose168/charm-experimental-packages-cn/examples; \
		echo "模块初始化完成，运行 make tidy 添加依赖..."; \
		$(GO) mod tidy; \
		echo "模块重新初始化完成！"; \
	else \
		echo "go.mod 文件不存在，运行 make init 初始化模块..."; \
		$(MAKE) init; \
	fi

# 整理依赖
.PHONY: tidy
tidy:
	@echo "整理 examples 目录的依赖..."
	@$(GO) mod tidy

# 格式化代码
.PHONY: format
format:
	@echo "格式化 examples 目录的代码..."
	@$(GOFMT) -s -w .

# 构建所有示例
.PHONY: build
build:
	@echo "构建所有示例..."
	@for example in $(EXAMPLES); do \
		echo "构建 $$example 示例..."; \
		cd $(CURDIR)/$$example && $(GO) build -o $$example .; \
	 done

# 清理构建产物
.PHONY: clean
clean:
	@echo "清理 examples 目录的构建产物..."
	@for example in $(EXAMPLES); do \
		echo "清理 $$example 示例的构建产物..."; \
		rm -f $(CURDIR)/$$example/$$example; \
	 done

# 为每个示例创建运行目标
.PHONY: run-charmtone
run-charmtone:
	@echo "运行 charmtone 示例..."
	@cd $(CURDIR)/charmtone && $(GO) run .

.PHONY: run-cellbuf
run-cellbuf:
	@echo "运行 cellbuf 示例..."
	@cd $(CURDIR)/cellbuf && $(GO) run .

.PHONY: run-faketty
run-faketty:
	@echo "运行 faketty 示例..."
	@cd $(CURDIR)/faketty && $(GO) run .

.PHONY: run-img2term
run-img2term:
	@echo "运行 img2term 示例..."
	@cd $(CURDIR)/img2term && $(GO) run .

.PHONY: run-layout
run-layout:
	@echo "运行 layout 示例..."
	@cd $(CURDIR)/layout && $(GO) run .

.PHONY: run-mosaic
run-mosaic:
	@echo "运行 mosaic 示例..."
	@cd $(CURDIR)/mosaic && $(GO) run .

.PHONY: run-parserlog
run-parserlog:
	@echo "运行 parserlog 示例..."
	@cd $(CURDIR)/parserlog && $(GO) run .

.PHONY: run-parserlog2
run-parserlog2:
	@echo "运行 parserlog2 示例..."
	@cd $(CURDIR)/parserlog2 && $(GO) run .

.PHONY: run-pen
run-pen:
	@echo "运行 pen 示例..."
	@cd $(CURDIR)/pen && $(GO) run .

.PHONY: run-toner
run-toner:
	@echo "运行 toner 示例..."
	@cd $(CURDIR)/toner && $(GO) run .
