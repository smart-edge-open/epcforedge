```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation and Smart-Edge.com, Inc.
```
# 1. Introduction
## Directory Structure
- `/cmd` : Main applications inside. 
- `/pkg` : Lib code used by applications. Perphaps common libs such as lib or utils in the folder. 
- `/dist` : Built golang binaries inside. 
- `/test` : Test apps and Test data inside. 
- `/configs` : Configuration files inside. 
- `/scripts` : Scripts files inside. 


## How to build oam sample bin and run it
### make oam
OAM sample code, generated bin will be put under ngc/dist


### make af
AF sample code, generated bin will be put under ngc/dist


### make test-unit
It is used to run unit test and generate coverage report.

Need to install ginkgo and generated unit test template as below:
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega/...
$ ginkgo bootstrap 

### make lint
need to install golangci-lint
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.21.0

### make help
IMPORTANT: understand make options by this help command.
