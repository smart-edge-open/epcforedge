# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
#

cmake_minimum_required(VERSION 2.8)
project(OAMAGENT_TEST)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -fprofile-arcs -ftest-coverage")
set(CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG}  -fprofile-arcs -ftest-coverage")
set(CMAKE_C_FLAGS_DEBUG "${CMAKE_C_FLAGS_DEBUG}  -fprofile-arcs -ftest-coverage")
set(CMAKE_EXE_LINKER_FLAGS_DEBUG "${CMAKE_EXE_LINKER_FLAGS_DEBUG}  -fprofile-arcs -ftest-coverage -lgcov")
set(JSONCPP_PATH /usr/include)
set(LIB64_PATH /usr/lib64)

include_directories(
        ${JSONCPP_PATH}/jsoncpp
        ./common
        ./test_cases
	./../../../utils/exception
        )

link_directories(
        /usr/local/lib
        )

set(TEST_SRC_FILES
        ./main.cpp
        ./common/TesterBase.cpp
        ./common/TesterFramework.cpp
        ./common/TestUtility.cpp
        ./../../../utils/exception/Exception.cpp
        ./test_cases/DelUserplanesTester.cpp
        ./test_cases/GetUserplanesTester.cpp
        ./test_cases/PostUserplanesTester.cpp	
        ./test_cases/PatchUserplanesTester.cpp		
        )


add_executable(oamagent-test ${TEST_SRC_FILES})

target_link_libraries(oamagent-test
        jsoncpp
        fcgi++
        fcgi
        /usr/lib64/libboost_system.so
        /usr/lib64/libboost_thread.so
        crypto
	ssl
        pthread
        )
