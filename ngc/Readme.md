```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation and Smart-Edge.com, Inc.
```
# 1. Introduction
## OAM Pkg source code
hello-world level sample code right now.
OAM core source files stay under ngc/pkg/oam directory.
OAM API-stub backend source files stay under ngc/test/oam/ directory.
OAM Flexcore backend source files stay under ngc/internal/flexcore/directory. 


## How to build oam sample bin and run it
* `make oam`
The bin will be generated under ngc/dist directory, used to build the OAM with API-stub backend (default mode).
Jus directly run bin as below:
`./oam-stub`


* `make oam-test-stub`
It is same with 'make oam', used to build the OAM with API-stub backend.

* `make oam-test-flexcore`
It is used to build the OAM with Flexcore backend.

### make af
AF sample code, generated bin will be put under ngc/dist

### make help
IMPORTANT: understand make options by this help command


### make test-unit
It is used to run unit test

need to install ginkgo and generated unit test template as below:
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega/...
$ ginkgo bootstrap 

### make lint
need to install golangci-lint
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.21.0
