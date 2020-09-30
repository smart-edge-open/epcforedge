#! /bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation


setup_dir=${PWD}
echo "$setup_dir"

set -e

curl -X PUT -i "Content-Type: application/json" --data @./json/AF_NEF_PUT_UDR_01.json http://localhost:8091/3gpp-traffic-influence/v1/AF_02/subscriptions/11111

exit 0

