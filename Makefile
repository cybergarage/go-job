# Copyright (C) 2025 The go-job Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

GOBIN := $(shell go env GOPATH)/bin
PATH := $(GOBIN):$(PATH)

MODULE_ROOT=github.com/cybergarage/go-job
PKG_NAME=job
PKG_COVER=${PKG_NAME}-cover

PKG_ID=${MODULE_ROOT}/${PKG_NAME}
PKG_SRC_DIR=${PKG_NAME}
PKG=${MODULE_ROOT}/${PKG_SRC_DIR}

TEST_PKG_NAME=${PKG_NAME}test
TEST_PKG_ID=${MODULE_ROOT}/${TEST_PKG_NAME}
TEST_PKG_DIR=${TEST_PKG_NAME}
TEST_PKG=${MODULE_ROOT}/${TEST_PKG_DIR}

.PHONY: format vet lint clean
.IGNORE: lint

all: test

version:
	@pushd ${PKG_SRC_DIR} && ./version.gen > version.go && popd
	-git commit ${PKG_SRC_DIR}/version.go -m "Update version"

format: version
	gofmt -s -w ${PKG_SRC_DIR} ${TEST_PKG_DIR}

vet: format
	go vet ${PKG_ID} ${TEST_PKG_ID}

lint: vet
	golangci-lint run ${PKG_SRC_DIR}/... ${TEST_PKG_DIR}/...

godoc:
	go install golang.org/x/tools/cmd/godoc@latest
	godoc -http=:6060 -play

test: lint
	go test -v -p 1 -timeout 10m -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

clean:
	go clean -i ${PKG}

# Protobuf generation

PKG_PROTO_ROOT=${PKG_SRC_DIR}/api
protopkg:
	go get -u google.golang.org/protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest	
%.pb.go : %.proto protopkg
	protoc -I=${PKG_PROTO_ROOT}/proto/v1 --go_out=paths=source_relative:${PKG_PROTO_ROOT}/gen/go/v1 --go-grpc_out=paths=source_relative:${PKG_PROTO_ROOT}/gen/go/v1 --plugin=protoc-gen-go=${GOBIN}/protoc-gen-go --plugin=protoc-gen-go-grpc=${GOBIN}/protoc-gen-go-grpc $<
protos=$(shell find ${PKG_PROTO_ROOT} -name '*.proto')
pbs=$(protos:.proto=.pb.go)
proto: protopkg $(pbs)

# Documentation generation

%.md : %.adoc
	asciidoctor -b docbook -a leveloffset=+1 -o - $< | pandoc -t markdown_strict --wrap=none -f docbook > $@

%.png : %.pu
	plantuml -tpng $<

images := $(wildcard doc/img/*.png)
docs := $(wildcard doc/*.md)
doc: $(docs) $(images)
	@echo "Generated: $(docs)"
	@echo "Generated: $(images)"
