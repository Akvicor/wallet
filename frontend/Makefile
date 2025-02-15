PROJECT=wallet

# 系统工具
MKDIR=/usr/bin/mkdir
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

.PHONY: build dev count

# 编译项目
build:
	sed -i 's/"version": ".*"/"version": "$(VERSION)"/' package.json
	yarn install
	yarn build
	sed -i 's/"version": ".*"/"version": "0.0.1"/' package.json

# 调试运行
dev:
	yarn start

