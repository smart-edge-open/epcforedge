# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation

#! /bin/sh


setup_dir=${PWD}

set -e

# Create 2 AF

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_01.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_02.json http://localhost:8091/3gpp-traffic-influence/v1/AF_02/subscriptions

exit 0

