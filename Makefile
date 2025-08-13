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


BIN_SRC_DIR=cmd
BIN_ID=${MODULE_ROOT}/${BIN_SRC_DIR}
BIN_CLI=${PKG_NAME}ctl
BIN_CLI_ID=${BIN_ID}/${BIN_CLI}
BIN_SERVER=${PKG_NAME}d
BIN_SERVER_ID=${BIN_ID}/${BIN_SERVER}
BIN_SRCS=\
        ${BIN_SRC_DIR}/${BIN_CLI} \
        ${BIN_SRC_DIR}/${BIN_SERVER}
BINS=\
        ${BIN_CLI_ID} \
        ${BIN_SERVER_ID}

.PHONY: format vet lint clean
.IGNORE: lint

all: test

version:
	@pushd ${PKG_SRC_DIR} && ./version.gen > version.go && popd
	-git commit ${PKG_SRC_DIR}/version.go -m "Update version"

format: version
	gofmt -s -w ${PKG_SRC_DIR} ${TEST_PKG_DIR} ${BIN_SRC_DIR}

vet: format
	go vet ${PKG_ID} ${TEST_PKG_ID}

lint: vet
	golangci-lint run ${PKG_SRC_DIR}/... ${TEST_PKG_DIR}/...

godoc:
	go install golang.org/x/tools/cmd/godoc@latest
	open http://localhost:6060/pkg/${PKG_ID}/ || xdg-open http://localhost:6060/pkg/${PKG_ID}/ || gnome-open http://localhost:6060/pkg/${PKG_ID}/
	godoc -http=:6060 -play

test: lint
	go test -v -p 1 -timeout 10m -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

cover: test
	open ${PKG_COVER}.html || xdg-open ${PKG_COVER}.html || gnome-open ${PKG_COVER}.html

build:
	go build -v -gcflags=${GCFLAGS} -ldflags=${LDFLAGS} ${BINS}

install:
	go install -v -gcflags=${GCFLAGS} -ldflags=${LDFLAGS} ${BINS}

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
	git commit ${PKG_PROTO_ROOT} -m "feat: update $(notdir $<)"
protos=$(shell find ${PKG_PROTO_ROOT} -name '*.proto')
pbs=$(protos:.proto=.pb.go)
proto: protopkg $(pbs)

# Documentation generation

DOC_ROOT=doc
DOC_CLI_ROOT=${DOC_ROOT}/cmd/cli
DOC_CLI_BIN=jobdoc
doc-cmd-cli:
	go build -o ${DOC_CLI_ROOT}/${DOC_CLI_BIN} ${MODULE_ROOT}/${DOC_CLI_ROOT}
	pushd ${DOC_CLI_ROOT} && ./${DOC_CLI_BIN} && popd
	rm ${DOC_CLI_ROOT}/${DOC_CLI_BIN}
	-git add ${DOC_CLI_ROOT}/*.md
	-git commit ${DOC_CLI_ROOT}/*.md -m "docs: update CLI documentation"

doc-proto:
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
	protoc --doc_out=./${DOC_ROOT} --doc_opt=markdown,grpc-api.md \
		--proto_path=${PKG_PROTO_ROOT}/proto/v1 \
		$(shell find ${PKG_PROTO_ROOT}/proto/v1 -name "*.proto")
	-git commit ${DOC_ROOT}/grpc-api.md -m "docs: update proto documentation"

cmd-docs: doc-cmd-cli

%.md : %.adoc
	asciidoctor -b html5 -a leveloffset=+1 -o - $< | \
	pandoc -t gfm --wrap=none -f html > $@
	git commit $@ $< -m "docs: update $(notdir $<)"

%.png : %.pu
	plantuml -tpng $<
	-git commit $@ $< -m "docs: update $(notdir $<)"

images := $(wildcard doc/img/*.png)
docs := $(wildcard doc/*.md)
doc: $(docs) $(images) cmd-docs doc-proto

# Valkey container management

.PHONY: valkey-start valkey-stop

VALKEY_CONTAINER_NAME ?= valkey
VALKEY_VERSION ?= 8.1.3
VALKEY_IMAGE ?= valkey/valkey:$(VALKEY_VERSION)
VALKEY_PORT ?= 6379

valkey-start:
	docker run -d --name $(VALKEY_CONTAINER_NAME) -p $(VALKEY_PORT):$(VALKEY_PORT) $(VALKEY_IMAGE)

valkey-stop:
	@docker stop $(VALKEY_CONTAINER_NAME) || true
	@docker rm $(VALKEY_CONTAINER_NAME) || true

# etcd container management

ETCD_CONTAINER_NAME ?= etcd
ETCD_VERSION ?= 3.6.4
ETCD_IMAGE ?= gcr.io/etcd-development/etcd:v$(ETCD_VERSION)
ETCD_PORT ?= 2379
ETCD_PEER_PORT ?= 2380

etcd-start:
	docker run -d --name $(ETCD_CONTAINER_NAME) -p $(ETCD_PORT):$(ETCD_PORT) -p $(ETCD_PEER_PORT):$(ETCD_PEER_PORT) \
	  $(ETCD_IMAGE) \
	  /usr/local/bin/etcd \
	  --name test-etcd \
	  --data-dir /etcd-data \
	  --advertise-client-urls http://0.0.0.0:2379 \
	  --listen-client-urls http://0.0.0.0:2379 \
	  --listen-peer-urls http://0.0.0.0:2380

etcd-stop:
	@docker stop etcd || true
	@docker rm etcd || true