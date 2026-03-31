#!/usr/bin/make -f

# Project paths
ROOT_PATH := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
BACKEND_PATH := $(ROOT_PATH)/backend
FRONTEND_PATH := $(ROOT_PATH)/frontend
DOCKERFILES_PATH := $(ROOT_PATH)/dockerfiles

# Version and build info
VERSION := $(shell cat $(ROOT_PATH)/VERSION 2>/dev/null || echo "v1.0.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date +%Y%m%d%H%M%S)

# Industry Standard Registry Variables
DOCK_REGISTRY=77kymo
TENCENT_REGISTRY=ccr.ccs.tencentyun.com/itqm-private

# Platform policies
PLATFORMS_DOCKER=linux/amd64,linux/arm64
PLATFORM_CCR=linux/amd64

# Go build environment (Local builds)
GO_PROXY ?= https://goproxy.cn/
GOARCH := $(shell go env GOARCH)
GOOS := linux
CGO_ENABLED ?= 0
GO_BUILD_ENV ?= GOPROXY=${GO_PROXY} GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED}
GO_VERSION := $(shell go version | awk '{print $$3}')

# Build tags based on CodeMode
CodeMode := $(strip $(CodeMode))
ifeq ($(CodeMode),)
    # Defaulting to OpenCode if not specified for buildx sanity, 
    # but still enforcing it for specific targets
    CodeMode := OpenCode
endif

# Build flags module github.com/kymo-mcp/mcpcan
LDFLAGS := -X 'github.com/kymo-mcp/mcpcan/pkg/version.Version=${VERSION}' \
		-X 'github.com/kymo-mcp/mcpcan/pkg/version.BuildTime=${BUILD_TIME}' \
		-X 'github.com/kymo-mcp/mcpcan/pkg/version.Commit=${COMMIT}' \
		-X 'github.com/kymo-mcp/mcpcan/pkg/version.GoVersion=${GO_VERSION}' \
		-X 'github.com/kymo-mcp/mcpcan/pkg/version.CodeMode=${CodeMode}'

# Docker Hub multi-arch push logic
# Paras: 1-DockerfileSuffix, 2-ImageName
define push_image_docker_multiarch
	@echo "Building and pushing multi-arch image $(DOCK_REGISTRY)/$(2):$(VERSION) platforms: $(PLATFORMS_DOCKER)..."
	docker buildx build --platform $(PLATFORMS_DOCKER) \
		--build-arg CodeMode=$(CodeMode) \
		-t $(DOCK_REGISTRY)/$(2):$(VERSION) \
		-t $(DOCK_REGISTRY)/$(2):latest \
		-f $(DOCKERFILES_PATH)/Dockerfile.$(1) \
		--push .
endef

# CCR amd64-only push logic
# Paras: 1-DockerfileSuffix, 2-ImageName
define push_image_ccr_amd64
	@echo "Building and pushing amd64 image $(TENCENT_REGISTRY)/$(2):$(VERSION) platform: $(PLATFORM_CCR)..."
	docker buildx build --platform $(PLATFORM_CCR) \
		--build-arg CodeMode=$(CodeMode) \
		-t $(TENCENT_REGISTRY)/$(2):$(VERSION) \
		-t $(TENCENT_REGISTRY)/$(2):latest \
		-f $(DOCKERFILES_PATH)/Dockerfile.$(1) \
		--push .
endef

# README Push logic (Industry Standard)
# Paras: 1-ImageName, 2-ShortDesc
define push_readme_doc
	@if [ -n "$(DOCKER_USER)" ] && [ -n "$(DOCKER_PASS)" ]; then \
		echo "Updating README for $(1)..."; \
		docker run --rm -v $(ROOT_PATH):/data \
			-e DOCKER_USER=$(DOCKER_USER) \
			-e DOCKER_PASS=$(DOCKER_PASS) \
			chko/docker-pushrm:1 \
			--file /data/README.md \
			--short "$(2)" \
			$(DOCK_REGISTRY)/$(1); \
	fi
endef

# Backend build targets (Local)
define go_build_service
	@echo "---------- Start Go build $(1) ----------"
	@cd $(BACKEND_PATH) && $(GO_BUILD_ENV) go build -tags OpenCode -ldflags "$(LDFLAGS)" -o $(BACKEND_PATH)/bin/$(1) $(BACKEND_PATH)/cmd/$(1)/main.go
	@echo "---------- End Go build $(1) ----------"
