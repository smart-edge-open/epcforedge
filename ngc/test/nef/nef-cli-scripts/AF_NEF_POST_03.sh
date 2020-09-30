#! /bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation


setup_dir=${PWD}
echo "$setup_dir"

set -e

# Create 5 subscription

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_01.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_02.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_03.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_04.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_05.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_06.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

exit 0

