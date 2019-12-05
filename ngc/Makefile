################################################################################
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
################################################################################


export GO111MODULE = on

.PHONY: build af oam nef lint test-unit-af test-unit-nef test-unit-af
TMP_DIR:=$(shell mktemp -d)
BUILD_DIR ?=dist

VER:=1.0

build: af nef oam

af:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/af" ./cmd/af

nef:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/nef" ./cmd/nef

oam:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/oam" ./cmd/oam

lint:
	golangci-lint run ./pkg/oam ./pkg/nef ./pkg/af

test-unit-af:
	ginkgo -v -cover ./pkg/af

test-unit-nef:
	ginkgo -v -cover ./pkg/nef

test-unit-oam:
	ginkgo -v -cover ./pkg/config ./pkg/oam

test-unit-nef:
	ginkgo -v -cover ./pkg/nef

test-unit-af:
	ginkgo -v -cover ./pkg/af

clean:
	rm -rf "${BUILD_DIR}"


help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "Building Packages:"
	@echo "  build                  to build all"
	@echo "  af                     to build the AF"
	@echo "  nef                    to build the NEF"
	@echo "  oam                    to build the OAM"
	@echo ""
	@echo "Executing Test:"
	@echo "  lint                   to run static code analysis with golangci-lint"
	@echo "  test-unit-af           to run AF unit tests with ginkgo"
	@echo "  test-unit-nef          to run NEF unit tests with ginkgo"
	@echo "  test-unit-oam          to run OAM unit tests with ginkgo"