endef

.PHONY: pnpm-build
pnpm-build:
	@echo "---------- Start build frontend ----------"
	@cd $(FRONTEND_PATH) && pnpm i && pnpm build
	@echo "---------- End build frontend ----------"

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  push-all           - Push CCR amd64 + Docker multi-arch images (Market, Authz, Web)"
	@echo "  push-market        - Push Market image (CCR amd64, Docker multi-arch)"
	@echo "  push-authz         - Push Authz image (CCR amd64, Docker multi-arch)"
	@echo "  push-frontend      - Push Frontend image (CCR amd64, Docker multi-arch)"
	@echo "  proto-buf          - Generate protobuf files"
	@echo "  go-build-market    - Local Go build for market"
	@echo "  sync-docs          - Sync documentation from tools submodule"
	@echo "  clean              - Remove build artifacts"

.PHONY: push-all
push-all: push-market push-authz push-frontend

.PHONY: push-market
push-market:
	$(call push_image_ccr_amd64,market,mcp-market)
	$(call push_image_docker_multiarch,market,mcp-market)
	$(call push_readme_doc,mcp-market,MCP Market Service)

.PHONY: push-authz
push-authz:
	$(call push_image_ccr_amd64,authz,mcp-authz)
	$(call push_image_docker_multiarch,authz,mcp-authz)
	$(call push_readme_doc,mcp-authz,MCP Authorization Service)

.PHONY: push-frontend
push-frontend:
	$(call push_image_ccr_amd64,frontend,mcp-web)
	$(call push_image_docker_multiarch,frontend,mcp-web)
	$(call push_readme_doc,mcp-web,MCP Web Frontend)

# Remaining Utility targets

.PHONY: proto-buf
proto-buf:
	@echo "---- Cleaning existing generated files ----"
	@rm -rf $(shell find $(BACKEND_PATH)/api -type f -name '*.go')
	@rm -rf $(shell find $(BACKEND_PATH)/api -type f -name '*.json')
	@echo "---- Generating protobuf files ----"
	@cd $(BACKEND_PATH)/api && buf --debug generate 
	@find $(BACKEND_PATH)/api -name "*.pb.go" -exec protoc-go-inject-tag -input={} \; || echo "No .pb.go files found for tag injection"
	@echo "---- Merging swagger files ----"
	@rm -f $(BACKEND_PATH)/api/merged.swagger.json
	@if [ -n "$$(find $(BACKEND_PATH)/api -name '*.json' -type f)" ]; then \
		swagger mixin $$(find $(BACKEND_PATH)/api -name '*.json' -type f) -o $(BACKEND_PATH)/api/merged.swagger.json; \
		echo "Swagger files merged successfully"; \
	else \
		echo "No swagger JSON files found to merge"; \
		touch $(BACKEND_PATH)/api/merged.swagger.json; \
	fi

.PHONY: go-mod-tidy
go-mod-tidy:
	@cd $(BACKEND_PATH) && go mod tidy

.PHONY: go-build-market
go-build-market:
	$(call go_build_service,market)

.PHONY: go-build-authz
go-build-authz:
	$(call go_build_service,authz)

.PHONY: clean
clean:
	@rm -rf $(BACKEND_PATH)/bin/*
	@rm -rf $(FRONTEND_PATH)/dist

.PHONY: sync-docs
sync-docs:
	@echo "---------- Syncing Documentation ----------"
	@mkdir -p $(ROOT_PATH)/docs/zh $(ROOT_PATH)/docs/en
	@if [ -d "../mcpcan-tools/mcpcan-platform-docs/zh/guide" ]; then \
		cp -r ../mcpcan-tools/mcpcan-platform-docs/zh/guide $(ROOT_PATH)/docs/zh/; \
		echo "Synced zh/guide to docs/zh/"; \
	fi
	@if [ -d "../mcpcan-tools/mcpcan-platform-docs/en/guide" ]; then \
		cp -r ../mcpcan-tools/mcpcan-platform-docs/en/guide $(ROOT_PATH)/docs/en/; \
		echo "Synced en/guide to docs/en/"; \
	fi
	@echo "---------- Sync Completed ----------"

# @Deprecated / Removed targets marked as error to guide migration
.PHONY: docker-build-init docker-build-gateway docker-build-sidecar
docker-build-init docker-build-gateway docker-build-sidecar:
	@echo "Error: This target is deprecated and has been migrated to mcpcan-tools or removed."
	@exit 1
