################################################################################
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
################################################################################


export GO111MODULE = on

.PHONY: build af oam nef lint test-unit-af test-unit-nef test-unit-oam
TMP_DIR:=$(shell mktemp -d)
BUILD_DIR ?=dist
CERTS_DIR ?=certs

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

# ./scripts/genCerts.sh -t DNS -n <HostName/KubernetesServiceName> -m <HostName/KubernetesServiceName>
# Replace the "localhost" in test-unit-af with hostname
# or Kubernetes service name as per user environment

test-unit-af:	nef
	@-pkill -f "nef"
	mkdir -p /etc/"${CERTS_DIR}"
	./scripts/genCerts.sh -t DNS -n localhost -m localhost
	mv root-ca-cert.pem root-ca-cert.srl root-ca-key.pem server-cert.pem server-key.pem server-request.csr /etc/${CERTS_DIR} 
	./dist/nef &
	ginkgo -v -cover ./pkg/af
	pkill -f "nef"
test-unit-nef:
	ginkgo -v -cover ./pkg/nef

test-unit-oam:
	ginkgo -v -cover ./pkg/config ./pkg/oam

clean:
	rm -rf "${BUILD_DIR}"
	rm -rf /etc/"${CERTS_DIR}"
	rm -rf extfile.cnf

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
