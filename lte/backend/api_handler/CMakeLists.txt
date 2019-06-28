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

#include_directories(
#        database
#        handler)

set(SOURCE_FILES
        main.cpp
        CpfInterface.cpp
        CupsMgmtInterface.cpp
        LocalConfig.cpp
        UpfController.cpp)

add_executable(oamagent ${SOURCE_FILES})
SET(EXECUTABLE_OUTPUT_PATH ${OAMAGENT_BINARY_DIR})

target_link_libraries(oamagent
        # oamagent libs
        liblog.a
        libexception.a
        libframework.a
        # system libs
        jsoncpp
        fcgi++
        fcgi
        curl 
        ssl
        crypto
        boost_system
        boost_thread
        pthread)
