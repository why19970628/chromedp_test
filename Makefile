.DEFAULT: help

IMAGE_NAME ?= lampnick/doctron
CENTOS_IMAGE_TAG ?= v0.3.3-centos
ALPINE_IMAGE_TAG ?= v0.3.3-alpine

help: Makefile
	@echo "Doctron is a document convert tools for html pdf image etc.\r\n"
	@echo "Usage: make <command>\r\n\r\nThe commands are:\r\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
# 	@sed -n 's/^##.*:/\033[34m &: \033[0m/p' $< | column -t -s ':' |sed -e 's/^/  /' | awk '{print $0}'

## build-runtime-alpine: build a runtime docker image with alpine.
build-alpine:
	@docker build -f Dockerfile_alpine -t chromium-alpine:latest .

build-debian:
	@docker build -f Dockerfile_debian -t chromium-debian:latest .

build-biz:
	@docker build -f Dockerfile -t cchormedp_test .
