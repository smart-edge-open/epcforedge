```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation.
```

# backend

This is the HTTP backend processing implementation for oamagent.

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
- To communcate with oamagent from remote server, need to use generated certification crt file and add hostname "epc.oam" into the /etc/hosts file. The IP Address should be oamagent IP Address. 
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
- API Test: Provides CURL based test scripts for MEC Controller CUPS  API testing
- Unit Test: Provides unit test and code coverage  
- The details refer to README.md in the folders: api_test and unit_test. 

## Controller Integration Test without EPC

- Macro INT_TEST is flag to enable oamagent run without EPC, and respond with controller with pre-defined message in json format
- Enable INT_TEST in CMakeLists.txt as below:
```text
add_definitions(-g -Wall)
add_definitions(-O3)
add_definitions(-DCUPS_API_INT64_TYPE)
#add_definitions(-DINT_TEST)
```
- Then make 
- Copy json_payload folder (backend/test/unit_test/test_app/json_payload) into the same folder with oamagent
- Run oamagent as guide
- Limitation: Actually on INT_TEST mode, all the user planes are pre-configured by json files in json_payload. 
  So it is not flexible and just for some basic interface testing as below:
  - GET ALL:       user planes configuration are from PgwGetAllRspData.json and SgwGetAllRspData.json. And they can be changed according to test requirements.
  - GET by ID:     must comply with : PgwGetOneRspData.json and SgwGetOneRspData.json。And they can be changed according to test requirements.
  - Delete by ID:  hardcoded success response . No corresponding json file
  - POST:  hardcoded success only for user plane id 5. No corresponding json file
  - PATCH:  hardcoded success only for user plane id 5. No corresponding json file

  
 
