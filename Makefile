SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = comfo
ALPINE_BINARY = comfo_alpine

VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date +%FT%T%z)
GITREV := $(shell git rev-parse HEAD)
GOVERSION := $(shell go version)

LDFLAGS := -ldflags '-X "main.GitRev=${GITREV}" -X "main.Version=${VERSION}" -X "main.BuildTime=${BUILD_TIME}" -X "main.GoVersion=${GOVERSION}"'

PACKAGES := $(shell go list ./... | grep -v '/vendor/')

# Require the Go compiler/toolchain to be installed
ifeq (, $(shell which go 2>/dev/null))
$(error No 'go' found in $(PATH), please install the Go compiler for your system)
endif

.DEFAULT_GOAL: $(BINARY)

# This target needs to be named after the file it generates
$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY}

.PHONY: clean
clean:
	@if [ -f ${BINARY} ] ; then rm -f ${BINARY} ; echo "Removed file '${BINARY}'." ; fi

.PHONY: clean_release
clean_release:
	@echo "Cleaning all *.tar.gz in repository.."
	rm -f *.tar.gz
	echo "Cleaning all kermit_*_amd64 in repository.."
	rm -f ${BINARY}_*_amd64

.PHONY: test
test:
	go test $(PACKAGES)

.ONESHELL:
.PHONY: cover
cover:
	go test -coverprofile=coverage.out -covermode=count
	go tool cover -func=coverage.out

.ONESHELL:
.PHONY: gox
gox:
	@if [ ! `command -v gox` ]; then
		echo "Installing gox.."
		go get -u github.com/mitchellh/gox
		echo "Successfully installed gox!"
	fi

.PHONY: release
release: gox clean_release
	@gox -osarch="darwin/amd64 linux/amd64" ${LDFLAGS} --output "${BINARY}_{{.OS}}_{{.Arch}}"

	echo "Archiving:" *_amd64
	tar -czvf ${BINARY}-linux-${VERSION}.tar.gz ${BINARY}_linux_amd64
	tar -czvf ${BINARY}-macos-${VERSION}.tar.gz ${BINARY}_darwin_amd64
	echo "Built artifacts:" *.tar.gz
