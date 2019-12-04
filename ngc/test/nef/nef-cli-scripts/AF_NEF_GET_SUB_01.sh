# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation

#! /bin/sh

setup_dir=${PWD}

set -e

curl -X GET -i "Content-Type: application/json"  http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions/11112

exit 0
