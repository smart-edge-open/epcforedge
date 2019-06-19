``text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation and Smart-Edge.com, Inc.
```

# backend

This is the HTTPS backend processing implementation for oamagent.

## Directory Structure

- api_handler: The sub-directory contains source files to hanlde API message from/to MEC controller and EPC.
- utils: The sub-directory contains the source file to provide  utitlity functions.
- test: The sub-directory contains api test suites and unit test suites.

## Build

- Run setup_env.sh to install dependencies required for the building. 
- Type command:  make to build it. The output binary: oamagent will be in the build directory. 
- localconfig.json contains configuration information for oamagent.

## Run

- Run nginx according to README.md in the  http folder.
- To communcate with oamagent from remote server, need to use generated certification crt file and add hostname "mec.oam" into the /etc/hosts file. The IP Address should be oamagent IP Address. 
- Before run oamagent, need to set envionment varible as: export LD_LIBRARY_PATH=LD_LIBRARY_PATH:/usr/local/lib/
- Enter into build directory where oamagent locates, copy localconfig.json into build directory, then run oamagent directly such as:
```text
   export LD_LIBRARY_PATH=LD_LIBRARY_PATH:/usr/local/lib/
   ./oamagent &  
```
- The debug log will be in the /var/log/message

### Expected result

The expected result of running oamagent will show as below in the /var/log/message
```
Func:main(Line:115)OAMAgent Backend MgmtAPI registering ....
Func:main(Line:120)OAMAgent Backend running ...
```
### Unexpected result

- script set_env.sh failed - probably due to network connectivity.
- running oamagent failed - check whether localconfig.json format correctness and it should be same folder with oamagent.

## Test

In the test directory , there two types of tests:
- API Test: Provides CURL based test scripts for MEC Controller  API testing
- Unit Test: Provides unit test and code coverage  
- The details refer to README.md in the folders: api_test and unit_test. 
