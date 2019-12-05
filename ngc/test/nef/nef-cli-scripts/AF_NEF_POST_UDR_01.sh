# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation

#! /bin/sh

setup_dir=${PWD}

set -e

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_UDR_01.json http://localhost:8091/3gpp-traffic-influence/v1/AF_02/subscriptions

#curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_02.json http://localhost:8090/3gpp-traffic-influence/v1/AF_01/subscriptions

exit 0

