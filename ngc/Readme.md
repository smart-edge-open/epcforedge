```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019-2020 Intel Corporation 
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

| Param          | Description                                                                   |
| -------------- | ----------------------------------------------------------------------------  |
| TlsEndpoint    | HTTPS(TLS) EndPoint. (Not support for this release)                           |
| OpenEndpoint   | HTTP2 EndPoint. Used by CNCA to access OAM via HTTP2                          |
| NgcEndpoint    | NGC EndPoint. Used by OAM to access NGC                                       |
| NgcType        | NGC Type. Now only support APISTUB test mode.                                 |
| NgcTestData    | NGC TestData Path. Used by APISTUB testdata                                   |
| ServerCertPath | Path to SSL certs used by OAM server for CNCA/UI HTTP2 Requests               |
| ServerKeyPath  | Path to keys used by OAM server for HTTP2 connection between CNCA/UI and OAM  |

To run oam, just execute as below:
```sh
./oam
```

> NOTE: The OAM bin will load configuration from `configs/oam.json`, so before execution please have nef configuration file in the configs folder

### OAM Unit and API Testing

The Unit and API Testing should be performed in a local development environment. 

To run only unit tests and generate coverage report:
```sh
make test-unit-oam
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
### Build

To build af:

```sh
make af
```
Generated bin will be put under `ngc/dist`

AF sample code, generated bin will be put under ngc/dist


AF configurable parameters list:

| Param               | Description                                                                         |
| ------------------- | ----------------------------------------------------------------------------------- |
| AfID                | AF ID provided by OAM during AF registration                                        |
| CNCAEndpoint        | HTTP2 EndPoint. Used by CNCA to access AF via HTTP2                                 |
| Hostname            | Provided by AF to NEF. Part of URL used by  notifications                           |
| NotifPort           | NGC EndPoint. Used by NEF to send notifications to AF                               |
| UIEndpoint          | CNCA UI EndPoint Used by AF                                                         |
| ServerCertPath      | Path to certs used by AF server for CNCA HTTP2 Requests and NEF HTTP2 notification  |
| ServerKeyPath       | Path to keys used for HTTP2 connection between AF-CNCA and AF-NEF                   |
| Protocol            | Protocol used between AF and NEF                                                    |
| NEFHostname         | NEF Hostname used by AF                                                             |
| NEFPort             | NEF Port used by AF for sending requests to NEF                                     |
| NEFBasePath         | URL used by AF to access NEF Traffic Influence                                      |
| UserAgent           | Used by AF client connecting to NEF                                                 |
| NEFClientCertPath   | Path to certs used by AF client                                                     |
| LocationPrefixPfd   | The API prefix for PFD management                                                   |
| NEFPFDBasePath      | URL used by AF to access NEF PFD management                                         |
| OAuth2Support       | OAuth2 support in AF                                                                |

AF configurable parameters list for Policy Authorization (Indicated by CliPaConfig):

| Param               | Description                                                                         |
| ------------------- | ----------------------------------------------------------------------------------- |
| Protocol            | Protocol used between AF and PCF (http/https)                                       |
| ProtocolVer         | HTTP/s Protocol version (1.1/2.0)                                                   |
| Hostname            | PCF Hostname used by AF                                                             |
| Port                | PCF Port number to which AF sends request                                           |
| BasePath            | URL used by AF to access PCF for Policy authorization                               |
| LocationPrefixURI   | The API prefix for Policy Authorization                                             |
| CliCertPath         | Path for certificates to be used by AF client to communicate with PCF               |
| OAuth2Support       | OAuth2 authorization support between AF and PCF                                     |
| NotifURI            | Notification URL to which AF will recieve notifications from PCF                    |

To run af, just execute as below:
```sh
./af
```

> NOTE: The AF bin will load configuration from `configs/af.json`, so before execution please have nef configuration file in the configs folder

To run AF Ginkgo test suite run

```sh
make test-unit-af
```

To run af PFD functional tests, execute the test scripts in ngc/test/af/pfd_tests/:

```sh
./af_pfd_trans_tests.sh  - These are the PFD Transaction level tests GET/POST/GET_ALL/PUT/DELETE
./af_pfd_max_trans.sh - These are tests for maximum number of transactions supported
./af_pfd_app_tests.sh - These are the PFD application level tests GET/PUT/DELETE/PATCH
./af_pfd_all_tests.sh - This is a script which calls all the above tests
```

> NOTE: The config file contains the url, hostname and port information for AF/NEF. 
> The AF bin and NEF bin must be started prior to test execution with the same config as in config file


## NEF

### Build

To build nef:

```sh
make nef
```
Generated bin will be put under `ngc/dist`

### Configure and Run NEF bin

NEF Config Example File - `nef.json` located at `ngc/configs`

#### The configurable parameters list:

