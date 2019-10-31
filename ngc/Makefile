################################################################################
# Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

#GOOS=linux go build -o "${BUILD_DIR}/oam-flexcore" ./pkg/oam ./internal/flexcore/oam_backend.go

export GO111MODULE = on

.PHONY: build af oam oam-test oam-test-flexcore lint test-unit
TMP_DIR:=$(shell mktemp -d)
BUILD_DIR ?=dist

VER:=1.0

build: af nef oam

af:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/af" ./pkg/af

nef:
	mkdir -p "${BUILD_DIR}"
	GOOS=linux go build -o "${BUILD_DIR}/nef" ./pkg/nef

oam:
	mkdir -p "${BUILD_DIR}"
	rm -f ./pkg/oam/oam_backend.go
	ln -s ../../internal/flexcore/oam_flexcore_backend.go ./pkg/oam/oam_backend.go
	GOOS=linux go build -o "${BUILD_DIR}/oam-stub" ./pkg/oam

oam-test-stub:
	mkdir -p "${BUILD_DIR}"
	rm -f ./pkg/oam/oam_backend.go
	ln -s ../../internal/flexcore/oam_flexcore_backend.go ./pkg/oam/oam_backend.go
	GOOS=linux go build -o "${BUILD_DIR}/oam-stub" ./pkg/oam

oam-test-flexcore:
	mkdir -p "${BUILD_DIR}"
	rm -f ./pkg/oam/oam_backend.go
	ln -s ../../test/oam/ngc_stub.go ./pkg/oam/oam_backend.go
	GOOS=linux go build -o "${BUILD_DIR}/oam-flexcore" ./pkg/oam


lint:
	golangci-lint run

test-unit:
	ginkgo ./pkg/oam

clean:
	rm -rf "${BUILD_DIR}"


help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "Building Packages:"
	@echo "  build                  to build all"
	@echo "  af                     to build the AF"
	@echo "  oam                    to build the OAM with API-stub backend (default mode)"
	@echo "  oam-test-stub          to build the OAM with API-stub backend"
	@echo "  oam-test-flexcore      to build the OAM with Flexcore backend"
	@echo ""
	@echo "Executing Test:"
	@echo "  lint                   to run static code analysis with golangci-lint"
	@echo "  test-unit              to run unit tests with ginkgo"
