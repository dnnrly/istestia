GO111MODULE=on

CURL_BIN ?= curl
GO_BIN ?= go
GORELEASER_BIN ?= goreleaser

PUBLISH_PARAM?=
GO_MOD_PARAM?=-mod vendor
TMP_DIR?=./tmp

BASE_DIR=$(shell pwd)

NAME=istestia

export GO111MODULE=on
export GOPROXY=https://proxy.golang.org
export PATH := ./bin:$(PATH)

.PHONY: install
install: deps

.PHONY: build
build:
	$(GO_BIN) build -v ./cmd/$(NAME)

.PHONY: clean
clean:
	rm -f $(NAME)
	rm -rf dist

.PHONY: clean-deps
clean-deps:
	rm -rf ./bin
	rm -rf ./tmp
	rm -rf ./libexec
	rm -rf ./share

./bin/bats:
	git clone https://github.com/bats-core/bats-core.git ./tmp/bats
	./tmp/bats/install.sh .

./bin/golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.17.1

.PHONY: test-deps
test-deps: ./bin/bats ./bin/golangci-lint
	$(GO_BIN) get -v ./...
	$(GO_BIN) mod tidy

./bin:
	mkdir ./bin

./tmp:
	mkdir ./tmp

./bin/goreleaser: ./bin ./tmp
	$(CURL_BIN) --fail -L -o ./tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/v0.117.2/goreleaser_Linux_x86_64.tar.gz
	gunzip -f ./tmp/goreleaser.tar.gz
	tar -C ./bin -xvf ./tmp/goreleaser.tar

.PHONY: build-deps
build-deps: ./bin/goreleaser

.PHONY: deps
deps: build-deps test-deps

.PHONY: test
test:
	$(GO_BIN) test ./...

.PHONY: acceptance-test
acceptance-test:
	bats -r --tap test/

.PHONY: ci-test
ci-test:
	$(GO_BIN) test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: release
release: clean
	$(GORELEASER_BIN) $(PUBLISH_PARAM)

.PHONY: update
update:
	$(GO_BIN) get -u
	$(GO_BIN) mod tidy
	make test
	make install
	$(GO_BIN) mod tidy
