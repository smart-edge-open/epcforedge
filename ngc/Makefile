################################################################################
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
################################################################################


export GO111MODULE = on

.PHONY: build af oam nef lint test-unit-af test-unit-nef test-unit-oam test-unit-cntf
TMP_DIR:=$(shell mktemp -d)
BUILD_DIR ?=dist
CERTS_DIR ?=/etc/certs

VER:=1.0

build: af nef oam cntf

af:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/af" ./cmd/af

nef:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/nef" ./cmd/nef

oam:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/oam" ./cmd/oam

cntf:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/cntf" ./cmd/cntf


lint:
	golangci-lint run ./pkg/oam ./pkg/nef ./pkg/af ./pkg/cntf

# ./scripts/genCerts.sh -t DNS -h localhost
# Replace the "localhost" in test-unit-af with hostname
# or Kubernetes service name as per user environment

test-unit-af:	nef
	@-pkill -f "nef"
	cp -rf configs ./pkg/af
	mkdir -p "${CERTS_DIR}"
	./scripts/genCerts.sh -t DNS -h localhost
	mv root-ca-cert.pem root-ca-cert.srl root-ca-key.pem server-cert.pem server-key.pem server-request.csr "${CERTS_DIR}"
	./dist/nef &
	ginkgo -v -cover ./pkg/af
	pkill -f "nef"
	rm -rf pkg/af/configs
test-unit-nef:
	ginkgo -v -cover ./pkg/nef

test-unit-oam:
	ginkgo -v -cover ./pkg/config ./pkg/oam

test-unit-cntf:
	ginkgo -v -cover ./pkg/cntf

clean:
	rm -rf "${BUILD_DIR}"
	rm -rf "${CERTS_DIR}"
	rm -rf extfile.cnf

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "Building Packages:"
	@echo "  build                  to build all"
	@echo "  af                     to build the AF"
	@echo "  nef                    to build the NEF"
	@echo "  oam                    to build the OAM"
	@echo "  cntf                	to build the cntf"
	@echo ""
	@echo "Executing Test:"
	@echo "  lint                   to run static code analysis with golangci-lint"
	@echo "  test-unit-af           to run AF unit tests with ginkgo"
	@echo "  test-unit-nef          to run NEF unit tests with ginkgo"
	@echo "  test-unit-oam          to run OAM unit tests with ginkgo"
	@echo "  test-unit-cntf      	to run cntf unit tests with ginkgo"
