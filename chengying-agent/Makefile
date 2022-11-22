
PROTOC ?= protoc
GO ?= go
GOFMT ?= gofmt
GODEP ?= godep
GOBINDATA ?= go-bindata

VERSION = $(shell git rev-parse --short HEAD)

PROD_PKG_PREFIX = easyagent/prod

OUT_DIR = build
PROTO_DIR = internal/proto

GOBUILD_OPTS = -ldflags "-s -X easyagent/internal.VERSION=$(VERSION)"

define BUILD_RELEASE_TARGET
	$(eval os := $1)
	$(eval arch := $2)
	$(eval target := $3)
	$(eval mode := $4)
	@echo "Build $(mode)-version for $(target) $(os)-$(arch) ..."
	$(eval output := $(if $(filter $(target),easy-agent-server),server,$(target)))
	$(eval target_dir := $(OUT_DIR)/$(mode)/$(output)$(if $(filter $(mode),release),/$1-$2,))
	@mkdir -p $(target_dir)
	$(eval bin_name := $(target)$(if $(filter $(os),windows),.exe,))
	$(if $(filter $(mode),release),
		GOOS=$(os) GOARCH=$(arch) $(GO) build $(GOBUILD_OPTS) -o $(target_dir)/$(bin_name) $(PROD_PKG_PREFIX)/$(target),
		$(GO) build $(GOBUILD_OPTS) -o $(target_dir)/$(bin_name) $(PROD_PKG_PREFIX)/$(target))
endef

asset: templates/*
	@echo "Generating templates..."
	$(GOBINDATA) -nocompress -o internal/server/asset/assets.go -pkg asset templates/


proto: $(PROTO_DIR)/*.proto
	@echo "Generating protocol codes ..."
	$(PROTOC) -I=$(PROTO_DIR) -I=$(GOPATH)/src/ -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf/ \
		--gogo_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,plugins=grpc:$(PROTO_DIR) $?

server-release:
	$(call BUILD_RELEASE_TARGET,linux,amd64,easy-agent-server,release)

sidecar-release: ## build sidecar client
	$(call BUILD_RELEASE_TARGET,linux,amd64,sidecar,release)
	$(call BUILD_RELEASE_TARGET,windows,amd64,sidecar,release)

clean-release:
	@rm -rf $(OUT_DIR)/release

release: clean-release proto server-release sidecar-release asset

server-debug:
	$(call BUILD_RELEASE_TARGET,,,easy-agent-server,debug)

sidecar-debug:
	$(call BUILD_RELEASE_TARGET,,,sidecar,debug)

clean-debug:
	@rm -rf $(OUT_DIR)/debug

debug: clean-debug proto sidecar-debug server-debug asset

clean: clean-release clean-debug

all: debug release

.DEFAULT_GOAL := debug

