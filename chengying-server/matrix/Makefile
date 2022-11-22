
PROTOC ?= protoc
GO ?= go
GOFMT ?= gofmt
GOMOD ?= go mod
GOBINDATA ?= go-bindata
export GO111MODULE=on

VERSION = $(shell git rev-parse --short HEAD)

assets: templates/*
	@echo "Generating templates..."
	$(GOBINDATA) -nocompress -o asset/assets.go -pkg asset templates/

build:
	@echo "go build..."
	go mod tidy
	go build

