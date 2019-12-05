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

# 2. Prerequisites
## OS
- CentOS (7.6.1810)

## Tools
- `make` (version 3.81 or higher)
- `golang` (version 1.12.4)
- `openssl` (version 1.0.2 or higher)
- `curl` (version 7.29.0 or higher)

# 3. Quick Start
## OAM
### Build

To build oam:

```sh
make oam
```
Generated bin will be put under `ngc/dist`

### Configure and Run OAM Bin

OAM Config Example File - `oam.json` locates at `ngc/configs`

The configurable parameters list:

| Param              | Description                                              |
|--------------------|----------------------------------------------------------|
| TlsEndpoint        | HTTPS(TLS) EndPoint. (Not support for this release)      |
| OpenEndpoint       | HTTP EndPoint. Used by CNCA to access OAM via HTTP       |
| NgcEndpoint        | NGC EndPoint. Used by OAM to access NGC                  |
| NgcType            | NGC Type. Now only support APISTUB test mode.            |
| NgcTestData        | NGC TestData Path. Used by APISTUB testdata              |

To run oam, just execute as below:
```sh
./dist/oam
```

> NOTE: The OAM bin will load configuration from `ngc/configs/oam.json`.

### OAM Unit and API Testing

The Unit and API Testing should be performed in a local development environment. 

To run only unit tests and generate coverage report:
```sh
make test-unit
```

> NOTE: Before executing unit test, need to install ginkgo as below:

```sh
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega/...
```


To run API Testing:

Step 1: Run OAM with NgcType as `APISTUB`.

Step 2: Run Curl Test Scripts to simulate HTTP Request to OAM.
```sh
cd ngc/test/oam/cnca-cli-scripts
./runAll.sh
```
> NOTE: Can use ./runAll.sh to run all methods testing automatically. Also can use ./cliTest.shto run different method manually.


## AF
### make af
AF sample code, generated bin will be put under ngc/dist



## Lint

```sh
make lint
```
> NOTE: Need to install golangci-lint as below:
```sh
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.21.0
```

## make help
IMPORTANT: understand make options by this help command.


## Certifications Generation

To generate certifications for TLS, can use shell script - `genCerts.sh` in `ngc/scripts`. 
Can get detail usage by executing the script as: `./genCerts.sh ?` .
