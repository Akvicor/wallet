PROJECT=wallet

# 系统工具
MKDIR=mkdir
CP=/usr/bin/cp
MV=/usr/bin/mv
RM=/usr/bin/rm
MAKE=/usr/bin/make
SED=/usr/bin/sed
DATE=/usr/bin/date
ARCH=/usr/bin/arch
# 其他工具
GIT=git
GO=go
GOFMT=gofmt

# 编译信息
BRANCH=$(shell $(GIT) rev-parse --abbrev-ref HEAD)
VERSION=$(shell $(GIT) describe --tags --always | $(SED) 's/^v//')
COMMIT=$(shell $(GIT) rev-parse --verify HEAD)
BUILD_TIME=$(shell $(DATE) +"%Y-%m-%d %H:%M:%S %z")
GO_FLAGS=-ldflags "-s -w -X '$(PROJECT)/cmd/def.Branch=$(BRANCH)' -X '$(PROJECT)/cmd/def.Version=$(VERSION)' -X '$(PROJECT)/cmd/def.Commit=$(COMMIT)' -X '$(PROJECT)/cmd/def.BuildTime=$(BUILD_TIME)'"

.PHONY: format build release rewrite example migrate dev count

# 格式化项目
format:
	@$(GO) mod tidy
	@$(GOFMT) -s -w .
	@echo "Version: $(VERSION)-$(BRANCH)_$(COMMIT) ($(BUILD_TIME))"

# 编译项目
build:
	$(MAKE) format
	CGO_ENABLED=0 $(GO) build $(GO_FLAGS) -o bin/$(PROJECT) main.go

# 编译为目标架构发行版
release:
	$(MAKE) format
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) $(GO) build $(GO_FLAGS) -o bin/$(PROJECT) main.go

# 重建数据
rewrite:
	$(MAKE) build
	- $(MKDIR) bin/data
	cd bin && ./$(PROJECT) rewrite -c ./data/config.toml

# 生成默认配置文件
example:
	$(MAKE) build
	- $(MKDIR) bin/data
	cd bin && ./$(PROJECT) example -c > ./data/config.toml

# 创建、初始化数据库
migrate:
	$(MAKE) build
	- $(MKDIR) bin/data
	- $(RM) bin/data/$(PROJECT).db
	cd bin && ./$(PROJECT) migrate -c ./data/config.toml

# 调试运行
dev:
	-$(RM) bin/data/web/*
	-$(CP) -r /kayuki/mvq/workspace/sync/wallet-web/build/* bin/data/web/
	$(MAKE) build
	cd bin && ./$(PROJECT) migrate && ./$(PROJECT) server

