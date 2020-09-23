#! /bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation


setup_dir=${PWD}
echo "$setup_dir"

set -e

# Create 5 subscription

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_UDR_01.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_UDR_02.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_UDR_03.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_UDR_04.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_UDR_05.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_UDR_06.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

exit 0

