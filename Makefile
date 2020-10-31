.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --no-builtin-rules
VERSION=`cat version`
BUILD=`date +%FT%T%z`
#COMMIT=`git rev-parse HEAD`
COMMIT=`date +%FT%T%z`
EXECUTABLE="dunzo"

all: build

.PHONY: build test clean generate dist init build_linux build_mac functional_test unit_test

build: #generate
	@go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go

test: unit_test functional_test

unit_test:
	go list ./... | grep -v tests | xargs go test -count 1 -cover -race -timeout 1m -tags=unit_test

functional_test:
	go list ./tests/... | xargs go test -count 1 -timeout 1m

run: build
	@./${EXECUTABLE}

generate:
	@echo " > generating resources"
	@go generate ./resources

clean:
	@rm -rf ${EXECUTABLE} dist/

build_nix:
	@env GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go

build_mac:
	@env GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go