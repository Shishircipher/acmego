.PHONY: default clean build image
export GO111MODULE=on
export CGO_ENABLED=0

LEGO_IMAGE := acmego/letsacme
MAIN_DIRECTORY := ./main.go

BIN_OUTPUT := $(if $(filter $(shell go env GOOS), windows), dist/lego.exe, dist/lego)

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

default: clean build

clean:
	@echo BIN_OUTPUT: ${BIN_OUTPUT}
	rm -rf dist/ builds/

build: clean
	@echo Version: $(VERSION)
	go build -trimpath -ldflags '-X "main.version=${VERSION}"' -o ${BIN_OUTPUT} ${MAIN_DIRECTORY}

image:
	@echo Version: $(VERSION)
	docker build -t $(LEGO_IMAGE) .
