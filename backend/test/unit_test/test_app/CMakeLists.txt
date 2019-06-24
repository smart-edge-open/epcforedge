############################################################################
# Copyright 2019 Intel Corporation. All rights reserved.
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
###########################################################################
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
