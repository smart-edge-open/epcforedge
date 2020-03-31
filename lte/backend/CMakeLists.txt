# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
#


cmake_minimum_required(VERSION 2.8)
project(oamagent)

message([STATUS|WARNING|AUTHOR_WARNING|FATAL_ERROR|SEND_ERROR]
          "message to display" ...)

add_definitions(-fstack-protector-strong -fPIC -D_FORTIFY_SOURCE=2 -Wformat -Wformat-security)
set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -z noexecstack -z relro -z now")

add_definitions(-g -Wall)
add_definitions(-O3)
add_definitions(-DCUPS_API_INT64_TYPE)
#add_definitions(-DINT_TEST)


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
