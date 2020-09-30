#! /bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation


setup_dir=${PWD}
echo "$setup_dir"

set -e

curl -X GET -i "Content-Type: application/json"  http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions/11

exit 0
