# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
#


cmake_minimum_required(VERSION 2.8)
project(exception)

aux_source_directory(. SOURCE_FILES)
add_library(exception STATIC ${SOURCE_FILES})
SET(LIBRARY_OUTPUT_PATH ${OAMAGENT_LIBRARY_DIR})
