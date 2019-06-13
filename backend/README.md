# backend

This is the HTTPS backend processing implementation for oamagent.

## Directory Structure

- api_handler: The sub-directory contains source files to hanlde API message from/to MEC controller and EPC.
- utils: The sub-directory contains the source file to provide  utitlity functions.
- test: The sub-directory contains api test suites and unit test suites.

## Build

1/ Run setup_env.sh to install dependencies required for the building. 
2/ Type command:  make to build it. The output binary: oamagent will be in the build directory. 
3/ localconfig.json contains configuration information for oamagent.

## Run

1/ Run nginx according to README.md in the  http folder
2/ Before run oamagent, need to set envionment varible as: export LD_LIBRARY_PATH=LD_LIBRARY_PATH:/usr/local/lib/
3/ Enter into build directory, copy localconfig.json , then run oamagent


## Test

In the test directory , there two types of tests:
1/ API Test: Provides CURL based test scripts for MEC Controller  API testing
2/ Unit Test: Provides unit test and code coverage  
The details refer to README.md in the folder: api_test and unit_test. 