| Param                     | Description                                                                                                                                                             |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| NefAPIRoot                | The API root of the NEF i.e. ip address or domain name                                                                                                                  |
| LocationPrefix            | The API prefix for the traffic influence subscription. The NefAPIRoot + Endpoint + LocationPrefix + subscription id generated by NEF form the subscription resource uri |
| MaxSubSupport             | The maximum number of subscriptions to be supported by NEF.                                                                                                             |
| MaxAFSupport              | The maximum number of AF's to be supported by NEF                                                                                                                       |
| SubStartId                | The start value of  the subscription ids                                                                                                                                |
| UpfNotificationResUriPath | The API path on which the the NEF would listen for UPF notifications from SMF. The NefAPI + EndPoint + UpfNotificationResUriPath together form the notification URI     |
| UserAgent                 | The user agent information to put in the HTTP requests                                                                                                                  |
| HTTPConfig                | The fields under this describe the configuration for the HTTP Endpoint                                                                                                  |
| Endpoint                  | The end point where the NEF server needs to listen for HTTP 1.1 requests.Format ipaddress:port                                                                          |
| HTTP2Config               | The fields under this describe the configuration for the HTTP2 Endpoint                                                                                                 |
| Endpoint                  | The end point where the NEF server needs to listen for HTTP 2.1 requests. Format ipaddress:port                                                                         |
| NefServerCert             | The file path containing the NEF Server public key                                                                                                                      |
| NefServerKey              | The file path containing the NEF Server private key                                                                                                                     |
| AfClientCert              | The file path containing the AF Server public key                                                                                                                       |
| AfServiceID               | List of mappings for AfServiceID to dnn and snssai, used by NEF when communicating with UDR and PCF                                                                     |
| id                        | The AF Service ID                                                                                                                                                       |
| dnn                       | Data network name                                                                                                                                                       |
| snssai                    | Single Network Slice Selection Assistance Information                                                                                                                   |
| LocationPrefixPfd         | The API prefix for PFD management. The NefAPIRoot + Endpoint + LocationPrefixPfd + transaction id generated by NEF forms the PFD resource uri                           |
| MaxPfdTransSupport        | The maximum number of PFD transactions to be supported by NEF.                                                                                                          |
| PfdTransStartID           | The start value of  the PFD transaction ids                                                                                                                             |
| OAuth2Support             | OAuth2 support in AF                                                                                                                                                    |

#### Run NEF
To run nef, just execute as below:
```sh
./nef
```
> NOTE: 
1. The NEF will load configuration from `/configs/nef.json`, so before execution please have nef configuration file in the configs folder
2. The NEF certificates need to be available in the location mentioned in the configuration

### NEF Unit and API Testing

The Unit and API Testing should be performed in a local development environment. 

To run only unit tests and generate coverage report:
```sh
make test-unit-nef
```

> NOTE: Before executing unit test, need to install ginkgo as below:

```sh
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega/...
```
To run API Testing:

Step 1 : Change to the nef folder
```sh
cd ngc/test/nef/nef-cli-scripts
```
Step 2 : Run Curl Test Scripts to simulate HTTP Request to NEF.

## CNTF

### Build

To build cntf:

```sh
make cntf
```
Generated bin will be put under `ngc/dist`

### Configure and Run CNTF bin

CNTF Config Example File - `cntf.json` located at `ngc/configs`

#### The configurable parameters list:

| Param                     | Description                                                                                                                                                             |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| CntfAPIRoot                | The API root of the CNTF i.e. ip address or domain name                                                                                                                  |
| LocationPrefix            | The API prefix  |
| MaxASCSupport             | The maximum number of Application Session Context supported by CNTF |
| HTTPConfig                | The fields under this describe the configuration for the HTTP Endpoint                                                                                                  |
| Endpoint                  | The end point where the CNTF server needs to listen for HTTP 1.1 requests.Format ipaddress:port                                                                          |
| HTTP2Config               | The fields under this describe the configuration for the HTTP2 Endpoint                                                                                                 |
| Endpoint                  | The end point where the CNTF server needs to listen for HTTP 2.1 requests. Format ipaddress:port                                                                         |
| CNTFServerCert             | The file path containing the CNTF Server public key                                                                                                                      |
| CNTFServerKey              | The file path containing the CNTF Server private key                                                                                                                     |
| AfClientCert              | The file path containing the AF Server public key                                                                                                                       |
| OAuth2Support             | OAuth2 support in CNTF                                                                                                                                                   |

#### Run CNTF
To run cntf, just execute as below:
```sh
./cntf
```
> NOTE: 
1. The CNTF will load configuration from `/configs/cntf.json`, so before execution please have nef configuration file in the configs folder
2. The CNTF certificates need to be available in the location mentioned in the configuration

### CNTF Unit and API Testing

The Unit and API Testing should be performed in a local development environment. 

To run only unit tests and generate coverage report:
```sh
make test-unit-cntf
```

> NOTE: Before executing unit test, need to install ginkgo as below:

```sh
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega/...
```
To run API Testing:

Step 1 : Change to the cntf folder
```sh
cd ngc/test/cntf/cntf-cli-scripts
```
Step 2 : Run Curl Test Scripts to simulate HTTP Request to CNTF.

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

## OAuth2

The OAuth2 is a go package used by NEF and AF for OAuth2 token generation and verification. By default OAuth2 is enabled between AF and NEF. OAuth2 can be disabled by respectinve AF and NEF component configuration. 

#### The configurable parameters list:
| Param      | Description                                                                                   |
| ---------- | --------------------------------------------------------------------------------------------- |
| SigningKey | The API root of the NEF i.e. ip address or domain name. The default signing key is "OPENNESS" |
| expiration | OAuth2 token expiration time                                                                  |
