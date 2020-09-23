#! /bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
#

setup_dir=${PWD}
echo "$setup_dir"
set -e

curl -v http://epc.oam:8080/userplanes/666
