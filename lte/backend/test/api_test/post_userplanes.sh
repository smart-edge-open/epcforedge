#! /bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
#

setup_dir=${PWD}
echo "$setup_dir"
set -e

curl -X POST -H "Content-Type: application/json" --data @post_userplanes.json \
http://epc.oam:8080/userplanes | json_reformat

exit 0
