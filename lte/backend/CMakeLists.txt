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
############################################################################

cmake_minimum_required(VERSION 2.8)
project(oamagent)

message([STATUS|WARNING|AUTHOR_WARNING|FATAL_ERROR|SEND_ERROR]
          "message to display" ...)

add_definitions(-fstack-protector-strong -fPIC -D_FORTIFY_SOURCE=2 -Wformat -Wformat-security)
set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -z noexecstack -z relro -z now")

add_definitions(-g -Wall)
add_definitions(-O3)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")

# oamagent source dir
set(OAMAGENT_DIR ${CMAKE_CURRENT_SOURCE_DIR})

# oamagent build dir
set(OAMAGENT_BUILD_DIR ${CMAKE_CURRENT_SOURCE_DIR}/build)
set(OAMAGENT_BINARY_DIR ${OAMAGENT_BUILD_DIR})
set(OAMAGENT_LIBRARY_DIR ${OAMAGENT_BUILD_DIR}/oamagent_utils)


# system lib dir
set(JSONCPP_PATH /usr/include)
set(LIB64_PATH /usr/lib64)

# include directories
include_directories(
        ${OAMAGENT_DIR}/utils/log
        ${OAMAGENT_DIR}/utils/exception
        ${OAMAGENT_DIR}/utils/framework
        ${JSONCPP_PATH}/jsoncpp)

# lib directories
LINK_DIRECTORIES(
        ${OAMAGENT_LIBRARY_DIR}
        ${JSONCPP_PATH}
        ${LIB64_PATH})

ADD_SUBDIRECTORY(utils)
ADD_SUBDIRECTORY(api_handler)
