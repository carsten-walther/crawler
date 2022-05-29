BINARY = crawler

GOARCH = amd64

BUILD_OS ?= $(shell uname -s)
VERSION ?= $(shell git describe --always)
COMMIT = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

EXECUTABLES = upx
K := $(foreach exec, $(EXECUTABLES), $(if $(shell which $(exec)), some string, $(error "No $(exec) in PATH")))

BUILD_DIR=$(shell pwd)/build
CURRENT_DIR=$(shell pwd)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-s -w -X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

OS_NAME = $(shell echo $(BUILD_OS) | tr '[:upper:]' '[:lower:]')

# Build the project
all: clean linux darwin windows upx install

linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-linux-${GOARCH} . ;

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-darwin-${GOARCH} . ;

windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-windows-${GOARCH}.exe . ;

upx:
	cd ${BUILD_DIR}; \
	upx --brute ${BINARY}* ; \
	cd - >/dev/null

clean:
	rm -f ${BUILD_DIR}/${BINARY}-* ; \
	rm -f /usr/local/bin/${BINARY}

install:
	cp ${BUILD_DIR}/${BINARY}-${OS_NAME}-${GOARCH} /usr/local/bin/${BINARY} ; \

.PHONY: linux darwin windows upx clean install